---
layout: "rightscale"
page_title: "Rightscale: security_group_rule"
sidebar_current: "docs-rightscale-resource-security_group_rule"
description: |-
  Create and maintain a RightScale security_group_rule.
---

# rightscale_security_group_rule

Use this resource to create, update or destroy RightScale [security group rules](http://reference.rightscale.com/api1.5/resources/ResourceSecurityGroupRules.html).

## Example Usage - Create a security group rule

```hcl
resource "rightscale_security_group_rule" "allow-ssh-from-all" {
  security_group_href = "${rightscale_security_group.us-oregon-vpc-devops-security-group.href}"
  direction = "ingress"
  protocol = "tcp"
  source_type = "cidr_ips"
  cidr_ips = "0.0.0.0/0"
  protocol_details {
    start_port = "22"
    end_port = "22"
  }
}
```

## Argument Reference

The following arguments are supported:

* `source_type` - (Required) Source type. May be a CIDR block or another Security Group. Options are 'cidr_ips' or 'group'.

* `protocol` - (Required) Protocol to filter on.  Options are 'tcp', 'udp', 'icmp' and 'all'.

* `security_group_href` - (Required) Href of parent security group.

* `protocol_details` - (Required) Block options include:

  * `start_port` (Contextual) - Start of port range (inclusive). Required if protocol is 'tcp' or 'udp'.

  * `end_port` (Contextual) - End of port range (inclusive). Required if protocol is 'tcp' or 'udp'.

  * `icmp_code` (Contextual) - ICMP code. Required if protocol is 'icmp'.

  * `icmp_type` (Contextual) - ICMP type. Required if protocol is 'icmp'.

* `cidr_ips` - (Contextual) An IP address range in CIDR notation. Required if source_type is 'cidr'. Conflicts with 'group_name' and 'group_owner'.

* `group_name` - (Contextual) Name of source Security Group. Required if source_type is 'group'.  Conflicts with 'cidr_ips'.

* `group_owner` - (Contexual) Owner of source Security Group. Required if source_type is 'group'. Conflicts with 'cidr_ips'.

* `direction` - (Optional) Direction of traffic to apply rule against.  Options are 'ingress' or 'egress'.

* `priority` - (Optional) Lower takes precedence. Supported by security group rules created in Microsoft Azure only.

## Attributes Reference

The following attributes are exported:

* `href` - Href of the security group rule.

* `resource_uid` - Cloud resource_uid.

* `links` - Hrefs of related API resources.
