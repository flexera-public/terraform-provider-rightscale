package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// resource "rightscale_subnet" "aws-us-east-1d" {
//     name = "subnet-aws-us-east-1d"
//     cloud_href = "$[data.rightscale_cloud.us-east.href}"
//     network_href = "${rightscale_network.aws-us-east-devops-vpc.href}"
//     description = "Subnet for aws us-east-1d for aws-us-east-devops vpc"
//     cidr_block = "192.168.1.0/24"
//     datacenter_href = "${data.rightscale_datacenter.ec2_us_east_1d.id}"
// }

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceSubnetCreate,
		Update: resourceUpdateFunc(subnetWriteFields),

		Schema: map[string]*schema.Schema{
			"cidr_block": &schema.Schema{
				Description: "range of IP addresses for subnet, this parameter is required for Amazon clouds",
				Type:        schema.TypeString,
				Required:    true,
			},
			"datacenter_href": &schema.Schema{
				Description: "ID of subnet datacenter",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": &schema.Schema{
				Description: "description of subnet",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": &schema.Schema{
				Description: "name of subnet",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"network_href": &schema.Schema{
				Description: "ID of subnet network",
				Type:        schema.TypeString,
				Required:    true,
			},
			"route_table_href": &schema.Schema{
				Description: "ID of subnet route table",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cloud_href": {
				Type:        schema.TypeString,
				Description: "Href of cloud resource",
				Required:    true,
				ForceNew:    true,
			},

			// Read-only fields
			"is_default": &schema.Schema{
				Description: "whether the subnet is the network default subnet",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"resource_uid": &schema.Schema{
				Description: "cloud ID of subnet resource",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": {
				Description: "indicates whether subnet is pending, available etc.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"href": {
				Type:        schema.TypeString,
				Description: "href of subnet",
				Computed:    true,
			},
		},
	}
}

func resourceSubnetCreate(d *schema.ResourceData, m interface{}) error {
	var rt string
	{
		if r, ok := d.GetOk("route_table_href"); ok {
			rt = r.(string)
		}
	}

	client := m.(rsc.Client)

	// first create subnet with no default route table
	fields := subnetWriteFields(d)
	delete(fields, "route_table_href")
	res, err := client.Create("rs_cm", "subnets", fields)
	if err != nil {
		return err
	}
	for k, v := range res.Fields {
		d.Set(k, v)
	}

	// Sets 'href' which is rightscale href (for stitching together cm resources IN rightscale) without namespace.
	d.Set("href", res.Locator.Href)
	// Sets 'id' which allows terraform to locate the objects created which includes namespace.
	d.SetId(res.Locator.Namespace + ":" + res.Locator.Href)

	// then update with default route table if any
	if rt != "" {
		d.Set("route_table_href", rt)
		if err := resourceUpdateFunc(subnetWriteFields)(d, client); err != nil {
			// Attempt to delete previously created subnet, ignore errors
			client.Delete(res.Locator)
			return err
		}
	}
	return nil
}

func subnetWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{
		"cidr_block":   d.Get("cidr_block"),
		"network_href": d.Get("network_href"),
	}
	for _, f := range []string{
		"datacenter_href", "description", "name",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return rsc.Fields{"cloud_href": d.Get("cloud_href"), "subnet": fields}
}
