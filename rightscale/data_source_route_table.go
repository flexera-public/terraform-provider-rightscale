package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_route_table" "infra_vpc_route_table" {
//   filter {
//     name = "rtb-123456"
//     network_href = ${data.rightscale_network.ec2_us_east_vpc.id}
//   }
// }

func dataSourceRouteTable() *schema.Resource {
	return &schema.Resource{
		Read: resourceRouteTableRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "name of route table, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"cloud_href": {
							Type:        schema.TypeString,
							Description: "Href of the route table's cloud",
							Optional:    true,
							ForceNew:    true,
						},
						"network_href": {
							Type:        schema.TypeString,
							Description: "Href of the route table's network",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},

			// Read-only fields
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"href": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRouteTableRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	loc := &rsc.Locator{Namespace: "rs_cm", Type: "route_tables"}

	res, err := client.List(loc, "", cmFilters(d))
	if err != nil {
		return err
	}

	if len(res) == 0 {
		return nil
	}
	for k, v := range res[0].Fields {
		d.Set(k, v)
	}
	d.SetId(res[0].Locator.Href)
	d.Set("href", res[0].Locator.Href)
	return nil
}
