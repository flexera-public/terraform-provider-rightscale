---
layout: "rightscale"
page_title: "Rightscale: route_table"
sidebar_current: "docs-rightscale-resource-route_table"
description: |-
  Create and maintain a RightScale route_table.
---

# rightscale_route_table

Use this resource to create, update or destroy RightScale [route tables](http://reference.rightscale.com/api1.5/resources/ResourceRouteTables.html).

## Example Usage

```hcl
resource "rightscale_route_table" "us-oregon-devops-vpc-route-table" {
  name = "us-oregon-devops-vpc-route-table"
  description = "AWS US Oregon vpc route table for devopery"
  cloud_href = "${data.rightscale_cloud.us-oregon.href}"
  network_href = "${rightscale_network.us-oregon-devops-vpc.href}
}

output "us-oregon-devops-vpc-route-table-aws-uid" {
  value = "${rightscale_network.us-oregon-devops-vpc-route-table.resource_uid}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` - (Required) Href of the cloud you want to create the route table in.

* `network_href` - (Required) Href of the network that owns the route table.

* `name` - (Required) Route table name.

* `description` - (Optional) Route table description.

## Attributes Reference

The following attributes are exported:

* `href` - Href of the route table.

* `resource_uid` - Cloud resource_uid.

* `links` - Hrefs of related API resources.

* `created_at` - Created at datestamp.

* `updated_at` - Last updated at datestamp.