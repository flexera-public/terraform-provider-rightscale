package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_deployment" "development_shared_resources_deployment" {
//   filter {
//     name = "dev assets do not touch"
//   }
// }

func dataSourceDeployment() *schema.Resource {
	return &schema.Resource{
		Read: datasourceDeploymentRead,

		Schema: map[string]*schema.Schema{
			"view": {
				Type:         schema.TypeString,
				Description:  "Filter at api level for the view: 'default,' 'inputs' or 'inputs_2_0' are valid options",
				Optional:     true,
				Default:      "default",
				ValidateFunc: validation.StringInSlice([]string{"default", "inputs", "inputs_2_0"}, false),
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
							Description: "name of deployment, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "description of deployment, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"resource_group_href": {
							Type:        schema.TypeString,
							Description: "resource group href to filter on",
							Optional:    true,
							ForceNew:    true,
						},
						"server_tag_scope": {
							Type:        schema.TypeString,
							Description: "tag routing scope of deployments to filter on, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},

			// Read-only fields
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_tag_scope": {
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

func datasourceDeploymentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	acParams := make(map[string]string)
	acParams["view"] = d.Get("view").(string)

	loc := &rsc.Locator{Namespace: "rs_cm", Type: "deployments", ActionParams: acParams}

	res, err := client.List(loc, "deployments", cmFilters(d))
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
