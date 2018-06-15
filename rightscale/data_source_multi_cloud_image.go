package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_multi_cloud_image" "centos_64" {
//   filter {
//     name = "RightImage_CentOS_6.4_x64_v13.5"
//     revision = 43
//   }
// }

func dataSourceMultiCloudImage() *schema.Resource {
	return &schema.Resource{
		Read: resourceMultiCloudImageRead,

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
							Description: "name of multi-cloud image, partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"revision": {
							Type:        schema.TypeInt,
							Description: "revision of multi-cloud image, use 0 to match latest non-committed version",
							Optional:    true,
							ForceNew:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "description of multi-cloud image, partial match",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"server_template_href": {
				Type:        schema.TypeString,
				Description: "ID of image's server template resource",
				Optional:    true,
				ForceNew:    true,
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"href": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMultiCloudImageRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	loc := &rsc.Locator{Namespace: "rs_cm"}
	link := ""
	if st, ok := d.GetOk("server_template_href"); ok {
		loc.Href = st.(string)
		link = "multi_cloud_images"
	} else {
		loc.Type = "multi_cloud_images"
	}

	res, err := client.List(loc, link, cmFilters(d))
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
