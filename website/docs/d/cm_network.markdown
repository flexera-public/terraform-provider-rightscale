---
layout: "rightscale"
page_title: "Rightscale: network"
sidebar_current: "docs-rightscale-datasource-network"
description: |-
  Defines a network datasource to operate against.
---

# rightscale_network

Use this data source to locate and extract info about an existing [network](http://reference.rightscale.com/api1.5/resources/ResourceNetworks.html) to pass to other rightscale resources.

## Example Usage: Get existing network resource_uid

```hcl
data "rightscale_network" "infrastructure-us-east" {
  filter {
    name = "Production Infrastructure US-East"
  }
}

output "prod-infra-us-east-aws-uid" {
  value = "${data.rightscale_network.infrastructure-us-east.resource_uid}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` (Optional) block supports:

  * `name` - (Optional) Network name.  Pattern match.

  * `cloud_href` - (Optional) Cloud Href of network.

  * `deployment_href` - (Optional) Deployment href associated with network.

  * `cidr_block` - (Optional) CIDR notation block of network.

  * `resource_uid` - (Optional) The resource_uid of the network.  If this filter option is set, additional retry logic will be enabled to wait up to 5 minutes for cloud resources to be polled and populated for use.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the network.

* `resource_uid` - Network resource_uid as reported by cm platform.

* `cidr_block` - Network CIDR notation block of network.

* `instance_tenancy` - Tenancy of instances on network.

* `is_default` - Reports if network is 'default' for a given cloud.

* `description` - The description of the network.

* `links` - Hrefs of related API resources.

* `href` - Href of the network.