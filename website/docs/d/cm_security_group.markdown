---
layout: "rightscale"
page_title: "Rightscale: security_group"
sidebar_current: "docs-rightscale-datasource-security_group"
description: |-
  Defines a security_group datasource to operate against.
---

# rightscale_security_group

Use this data source to locate and extract info about an existing [security group](http://reference.rightscale.com/api1.5/resources/ResourceSecurityGroups.html) to pass to other rightscale resources.

## Example Usage: Get existing security group resource_uid

```hcl
data "rightscale_security_group" "infrastructure-us-east-security-group" {
  cloud_href = "${data.rightscale_cloud.infrastructure-aws-us-east.id}"
  filter {
    name = "Infrastructure SG"
    network_href = "${data.rightscale_network.infrastructure-us-east.id}"
  }
}

output "prod-infra-us-east-aws-sg-uid" {
  value = "${data.rightscale_security_group.infrastructure-us-east-security-group.resource_uid}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` (Required) Cloud href that the security group exists in.

* `filter` (Optional) block supports:

  * `name` - (Optional) Security group name.  Pattern match.

  * `resource_uid` - (Optional) Cloud resource uid for security group.  If this filter option is set, additional retry logic will be enabled to wait up to 5 minutes for cloud resources to be polled and populated for use.

  * `network_href` - (Optional) Network href that security group is created in.

  * `deployment_href` - (Optional) Href of the deployment that owns the security group.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the security group.

* `resource_uid` - The cloud resource_uid of the security group.

* `description` - The description of the security group.

* `links` - Hrefs of related API resources.
