package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_server" "web_server" {
//   filter {
//     name = "web"
//   }
// }

func dataSourceServer() *schema.Resource {
	return &schema.Resource{
		Read: resourceServerRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployment_href": &schema.Schema{
							Description: "ID of deployment the server is in",
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
						},
						"name": &schema.Schema{
							Description: "name of server",
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
						},
						"cloud_href": &schema.Schema{
							Description: "ID of cloud the server is in",
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"description": &schema.Schema{
				Description: "description of server",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"instance": &schema.Schema{
				Description: "server instance details",
				Type:        schema.TypeList,
				Elem:        resourceInstance(),
				Computed:    true,
			},
			"name": &schema.Schema{
				Description: "name of server",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"optimized": &schema.Schema{
				Description: "A flag indicating whether Instances of this Server should be optimized for high-performance volumes (e.g. Volumes supporting a specified number of IOPS). Not supported in all Clouds.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"links": &schema.Schema{
				Description: "Hrefs of related API resources",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeMap},
				Computed:    true,
			},
			"href": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	var loc *rsc.Locator
	loc = &rsc.Locator{Namespace: "rs_cm", Type: "servers"}
	res, err := client.List(loc, "servers", cmFilters(d))
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
