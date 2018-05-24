---
layout: "rightscale"
page_title: "Rightscale: datacenter"
sidebar_current: "docs-rightscale-datasource-datacenter"
description: |-
  Defines a datacenter datasource to operate against.
---

# rightscale_datacenter

Use this data source to get the ID or other attributes of an existing datacenter for use in other resources.

Filter block is optional - ommitting it will result in the first available datacenter in a given cloud.

## Example Usage 1: Basic configuration of a datacenter data source

```hcl
 data "rightscale_datacenter" "ec2-us-east-1a" {
   cloud_href = "/api/clouds/1"
   filter {
     name = "us-east-1a"
   }
 }

output "datacenter name" {
  value = "${data.rightscale_datacenter.ec2-us-east-1a.name}"
}

output "datacenter ID" {
  value = "${data.rightscale_datacenter.ec2-us-east-1a.id}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` (REQUIRED) - The cloud_href the datacenter belongs to

* `filter` (OPTIONAL) - The filter block supports:

  * `name` - The name of the datacenter

  * `resource_id` - The resource_uid of the datacenter.  If this filter option is set, additional retry logic will be enabled to wait up to 5 minutes for cloud resources to be polled and populated for use.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the datacenter

* `description` - The description of the datacenter

* `resource_id` - The resource_uid of the datacenter as reported by the rightscale platform

* `links` - Hrefs of related API resources

