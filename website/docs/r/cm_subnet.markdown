---
layout: "rightscale"
page_title: "Rightscale: subnet"
sidebar_current: "docs-rightscale-resource-subnet"
description: |-
  Create and maintain a RightScale subnet.
---

# rightscale_subnet

Use this resource to create, update or destroy RightScale [subnets](http://reference.rightscale.com/api1.5/resources/ResourceSubnets.html).

## Example Usage

```hcl
resource "rightscale_subnet" "devops-oregon-subnet-a" {
  name = "devops-oregon-vpc-a"
  description = "AWS US Oregon Subnet for devopery in az 'a'"
  cloud_href = "${data.rightscale_cloud.aws-oregon.id}"
  datacenter_href = "${data.rightscale_datacenter.ec2_us_oregon_a.id}"
  network_href = "${rightscale_network.aws-oregon-devops-vpc.href}"
  cidr_block = "192.168.8.0/24"
}

output "us-oregon-devops-subnet-a-aws-href" {
  value = "${rightscale_network.devops-oregon-subnet-a.href}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` - (Required) Href of cloud you want to create the subnet in.

* `network_href` - (Required) Href of network to create subnet in.

* `cidr_block` - (Required) Subnet allocation range in CIDR notation.

* `name` - (Optional) Subnet name.

* `description` - (Optional) Subnet description.

* `datacenter_href` - (Optional) Href of cloud datacenter to assign subnet to.

* `route_table_href` - (Optional) Sets the default route table for this subnet, useful if you create the route table with a different resource.

## Attributes Reference

The following attributes are exported:

* `href` - Href of the subnet.

* `resource_uid` - Cloud resource_uid.

* `is_default` - Indicates whether the subnet is the network default subnet. (true or false)

* `state` - Indicates whether subnet is pending, available etc.

* `links` - Hrefs of related API resources.
