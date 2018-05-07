package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// resource "rightscale_network_gateway" "my_network_gateway" {
//   name        = "aws-us-oregon-dev-vpc-gateway"
//   description = "Development vpc internet network gateway in aws us-oregon"
//   cloud_href = "/api/clouds/6"
//	 type = "internet"
//   network_href = "${rightscale_network.my_network.href}"
//
// }

func resourceNetworkGateway() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: networkGatewayDelete,
		Create: resourceNetworkGatewayCreate,
		Update: resourceUpdateFunc(networkGatewayWriteFields),

		Schema: map[string]*schema.Schema{
			"cloud_href": &schema.Schema{
				Description: "ID of cloud in which to create network gateway",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": &schema.Schema{
				Description: "description of network gateway",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": &schema.Schema{
				Description: "name of network gateway",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": &schema.Schema{
				Description:  "type of network gateway",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"internet", "vpc"}, false),
			},
			"network_href": &schema.Schema{
				Description: "network href to attach to",
				Type:        schema.TypeString,
				Optional:    true,
			},

			// Read-only fields
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"href": {
				Type:        schema.TypeString,
				Description: "href of network gateway",
				Computed:    true,
			},
		},
	}
}

func resourceNetworkGatewayCreate(d *schema.ResourceData, m interface{}) error {
	var network string
	{
		if r, ok := d.GetOk("network_href"); ok {
			network = r.(string)
		}
	}

	client := m.(rsc.Client)

	// create network initially with no specific network attachment
	fields := networkGatewayWriteFields(d)
	delete(fields, "network_href")
	res, err := client.Create("rs_cm", "network_gateways", fields)
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

	// now update object with network_href if provided
	if network != "" {
		d.Set("network_href", network)
		if err := resourceUpdateFunc(networkGatewayWriteFields)(d, client); err != nil {
			// On error, delete the locator so the overall operation will fail
			client.Delete(res.Locator)
			return err
		}
	}
	return nil
}

func networkGatewayDelete(d *schema.ResourceData, m interface{}) error {
	// Network Gateway might be attached to a network per declaration.
	// Update object first to disconnect network gateway from network, or deletes will fail on dependency 4xx.
	var network string
	{
		if r, ok := d.GetOk("network_href"); ok {
			network = r.(string)
		}
	}

	client := m.(rsc.Client)
	loc, err := locator(d)
	if err != nil {
		return err
	}

	if network != "" {
		// set network_href to empty string so api triggers disconnection from network
		d.Set("network_href", "")
		if err := resourceUpdateFunc(networkGatewayWriteFields)(d, client); err != nil {
			return err
		}
	}

	// now delete the network gateway
	return client.Delete(loc)
}

func networkGatewayWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{
		"cloud_href": d.Get("cloud_href"),
		"name":       d.Get("name"),
		"type":       d.Get("type"),
	}
	if desc, ok := d.GetOk("description"); ok {
		fields["description"] = desc
	}
	if network, ok := d.GetOk("network_href"); ok {
		fields["network_href"] = network
	}
	return rsc.Fields{"network_gateway": fields}
}
