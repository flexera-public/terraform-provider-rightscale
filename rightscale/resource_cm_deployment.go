package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

func resourceCMDeployment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete, // can fail if deployment is locked - that's what we want
		Create: resourceCMDeploymentCreate,
		Update: resourceCMDeploymentUpdate,

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
				InputDefault: "deployment",
				ValidateFunc: validation.StringInSlice([]string{"account", "deployment"}, false),
			},

			// Read-only fields
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
		},
	}
}

func resourceCMDeploymentCreate(d *schema.ResourceData, m interface{}) error {
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

	return nil
}

func resourceCMDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
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
