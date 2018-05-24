---
layout: "rightscale"
page_title: "Rightscale: instance"
sidebar_current: "docs-rightscale-datasource-instance"
description: |-
  Defines an instance datasource to operate against.
---

# rightscale_instance

Use this data source to get the ID or other attributes of an existing instance for use in other resources.

Filter block is optional - ommitting it will result in the first available instance in a given cloud.

## Example Usage 1: Basic configuration of a instance data source

```hcl
data "rightscale_instance" "an_instance" {
  cloud_href = "/api/clouds/1"

  filter {
    name = "my_instance"
  }
}

output "instance name" {
  value = "${data.rightscale_instance.an_instance.name}"
}

output "instance ID" {
  value = "${data.rightscale_instance.an_instance.id}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` (REQUIRED unless server_array_href specified) - The cloud_href the instance belongs to (mutually exclusive with server_array_href, specify only one)

* `server_array_href` (REQUIRED unless cloud_href specified) - The server_array_href the instance belongs to (mutually exclusive with cloud_href, specify only one)

* `filter` (OPTIONAL) - The filter block supports:

  * `id` - The instance ID (e.g. /api/clouds/1/instances/63NFHKF8B7RP4)

  * `name` - The name of the instance

  * `state` - The state of the instance (e.g.: operational, terminated, stranded, ...)

  * `os_platform` - The OS platform of the instance. One of "linux" or "windows"

  * `parent_href` - The ID of instance server or server array parent resource.

  * `server_template_href` - The ID of the instance server template resource

  * `public_dns_name` - The public DNS name of the instance

  * `private_dns_name` - The private DNS name of the instance

  * `public_ip` - The public IP of the instance

  * `private_ip` - The private IP of the instance

  * `resource_uid` - The resource_uid of the instance.  If this filter option is set, additional retry logic will be enabled to wait up to 5 minutes for cloud resources to be polled and populated for use.

  * `deployment_href` - The href of the [deployment](http://docs.rightscale.com/cm/dashboard/manage/deployments/) that contains the instance (e.g. /api/deployments/594684003)

  * `placement_group_href` - The href of the [placement_group](http://docs.rightscale.com/cm/dashboard/clouds/aws/ec2_placement_groups.html) that contains the instance (e.g. /api/placement_groups/512SV3FUJA7OO)

  * `datacenter_href` - The href of the [datacenter](http://docs.rightscale.com/cm/dashboard/clouds/generic/datacenter_zones_concepts.html) that holds the instance (e.g. /api/clouds/6/datacenters/6IHONC8ANOUHI)

## Attributes Reference

The following attributes are exported:

* `associate_public_ip_address` - Indicates if the instance will get a Public IP address

* `cloud_href` - The cloud_href the instance belongs to (mutually exclusive with server_array_href)

* `server_array_href` - The server_array_href the instance belongs to (mutually exclusive with cloud_href)

* `cloud_specific_attributes` - Attributes specific to the cloud the instance belongs to

* `id` - The instance ID (e.g. /api/clouds/1/instances/63NFHKF8B7RP4)

* `name` - The name of the instance

* `pricing_type` - Pricing type of the instance (e.g. fixed, spot)

* `resource_uid` - The resource_uid of the instance (e.g. e0bf62bc-4e35-11e8-9f1f-0242ac110003)

* `links` - Hrefs of related API resources

* `locked` - Whether instance is locked, a locked instance cannot be terminated or deleted

* `private_ip_addresses` - List of private IP addresses of the instance

* `public_ip_addresses` - List of public IP addresses of the instance

* `state` - The instance state (e.g. operational, terminated, stranded, ...)

* `created_ at` - Time of creation of the instance

* `updated_at` - Last update of the instance
