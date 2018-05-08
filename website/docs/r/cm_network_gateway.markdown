---
layout: "rightscale"
page_title: "Rightscale: network_gateway"
sidebar_current: "docs-rightscale-resource-network_gateway"
description: |-
  Create and maintain a RightScale network gateway.
---

# rightscale_network_gateway

Use this resource to create, update or destroy RightScale [network gateways](http://reference.rightscale.com/api1.5/resources/ResourceNetworkGateways.html) in cloud management.

## Example Usage #1 - Create an internet gateway

```hcl
resource "rightscale_network_gateway" "us-oregon-devops-vpc-gateway" {
  name = "us-oregon-devops-vpc-gateway"
  description = "AWS US Oregon vpc gateway for devopery"
  cloud_href = "/api/clouds/6"
  type = "internet"
}

output "us-oregon-devops-vpc-gateway-aws-uid" {
  value = "${rightscale_network_gateway.us-oregon-devops-vpc-gateway.resource_uid}"
}
```

## Example Usage #2 - Create an internet gateway and attach it to a network

```hcl
resource "rightscale_network_gateway" "us-oregon-devops-vpc-gateway" {
  name = "us-oregon-devops-vpc-gateway"
  description = "AWS US Oregon vpc gateway for devopery"
  cloud_href = "/api/clouds/6"
  type = "internet"
  network_href = "${rightscale_network.us-oregon-devops-vpc.href}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` - (Required) Cloud you want to create the network gateway in.

* `type` - (Required) Type of network gateway.  Options are "internet" or "vpc".

* `name` - (Required) Network gateway name.

* `description` - (Optional) Network gateway description.

* `network_href` - (Optional) Href of network you want to attach the network gateway to.

## Attributes Reference

The following attributes are exported:

* `href` - Href of the network gateway.

* `created_at` - Date the network gateway was created at.

* `updated_at` - Date the network gateway was updated at.

* `state` - State of the network gateway.  ("available" means attached to a network)

* `resource_uid` - Cloud resource_uid.

* `links` - Hrefs of related API resources.