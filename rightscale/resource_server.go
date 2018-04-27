package rightscale

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceCreateFunc("rs_cm", "servers", serverWriteFields),
		Update: resourceUpdateFunc(serverWriteFields),

		Schema: map[string]*schema.Schema{
			"deployment_href": &schema.Schema{
				Description: "ID of deployment in which to create server",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": &schema.Schema{
				Description: "description of server",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance": &schema.Schema{
				Description: "server instance details",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem:        resourceInstance(),
			},
			"inputs": &schema.Schema{
				Description: "Inputs associated with an instance when incarnated from a server or server array - should be rendered via template provider - see docs",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				// Uncomment below then they start supporting this on TypeSet
				//ValidateFunc: validation.StringMatch(regexp.MustCompile("\\w+=\\w+:\\w+"), "values must be in format of 'KEY=type:value'"),
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
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
		},
	}
}

func serverWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{}
	// construct 'instance' hash so we end up with a server WITH a running instance
	if i, ok := d.GetOk("instance"); ok {
		fields["instance"] = i.([]interface{})[0].(map[string]interface{})
		if fields["instance"].(map[string]interface{})["associate_public_ip_address"].(bool) == false {
			delete(fields["instance"].(map[string]interface{}), "associate_public_ip_address")
		}
	}
	// if inputs are defined in resource, add those to instance field as proper format
	if _, ok := d.GetOk("inputs"); ok {
		if r, ok := cmInputs(d); ok != nil {
			log.Printf("[ERROR]: %v", ok)
		} else {
			fields["instance"].(map[string]interface{})["inputs"] = r["inputs"]
		}
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
		"server_tag_scope",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return rsc.Fields{"server": fields}
}
