package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_volume" "mysql_master" {
//   filter {
//     name = "mysql_master"
//   }
//   cloud_href = ${data.rightscale_cloud.ec2_us_east_1.id}
// }

func dataSourceVolume() *schema.Resource {
	return &schema.Resource{
		Read: resourceVolumeRead,

		Schema: map[string]*schema.Schema{
			"cloud_href": {
				Type:        schema.TypeString,
				Description: "ID of the volume cloud",
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
							Description: "name of volume, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "description of volume, uses partial match",
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
							Description: "ID of the volume datacenter resource",
							Optional:    true,
							ForceNew:    true,
						},
						"deployment_href": {
							Type:        schema.TypeString,
							Description: "ID of deployment resource that owns volume",
							Optional:    true,
							ForceNew:    true,
						},
						"parent_volume_snapshot_href": {
							Type:        schema.TypeString,
							Description: "ID of volume snapshot that volume was created from",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},

			// Read-only fields
			"cloud_specific_attributes": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
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
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	cloud := d.Get("cloud_href").(string)
	loc := &rsc.Locator{Namespace: "rs_cm", Href: cloud}

	// if 'resource_uid' filter is set, we expect it to show up.
	// retry for 5 min to allow rightscale to poll cloud to discover.
	if uidset := cmUIDSet(d); uidset {
		timeout := 300
		err := cmIndexRetry(client, loc, "volumes", d, timeout)
		if err != nil {
			return err
		}
	}

	res, err := client.List(loc, "volumes", cmFilters(d))
	if err != nil {
		return err
	}

	if len(res) == 0 {
		return nil
	}
	for k, v := range res[0].Fields {
		if k == "cloud_specific_attributes" {
			d.Set(k, []interface{}{v})
			continue
		}
		d.Set(k, v)
	}
	d.SetId(res[0].Locator.Href)
	return nil
}
