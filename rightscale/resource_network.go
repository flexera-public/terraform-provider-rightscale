package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// resource "rightscale_network" "my_network" {
//   name        = "aws-us-oregon-dev-vpc"
//   description = "Development vpc network in aws us-oregon"
//   cloud_href = "/api/clouds/6"
//	 cidr_block = "192.168.0.0/16"
//
// }

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceNetworkCreate,
		Update: resourceUpdateFunc(networkWriteFields),

		Schema: map[string]*schema.Schema{
			"cidr_block": &schema.Schema{
				Description: "range of IP addresses for network, this parameter is required for Amazon clouds",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"cloud_href": &schema.Schema{
				Description: "ID of cloud to create network in",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"deployment_href": &schema.Schema{
				Description: "Optional href of deployment that owns the network.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": &schema.Schema{
				Description: "description of network",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance_tenancy": &schema.Schema{
				Description:  "launch policy for AWS instances in the network. Specify 'dedicated' to force all instances to be launched as 'dedicated'.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "default",
				ValidateFunc: validation.StringInSlice([]string{"default", "dedicated"}, false),
			},
			"name": &schema.Schema{
				Description: "name of network",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"route_table_href": &schema.Schema{
				Description: "sets the default route table for this network",
				Type:        schema.TypeString,
				Optional:    true,
			},

			// Read-onyl fields
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"href": {
				Type:        schema.TypeString,
				Description: "href of network",
				Computed:    true,
			},
		},
	}
}

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	var rt string
	{
		if r, ok := d.GetOk("route_table_href"); ok {
			rt = r.(string)
		}
	}

	client := m.(rsc.Client)

	// first create network with no default route table
	fields := networkWriteFields(d)
	delete(fields, "route_table_href")
	res, err := client.Create("rs_cm", "networks", fields)
	if err != nil {
		return err
	}
	for k, v := range res.Fields {
		d.Set(k, v)
	}

	// then update with default route table if any
	if rt != "" {
		d.Set("route_table_href", rt)
		if err := resourceUpdateFunc(networkWriteFields)(d, client); err != nil {
			// Attempt to delete previously created network, ignore errors
			client.Delete(res.Locator)
			return err
		}
	}

	// set ID last so Terraform does not assume the network has been
	// created until all operations have completed successfully.
	d.SetId(res.Locator.Namespace + ":" + res.Locator.Href)
	d.Set("href", res.Locator.Href)
	return nil
}

func networkWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{"cloud_href": d.Get("cloud_href")}
	for _, f := range []string{
		"cidr_block", "deployment_href", "description",
		"instance_tenancy", "name", "route_table_href",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return rsc.Fields{"network": fields}
}
