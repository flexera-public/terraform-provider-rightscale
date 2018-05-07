package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// resource "rightscale_route_table" "my_route_table" {
//   name        = "my-sweet-route-table"
//   description = "This is the best route table ever"
//   cloud_href = "${data.rightscale_cloud.us-oregon.href}"
//	 network_href = "${rightscale_network.us-oregon-vpc-network.href}"
// }

func resourceRouteTable() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceCreateFunc("rs_cm", "route_tables", routeTableWriteFields),
		Update: resourceUpdateFunc(routeTableWriteFields),

		Schema: map[string]*schema.Schema{
			"cloud_href": &schema.Schema{
				Description: "Href of cloud in which to create route table",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": &schema.Schema{
				Description: "Description of route table",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": &schema.Schema{
				Description: "Name of route table",
				Type:        schema.TypeString,
				Required:    true,
			},
			"network_href": &schema.Schema{
				Description: "Href of network in which to create route table",
				Type:        schema.TypeString,
				Required:    true,
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
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"href": {
				Type:        schema.TypeString,
				Description: "Href of route table",
				Computed:    true,
			},
		},
	}
}

func routeTableWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{
		"cloud_href":   d.Get("cloud_href"),
		"name":         d.Get("name"),
		"network_href": d.Get("network_href"),
	}
	if desc, ok := d.GetOk("description"); ok {
		fields["description"] = desc
	}
	return rsc.Fields{"route_table": fields}
}
