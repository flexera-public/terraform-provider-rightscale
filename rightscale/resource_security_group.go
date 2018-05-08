package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// resource "rightscale_security_group" "ssh" {
//     cloud_href = ${data.rightscale_cloud.ec2_us_east_1.id}
//     network_href = ${resource.network.my_network.id}
//     description = "my security group"
// }

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceCreateFunc("rs_cm", "security_groups", securityGroupWriteFields),

		Schema: map[string]*schema.Schema{
			"cloud_href": {
				Type:        schema.TypeString,
				Description: "href of the security group cloud",
				Required:    true,
				ForceNew:    true,
			},
			"deployment_href": {
				Type:        schema.TypeString,
				Description: "Optional href of deployment that owns the security group",
				Optional:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description of the security group",
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of the security group",
				Required:    true,
				ForceNew:    true,
			},
			"network_href": {
				Type:        schema.TypeString,
				Description: "href of the network to create the security group in",
				Required:    true,
				ForceNew:    true,
			},

			// Read-only fields
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"href": {
				Type:        schema.TypeString,
				Description: "href of security group",
				Computed:    true,
			},
		},
	}
}

func securityGroupWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{
		"cloud_href":   d.Get("cloud_href"),
		"name":         d.Get("name"),
		"network_href": d.Get("network_href"),
	}
	for _, f := range []string{
		"deployment_href", "description",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return rsc.Fields{"cloud_href": d.Get("cloud_href"), "security_group": fields}
}
