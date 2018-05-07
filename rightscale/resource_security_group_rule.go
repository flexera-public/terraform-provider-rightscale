package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// resource "rightscale_security_group_rule" "allow-ssh-from-all" {
//     security_group_href = "${rightscale_security_group.my_security_group.href}"
//	   direction = "ingress"
//     protocol = "tcp"
//	   source_type = "cidr_ips"
//     cidr_ips = "0.0.0.0/0"
//	   protocol_details {
//       start_port = "22"
//       end_port = "22"
//     }
// }

func resourceSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSecurityGroupRuleRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceSecurityGroupRuleCreate,

		Schema: map[string]*schema.Schema{
			"source_type": {
				Type:         schema.TypeString,
				Description:  "Source type. May be a CIDR block or another Security Group. Options are 'cidr_ips' or 'group'.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cidr_ips", "group"}, false),
			},
			"cidr_ips": {
				Type:          schema.TypeString,
				Description:   "An IP address range in CIDR notation. Required if source_type is 'cidr'. Conflicts with 'group_name' and 'group_owner'",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"group_name", "group_owner"},
			},
			"direction": {
				Type:         schema.TypeString,
				Description:  "Direction of traffic.  Options are 'ingress' or 'egress.'",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false),
			},
			"group_name": {
				Type:          schema.TypeString,
				Description:   "Name of source Security Group. Required if source_type is 'group'.  Conflicts with 'cidr_ips'.",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cidr_ips"},
			},
			"group_owner": {
				Type:          schema.TypeString,
				Description:   "Owner of source Security Group. Required if source_type is 'group'. Conflicts with 'cidr_ips'.",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cidr_ips"},
			},
			"priority": {
				Type:        schema.TypeInt,
				Description: "Lower takes precedence. Supported by security group rules created in Microsoft Azure only.",
				Optional:    true,
				ForceNew:    true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "Protocol to filter on.  Options are 'tcp', 'udp', 'icmp' and 'all'.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp", "all"}, false),
			},
			"protocol_details": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_port": {
							Type:        schema.TypeString,
							Description: "End of port range (inclusive). Required if protocol is 'tcp' or 'udp'.",
							Optional:    true,
							ForceNew:    true,
						},
						"icmp_code": {
							Type:        schema.TypeString,
							Description: "ICMP code. Required if protocol is 'icmp'.",
							Optional:    true,
							ForceNew:    true,
						},
						"icmp_type": {
							Type:        schema.TypeString,
							Description: "ICMP type. Required if protocol is 'icmp'.",
							Optional:    true,
							ForceNew:    true,
						},
						"start_port": {
							Type:        schema.TypeString,
							Description: "Start of port range (inclusive). Required if protocol is 'tcp' or 'udp'.",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"security_group_href": {
				Type:        schema.TypeString,
				Description: "Href of parent security group",
				Required:    true,
				ForceNew:    true,
			},

			// Read-only fields
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"href": {
				Type:        schema.TypeString,
				Description: "Href of security group rule",
				Computed:    true,
			},
		},
	}
}

func resourceSecurityGroupRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	res, err := client.Create("rs_cm", "security_group_rules", securityGroupRuleWriteFields(d))
	if err != nil {
		return err
	}
	for k, v := range res.Fields {
		// for some reason the api requires 'cidr_ips' for source_type, but returns 'cidr' in the response.
		if k == "source_type" && v == "cidr" {
			d.Set(k, "cidr_ips")
		} else {
			d.Set(k, v)
		}
	}
	// Sets 'href' which is rightscale href (for stitching together cm resources IN rightscale) without namespace.
	d.Set("href", res.Locator.Href)
	// Sets 'id' which allows terraform to locate the objects created which includes namespace.
	d.SetId(res.Locator.Namespace + ":" + res.Locator.Href)
	return nil
}

func resourceSecurityGroupRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	loc, err := locator(d)
	if err != nil {
		return err
	}
	res, err := client.Get(loc)
	if err != nil {
		return handleRSCError(d, err)
	}
	for k, v := range res.Fields {
		// for some reason the api requires 'cidr_ips' for source_type, but returns 'cidr' in the response.
		if k == "source_type" && v == "cidr" {
			d.Set(k, "cidr_ips")
		} else {
			d.Set(k, v)
		}
	}
	return nil
}

func securityGroupRuleWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{
		"protocol":            d.Get("protocol"),
		"security_group_href": d.Get("security_group_href"),
		"source_type":         d.Get("source_type"),
	}
	if i, ok := d.GetOk("protocol_details"); ok {
		fields["protocol_details"] = i.([]interface{})[0].(map[string]interface{})
	}
	for _, f := range []string{
		"action", "cidr_ips", "direction", "group_name", "group_owner", "priority",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return rsc.Fields{"security_group_href": d.Get("security_group_href"), "security_group_rule": fields}
}
