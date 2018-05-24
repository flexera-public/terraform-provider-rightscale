package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_network" "infra_vpc" {
//   filter {
//     resource_uid = "vpc-c31ee987"
//     cloud_href = ${data.rightscale_cloud.ec2_us_east_1.id}
//   }
// }

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read: resourceNetworkRead,

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
							Description: "name of network, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"cloud_href": {
							Type:        schema.TypeString,
							Description: "ID of the network cloud",
							Optional:    true,
							ForceNew:    true,
						},
						"deployment_href": {
							Type:        schema.TypeString,
							Description: "ID of deployment resource that owns network",
							Optional:    true,
							ForceNew:    true,
						},
						"resource_uid": {
							Type:        schema.TypeString,
							Description: "cloud ID - if this filter is set additional retry logic will fire to allow for cloud resource discovery",
							Optional:    true,
							ForceNew:    true,
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Description: "CIDR of the network resource",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},

			// Read-only fields
			"cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_tenancy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
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
		},
	}
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	loc := &rsc.Locator{Namespace: "rs_cm", Type: "networks"}

	// if 'resource_uid' filter is set, we expect it to show up.
	// retry for 5 min to allow rightscale to poll cloud to discover.
	if uidset := cmUIDSet(d); uidset {
		timeout := 300
		err := cmIndexRetry(client, loc, "networks", d, timeout)
		if err != nil {
			return err
		}
	}

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
