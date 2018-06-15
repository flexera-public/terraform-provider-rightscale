package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_image" "centos_7" {
//   cloud_href = ${data.rightscale_cloud.gce.id}
//   filter {
//     name = "centos 7"
//     visibility = "public"
//   }
// }

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		Read: resourceImageRead,

		Schema: map[string]*schema.Schema{
			"cloud_href": {
				Type:        schema.TypeString,
				Description: "ID of image cloud resource",
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
							Description: "name of image, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "description of image, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"resource_uid": {
							Type:        schema.TypeString,
							Description: "cloud ID - if this filter is set additional retry logic will fire to allow for cloud resource discovery",
							Optional:    true,
							ForceNew:    true,
						},
						"os_platform": {
							Type:         schema.TypeString,
							Description:  "The image's operating system to filter on. Examples: Linux or Windows.",
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"windows", "linux"}, false),
						},
						"cpu_architecture": {
							Type:        schema.TypeString,
							Description: "CPU architecture of image, e.g. 'x86_64', uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"image_type": {
							Type:         schema.TypeString,
							Description:  `The Image Type to filter on. This will be either "machine", "machine_azure", "ramdisk" or "kernel"`,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"machine", "machine_azure", "ramdisk", "kernel"}, false),
						},
						"visibility": {
							Type:         schema.TypeString,
							Description:  "The visibility of the Image in the cloud to filter on, defaults to 'private.' Options: private, public.",
							Optional:     true,
							ForceNew:     true,
							Default:      "private",
							ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
						},
					},
				},
			},

			// Read-only fields
			"cpu_architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_type": {
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
			"os_platform": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_device_storage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtualization_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
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

func resourceImageRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	cloud := d.Get("cloud_href").(string)
	loc := &rsc.Locator{Namespace: "rs_cm", Href: cloud}

	// if 'resource_uid' filter is set, we expect it to show up.
	// retry for 5 min to allow rightscale to poll cloud to discover.
	if uidset := cmUIDSet(d); uidset {
		timeout := 300
		err := cmIndexRetry(client, loc, "images", d, timeout)
		if err != nil {
			return err
		}
	}

	res, err := client.List(loc, "images", cmFilters(d))
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
