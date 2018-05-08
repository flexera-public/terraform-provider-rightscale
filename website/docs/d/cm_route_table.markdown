---
layout: "rightscale"
page_title: "Rightscale: route_table"
sidebar_current: "docs-rightscale-datasource-route_table"
description: |-
  Defines a route_table datasource to operate against.
---

# rightscale_route_table

Use this data source to locate and extract info about an existing [route table](http://reference.rightscale.com/api1.5/resources/ResourceRouteTables.html) to pass to other rightscale resources.

## Example Usage: Get existing route table resource_uid

```hcl
data "rightscale_route_table" "infrastructure-us-east-route-table" {
  filter {
    name = "Production Infrastructure US-East"
    network_href = "${data.rightscale_network.infrastructure-us-east.id}"
  }
}

output "prod-infra-us-east-route-table-aws-uid" {
  value = "${data.rightscale_route_table.infrastructure-us-east-route-table.resource_uid}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` (Optional) block supports:

  * `name` - (Optional) Route table name.  Pattern match.

  * `cloud_href` - (Optional) Cloud href of route table.

  * `network_href` - (Optional) Network href that owns the route table.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the route table.

* `resource_uid` - Cloud resource_uid.

* `description` - The description of the route table.

* `routes` - Associated routes.

* `links` - Hrefs of related API resources.