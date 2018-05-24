package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_subnet" "infrastructure" {
//   cloud_href = ${data.rightscale_cloud.ec2_us_east_1.id}
//   filter {
//     name = "infrastructure"
//   }
// }

func dataSourceSubnet() *schema.Resource {
	return &schema.Resource{
		Read: resourceSubnetRead,

		Schema: map[string]*schema.Schema{
			"cloud_href": {
				Type:        schema.TypeString,
				Description: "ID of the subnet cloud",
				Required:    true,
				ForceNew:    true,
			},
			"filter": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of subnet, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"resource_uid": {
							Type:        schema.TypeString,
							Description: "cloud ID - if this filter is set additional retry logic will fire to allow for cloud resource discovery",
							Optional:    true,
							ForceNew:    true,
						},
						"datacenter_href": {
							Type:        schema.TypeString,
							Description: "Href of the subnet datacenter resource",
							Optional:    true,
							ForceNew:    true,
						},
						"instance_href": {
							Type:        schema.TypeString,
							Description: "Href of instance resource attached to subnet",
							Optional:    true,
							ForceNew:    true,
						},
						"network_href": {
							Type:        schema.TypeString,
							Description: "Href of network resource that owns subnet",
							Optional:    true,
							ForceNew:    true,
						},
						"visibility": {
							Type:        schema.TypeString,
							Description: "Visibility of the subnet to filter by (private, shared, etc)",
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSubnetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	cloud := d.Get("cloud_href").(string)
	loc := &rsc.Locator{Namespace: "rs_cm", Href: cloud}

	// if 'resource_uid' filter is set, we expect it to show up.
	// retry for 5 min to allow rightscale to poll cloud to discover.
	if uidset := cmUIDSet(d); uidset {
		timeout := 300
		err := cmIndexRetry(client, loc, "subnets", d, timeout)
		if err != nil {
			return err
		}
	}

	res, err := client.List(loc, "subnets", cmFilters(d))
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
