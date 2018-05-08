package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// resource "rightscale_route" "my_route" {
//   description = "A route to the internet gateway"
//   route_table_href = "${rightscale_route_table.my_route_table.href}"
//   destination_cidr_block = "0.0.0.0/0"
//   next_hop_type = "network_gateway"
//   next_hop_href = "${rightscale_network_gateway.my_network_gateway.href}"
// }

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceCreateFunc("rs_cm", "routes", routeWriteFields),
		Update: resourceUpdateFunc(routeWriteFields),

		Schema: map[string]*schema.Schema{
			"route_table_href": &schema.Schema{
				Description: "Href of route table in which to create new route",
				Type:        schema.TypeString,
				Required:    true,
			},
			"destination_cidr_block": &schema.Schema{
				Description: "Destination network in CIDR nodation",
				Type:        schema.TypeString,
				Required:    true,
			},
			"next_hop_type": &schema.Schema{
				Description:  "The route next hop type.  Options are 'instance', 'network_interface', 'network_gateway', 'ip_string', and 'url'",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"instance", "network_interface", "network_gateway", "ip_string", "url"}, false),
			},
			"next_hop_href": &schema.Schema{
				Description:   "The href of the Route's next hop. Required if 'next_hop_type' is 'instance', 'network_interface', or 'network_gateway'. Not allowed otherwise.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"next_hop_ip", "next_hop_url"},
			},
			"next_hop_ip": &schema.Schema{
				Description:   "The IP Address of the Route's next hop. Required if 'next_hop_type' is 'ip_string'. Not allowed otherwise.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"next_hop_url", "next_hop_href"},
			},
			"next_hop_url": &schema.Schema{
				Description:   "The URL of the Route's next hop. Required if 'next_hop_type' is 'url'. Not allowed otherwise.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"next_hop_ip", "next_hop_href"},
			},
			"description": &schema.Schema{
				Description: "Description of route",
				Type:        schema.TypeString,
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

func routeWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{
		"route_table_href":       d.Get("route_table_href"),
		"destination_cidr_block": d.Get("destination_cidr_block"),
		"next_hop_type":          d.Get("next_hop_type"),
	}
	for _, f := range []string{
		"description", "next_hop_href", "next_hop_ip", "next_hop_url",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return rsc.Fields{"route": fields}
}
