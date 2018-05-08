---
layout: "rightscale"
page_title: "Rightscale: route"
sidebar_current: "docs-rightscale-resource-route"
description: |-
  Create and maintain a RightScale route.
---

# rightscale_route

Use this resource to create, update or destroy RightScale [routes](http://reference.rightscale.com/api1.5/resources/ResourceRoutes.html).

## Example Usage

```hcl
resource "rightscale_route" "us-oregon-devops-vpc-route" {
  description = "A route to the internet through the internet gateway"
  destination_cidr_block = "0.0.0.0/0"
  next_hop_type = "network_gateway"
  next_hop_href = "${rightscale_network_gateway.my_network_gateway.href}"
  route_table_href = "${rightscale_route_table.my_route_table.href}"
}
```

## Argument Reference

The following arguments are supported:

* `route_table_href` - (Required) Href of route table in which to create new route.

* `destination_cidr_block` - (Required) Destination network in CIDR nodation.

* `next_hop_type` - (Required) The route next hop type.  Options are 'instance', 'network_interface', 'network_gateway', 'ip_string', and 'url'.

* `next_hop_href` - (Contextual) The href of the Route's next hop. Required if 'next_hop_type' is 'instance', 'network_interface', or 'network_gateway'.

* `next_hop_ip` - (Contextual) The IP Address of the Route's next hop. Required if 'next_hop_type' is 'ip_string'.

* `next_hop_url` - (Contextual) The URL of the Route's next hop. Required if 'next_hop_type' is 'url'.

* `description` - (Optional) Route description.

## Attributes Reference

The following attributes are exported:

* `href` - Href of the route. 

* `resource_uid` - Route resource_uid.

* `links` - Hrefs of related API resources.

* `created_at` - Created at datestamp.

* `updated_at` - Last updated at datestamp.
