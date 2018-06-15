---
layout: "rightscale"
page_title: "Rightscale: network_gateway"
sidebar_current: "docs-rightscale-datasource-network_gateway"
description: |-
  Defines a network gateway datasource to operate against.
---

# rightscale_network_gateway

Use this data source to locate and extract info about an existing [network gateway](http://reference.rightscale.com/api1.5/resources/ResourceNetworkGateways.html) to pass to other rightscale resources.

## Example Usage: Get existing network gateway resource_uid

```hcl
data "rightscale_network_gateway" "infrastructure-us-east" {
  filter {
    name = "Production Infrastructure US-East"
  }
}

output "prod-infra-us-east-aws-uid" {
  value = "${data.rightscale_network_gateway.infrastructure-us-east.resource_uid}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` (Optional) block supports:

  * `name` - (Optional) Network gateway name.  Pattern match.

  * `cloud_href` - (Optional) Cloud href of network gateway.

  * `network_href` - (Optional) Network href that network gateway is attached to.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the network gateway.

* `resource_uid` - Network gateway resource_uid from cloud.

* `type` - Type of network gateway.  Options are "internet" or "vpc."

* `state` - State of the network gateway.  ("available" means attached to a network)

* `description` - The description of the network gateway.

* `links` - Hrefs of related API resources.

* `href` - Href of the network gateway.