package rightscale

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// resource "rightscale_deployment" "my_deployment" {
//   name        = "my-test-deployment"
//   description = "The quick brown fox jumped over the lazy dogs"
// }

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDeploymentDelete,
		Create: resourceDeploymentCreate,
		Update: resourceDeploymentUpdate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Description: "name of deployment",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": &schema.Schema{
				Description: "description of deployment",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"resource_group_href": &schema.Schema{
				Description: "ID of the Windows Azure Resource Group attached to the deployment",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"locked": &schema.Schema{
				Description: "whether deployment is locked",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"server_tag_scope": &schema.Schema{
				Description:  "routing scope for tags for servers in the deployment",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "deployment",
				ValidateFunc: validation.StringInSlice([]string{"account", "deployment"}, false),
			},

			// Read-only fields
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"href": &schema.Schema{
				Description: "href of deployment",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	var mustLock bool
	{
		locked, ok := d.GetOk("locked")
		mustLock = ok && locked.(bool)
	}

	client := m.(rsc.Client)
	res, err := client.Create("rs_cm", "deployments", deploymentWriteFields(d))
	if err != nil {
		return err
	}
	d.SetId(res.Locator.Namespace + ":" + res.Locator.Href)
	for k, v := range res.Fields {
		d.Set(k, v)
	}

	if mustLock {
		d.Set("locked", true)
		if err := updateLock(d, client, "deployments"); err != nil {
			d.SetId("")
			d.Set("locked", false)
			// Attempt to delete previously created deployment, ignore errors
			client.Delete(res.Locator)
			return err
		}
	}
	d.Set("href", res.Locator.Href)
	return nil
}

func resourceDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	client := m.(rsc.Client)
	loc, err := locator(d)
	if err != nil {
		return err
	}

	// update lock
	if err := updateLock(d, client, "deployments"); err != nil {
		return handleRSCError(d, err)
	}
	d.SetPartial("locked")

	// then the other fields
	if err := client.Update(loc, deploymentWriteFields(d)); err != nil {
		return handleRSCError(d, err)
	}

	d.Partial(false)
	return nil
}

func deploymentWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{"name": d.Get("name")}
	for _, f := range []string{"description", "resource_group_href", "server_tag_scope"} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return rsc.Fields{"deployment": fields}
}

// can fail if deployment is locked - that's what we want
// however assets in 'terminating' state will eventually clear so retry on 422 for period of time
func resourceDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	loc, err := locator(d)
	if err != nil {
		return err
	}

	// wrap in rescue/retry if response is 422 with max ttl
	timeout := time.After(5 * time.Minute)
	tick := time.Tick(10 * time.Second)
	log.Printf("[INFO] Deleting Deployment - waiting up to 5 min for objects to finish being destroyed in deployment: %s", d.Id())
	for {
		select {
		case <-timeout:
			// 5 minutes expired - raise and exit
			return client.Delete(loc)

		case <-tick:
			err := client.Delete(loc)
			if err == nil {
				// successful deletion - exit retry/timeout loop
				return nil
			}
			// Search errorresponse for specific string indicating the deployment still contains objects that are still 'terminating'
			// If error message contains 'ActionNotAllowed: This deployment cannot be deleted because it contains running servers and/or active arrays.' retry,
			// otherwise on any other error raise and exit.
			if strings.Contains(err.Error(), "ActionNotAllowed: This deployment cannot be deleted because it contains running servers and/or active arrays.") {
				log.Printf("[INFO] Deleting Deployment - 422 from cm api - instances still active in deployment - try again later")
			} else {
				// Unhandled error from API that should be immediately returned eg deployment locked
				return client.Delete(loc)
			}
		}
	}
}
