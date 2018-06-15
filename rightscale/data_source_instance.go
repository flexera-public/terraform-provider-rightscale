package rightscale

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Examples:
//
// data "rightscale_instance" "mysql" {
//   filter {
//     resource_uid = "vpc-c31ee987"
//   }
//   cloud_href = "${data.rightscale_cloud.ec2_us_east_1.id}"
// }
//
// data "rightscale_instance" "worker_2" {
//   filter {
//     name = "Worker #2"
//   }
//   server_array = "${data.rightscale_server_array.workers}"
// }

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Read: resourceInstanceRead,

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
							Description: "name of instance, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"datacenter_href": {
							Type:        schema.TypeString,
							Description: "ID of the instance datacenter resource",
							Optional:    true,
							ForceNew:    true,
						},
						"os_platform": {
							Type:         schema.TypeString,
							Description:  "OS platform of instance",
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"windows", "linux"}, false),
						},
						"parent_href": {
							Type:        schema.TypeString,
							Description: "ID of instance server or server array parent resource",
							Optional:    true,
							ForceNew:    true,
						},
						"server_template_href": {
							Type:        schema.TypeString,
							Description: "ID of instance server template resource",
							Optional:    true,
							ForceNew:    true,
						},
						"state": {
							Type:        schema.TypeString,
							Description: "state of instance",
							Optional:    true,
							ForceNew:    true,
						},
						"placement_group_href": {
							Type:        schema.TypeString,
							Description: "ID of instance placement group resource",
							Optional:    true,
							ForceNew:    true,
						},
						"public_dns_name": {
							Type:        schema.TypeString,
							Description: "Public DNS name of instance",
							Optional:    true,
							ForceNew:    true,
						},
						"private_dns_name": {
							Type:        schema.TypeString,
							Description: "Private DNS name of instance",
							Optional:    true,
							ForceNew:    true,
						},
						"public_ip": {
							Type:        schema.TypeString,
							Description: "Public IP of instance",
							Optional:    true,
							ForceNew:    true,
						},
						"private_ip": {
							Type:        schema.TypeString,
							Description: "Private IP instance",
							Optional:    true,
							ForceNew:    true,
						},
						"deployment_href": {
							Type:        schema.TypeString,
							Description: "ID of deployment resource that owns instance",
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
			"associate_public_ip_address": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cloud_href": {
				Type:          schema.TypeString,
				Description:   "ID of instance cloud resource, exclusive with 'server_array_href'",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"server_array_href"},
			},
			"cloud_specific_attributes": instanceCloudAttributes,
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pricing_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_addresses": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"public_ip_addresses": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_array_href": {
				Type:          schema.TypeString,
				Description:   "ID of instance server array resource, exclusive with 'cloud_href'",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cloud_href"},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceInstanceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	var loc *rsc.Locator
	if cloud, ok := d.GetOk("cloud_href"); ok {
		loc = &rsc.Locator{Namespace: "rs_cm", Href: cloud.(string)}
	} else if array, ok := d.GetOk("server_array_href"); ok {
		loc = &rsc.Locator{Namespace: "rs_cm", Href: array.(string)}
	} else {
		return fmt.Errorf("instance data source must specify one of 'cloud_href' or 'server_array_href'")
	}

	// if 'resource_uid' filter is set, we expect it to show up.
	// retry for 5 min to allow rightscale to poll cloud to discover.
	if uidset := cmUIDSet(d); uidset {
		timeout := 300
		err := cmIndexRetry(client, loc, "instances", d, timeout)
		if err != nil {
			return err
		}
	}

	res, err := client.List(loc, "instances", cmFilters(d))
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
	d.Set("href", res[0].Locator.Href)
	return nil
}
