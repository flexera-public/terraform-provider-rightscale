package rightscale

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// resource "rightscale_server" "web_server" {
//   name = "web_server"
//   deployment_href = "/api/deployments/1234"
//   instance {
//     cloud_href = "/api/clouds/1234"
//     image_href = "/api/clouds/1234/images/1234"
//     instance_type_href = "/api/clouds/1234/instance_types/1234"
//     name = "web_instance"
//     tags = [ "pantaloons"]
//     server_template_href = "/api/server_templates/1234"
//     inputs {
//       FOO = "text:bar"
//       BAZ = "cred:Bangarang"
//     }
//   }
// }

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDeleteServer,
		Create: resourceCreateServer(serverWriteFields),
		Update: resourceUpdateFunc(serverWriteFields),

		Schema: map[string]*schema.Schema{
			"deployment_href": &schema.Schema{
				Description: "ID of deployment in which to create server",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tags": &schema.Schema{
				Description: "Tags to apply to server and associated instance",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
			},
			"description": &schema.Schema{
				Description: "description of server",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance": &schema.Schema{
				Description: "server instance details",
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Elem:        resourceInstance(),
			},
			"name": &schema.Schema{
				Description: "name of server",
				Type:        schema.TypeString,
				Required:    true,
			},
			"optimized": &schema.Schema{
				Description: "A flag indicating whether Instances of this Server should be optimized for high-performance volumes (e.g. Volumes supporting a specified number of IOPS). Not supported in all Clouds.",
				Type:        schema.TypeBool,
				Optional:    true,
			},

			// Read-only fields
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"href": &schema.Schema{
				Description: "href of server",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func serverWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{}
	// construct 'instance' hash so we end up with a server WITH a running instance
	if i, ok := d.GetOk("instance"); ok {
		fields["instance"] = instanceWriteFieldsFromMap(i.([]interface{})[0].(map[string]interface{}))
	}
	if o, ok := d.GetOk("optimized"); ok {
		if o.(bool) {
			fields["optimized"] = "true"
		} else {
			fields["optimized"] = "false"
		}
	}
	for _, f := range []string{
		"deployment_href", "description", "name", "resource_group_href",
		"server_tag_scope", "tags",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return rsc.Fields{"server": fields}
}

func resourceCreateServer(fieldsFunc func(*schema.ResourceData) rsc.Fields) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		client := m.(rsc.Client)
		res, err := client.CreateServer("rs_cm", "servers", fieldsFunc(d))
		if err != nil {
			// Depending on where this failed we may or may not have an active cloud instance attached to the server object
			// Set partial for ID so we don't leave orphan instances for next apply operation.
			if res.Locator != nil && res.Locator.Href != "" {
				d.Partial(true)
				d.SetId(res.Locator.Namespace + ":" + res.Locator.Href)
				d.SetPartial("ID")
			}
			return err
		}
		for k, v := range res.Fields {
			d.Set(k, v)
		}
		// Sets 'href' which is rightscale href (for stitching together cm resources IN rightscale) without namespace.
		d.Set("href", res.Locator.Href)
		// Sets 'id' which allows terraform to locate the objects created which includes namespace.
		d.SetId(res.Locator.Namespace + ":" + res.Locator.Href)
		return nil
	}
}

func resourceDeleteServer(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	loc, err := locator(d)
	if err != nil {
		return err
	}

	// wrap in rescue/retry if response is 422 with max ttl
	timeout := time.After(5 * time.Minute)
	tick := time.Tick(10 * time.Second)
	log.Printf("[INFO] Deleting Server - will retry for 5 minutes in case connectivity to the server is lost during termination: %s", d.Id())
	for {
		select {
		case <-timeout:
			// 5 minutes expired - this is the last try
			return client.Delete(loc)

		case <-tick:
			err := client.Delete(loc)
			if err == nil {
				// successful deletion - exit retry/timeout loop
				return nil
			}
			// If the server's network is being destroyed while terminating the server, we get an error message
			// containing 'ActionNotAllowed: The server has a current instance.'
			// So if we see this error, we retry the Delete action during 5 minutes
			if strings.Contains(err.Error(), "ActionNotAllowed: The server has a current instance.") {
				log.Printf("[INFO] Deleting Server - received 422 from CM API - probably we lost connection to server - trying again")
			} else {
				// Unhandled error from API that should be immediately returned
				return err
			}
		}
	}
}
