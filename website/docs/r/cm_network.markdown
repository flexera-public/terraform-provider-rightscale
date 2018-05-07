---
layout: "rightscale"
page_title: "Rightscale: network"
sidebar_current: "docs-rightscale-resource-network"
description: |-
  Create and maintain a rightscale cloud management network.
---

# rightscale_network

Use this resource to create, update or destroy rightscale [networks](http://reference.rightscale.com/api1.5/resources/ResourceNetworks.html) in cloud management.

## Example Usage

```hcl
resource "rightscale_network" "us-oregon-devops-vpc" {
  name = "us-oregon-devops-vpc"
  description = "AWS US Oregon vpc for devopery"
  cloud_href = "/api/clouds/6"
  cidr_block = "192.168.0.0/16"
}

output "us-oregon-devops-vpc-aws-uid" {
  value = "${rightscale_network.us-oregon-devops-vpc.resource_uid}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` - (Required) Cloud you want to create the network in.

* `cidr_block` - (Optional*) Cloud-specific.  Some clouds require this field, others do not.

* `name` - (Optional) Network name.

* `description` - (Optional) Network description.

* `instance_tenancy` - (Optional) Launch policy for AWS instances in the network. Specify 'dedicated' to force all instances to be launched as 'dedicated'.  Defaults to 'default.'

* `route_table_href` - (Optional) Sets the default route table for this network, useful if you create the route table with a different resource.

* `deployment_href` - (Optional) HREF of the deployment that owns the network.  If you wish to use a deployment object as top level ownership construct, perhaps allocating the new network to a single deployment, then provide this href.

## Attributes Reference

The following attributes are exported:

* `href` - Href of the network.

* `resource_uid` - Cloud resource_uid as reported by cm platform.

* `links` - Hrefs of related API resources.