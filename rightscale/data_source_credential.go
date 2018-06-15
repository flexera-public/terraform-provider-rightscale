package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_credential" "account_aws_access_key_id" {
//   filter {
//     name = "AWS_ACCESS_KEY_ID"
//   }
// }

func dataSourceCredential() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCredentialRead,

		Schema: map[string]*schema.Schema{
			"view": {
				Type:         schema.TypeString,
				Description:  "Filter at api level for the view: 'default' or 'sensitive' are valid options",
				Optional:     true,
				Default:      "sensitive",
				ValidateFunc: validation.StringInSlice([]string{"default", "sensitive"}, false),
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
							Description: "name of credential, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "description of credential",
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"updated_at": {
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

func datasourceCredentialRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	acParams := make(map[string]string)
	acParams["view"] = d.Get("view").(string)

	loc := &rsc.Locator{Namespace: "rs_cm", Type: "credential", ActionParams: acParams}

	res, err := client.List(loc, "credential", cmFilters(d))
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
