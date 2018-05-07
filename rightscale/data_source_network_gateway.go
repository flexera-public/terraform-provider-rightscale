package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_network_gateway" "infra_vpc_gateway" {
//   filter {
//     name = "infrastructure-vpc-gateway"
//     cloud_href = ${data.rightscale_cloud.ec2_us_east_1.id}
//   }
// }

func dataSourceNetworkGateway() *schema.Resource {
	return &schema.Resource{
		Read: resourceNetworkGatewayRead,

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
							Description: "filter by name of network gateway, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"cloud_href": {
							Type:        schema.TypeString,
							Description: "filter by href of the specified cloud",
							Optional:    true,
							ForceNew:    true,
						},
						"network_href": {
							Type:        schema.TypeString,
							Description: "filter by network href",
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkGatewayRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	loc := &rsc.Locator{Namespace: "rs_cm", Type: "network_gateways"}

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
	return nil
}
