package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_ssh_key" "ssh" {
//   filter {
//     name = "infra"
//   }
//   cloud_href = ${data.rightscale_cloud.ec2_us_east_1.id}
//   # 'sensitive' view returns private key material with api call; assumes rs account privs sufficient to do so.
//	 view = "sensitive"
// }

func dataSourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: datasourceSSHKeyRead,

		Schema: map[string]*schema.Schema{
			"cloud_href": {
				Type:        schema.TypeString,
				Description: "ID of the SSH key cloud",
				Required:    true,
				ForceNew:    true,
			},
			"view": {
				Type:         schema.TypeString,
				Description:  "Filter at api level for the view: 'default' or 'sensitive' are valid options",
				Optional:     true,
				Default:      "default",
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
							Description: "name of SSH key, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"resource_uid": {
							Type:        schema.TypeString,
							Description: "cloud ID - if this filter is set additional retry logic will fire to allow for cloud resource discovery",
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"material": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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

func datasourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	cloud := d.Get("cloud_href").(string)
	acParams := make(map[string]string)
	acParams["view"] = d.Get("view").(string)

	loc := &rsc.Locator{Namespace: "rs_cm", Href: cloud, ActionParams: acParams}

	// if 'resource_uid' filter is set, we expect it to show up.
	// retry for 5 min to allow rightscale to poll cloud to discover.
	if uidset := cmUIDSet(d); uidset {
		timeout := 300
		err := cmIndexRetry(client, loc, "ssh_keys", d, timeout)
		if err != nil {
			return err
		}
	}

	res, err := client.List(loc, "ssh_keys", cmFilters(d))
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
