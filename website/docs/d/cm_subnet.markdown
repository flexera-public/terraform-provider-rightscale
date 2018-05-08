---
layout: "rightscale"
page_title: "Rightscale: subnet"
sidebar_current: "docs-rightscale-datasource-subnet"
description: |-
  Defines a subnet datasource to operate against.
---

# rightscale_subnet

Use this data source to locate and extract info about an existing [subnet](http://reference.rightscale.com/api1.5/resources/ResourceSubnets.html) to pass to other rightscale resources.

## Example Usage: Get existing subnet resource_uid

```hcl
data "rightscale_subnet" "infrastructure-aws-us-east-subnet-b" {
  cloud_href = "/api/clouds/1"
  filter {
    name = "Production Infrastructure Subnet US-East B"
  }
}

output "prod-infra-us-east-subnet-b-aws-uid" {
  value = "${data.rightscale_subnet.infrastructure-aws-us-east-subnet-b.resource_uid}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` (Required) Cloud the subnet exists in.

* `filter` (Optional) block supports:

  * `name` - (Optional) Subnet name.  Pattern match.

  * `network_href` - (Optional) Network href the the subnet exists in.

  * `resource_uid` - (Optional) The resource_uid of the subnet.

  * `datacenter_href` - (Optional) Href of the subnet datacenter resource.

  * `instance_href` - (Optional) Href of instance resource attached to subnet.

  * `visibility` - (Optional) Visibility of the subnet to filter by (private, shared, etc).

## Attributes Reference

The following attributes are exported:

* `name` - Name of the subnet.

* `resource_uid` - Subnet resource_uid.

* `cidr_block` - Subnet allocation range in CIDR notation.

* `is_default` - Reports if subnet is 'default' for a given subnet.

* `description` - The description of the subnet.

* `state` - Indicates whether subnet is pending, available etc.

* `visibility` - Visibility of the subnet.

* `links` - Hrefs of related API resources.