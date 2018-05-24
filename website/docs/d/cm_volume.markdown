---
layout: "rightscale"
page_title: "Rightscale: volume"
sidebar_current: "docs-rightscale-datasource-volume"
description: |-
  Defines a volume datasource to operate against.
---

# rightscale_volume

Use this data source to get the ID or other attributes  of an existing volume for use in other resources.

## Example Usage 1: Basic configuration of a volume data source

```hcl
data "rightscale_volume" "a_volume" {
  cloud_href = "/api/clouds/1"

  filter {
    name = "my_volume"
  }
}

output "volume name" {
  value = "${data.rightscale_volume.a_volume.name}"
}

output "volume ID" {
  value = "${data.rightscale_volume.a_volume.id}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` (REQUIRED) - The cloud_href the volume belongs to

* `filter` (OPTIONAL) - The filter block supports:

  * `id` - The volume ID (e.g. /api/clouds/1/volumes/63NFHKF8B7RP4)

  * `name` - The name of the volume

  * `description` - The description of the volume

  * `resource_uid` - The resource_uid of the volume.  If this filter option is set, additional retry logic will be enabled to wait up to 5 minutes for cloud resources to be polled and populated for use.

  * `deployment_href` - The href of the [deployment](http://docs.rightscale.com/cm/dashboard/manage/deployments/) that contains the volume (e.g. /api/deployments/594684003)

  * `datacenter_href` - The href of the [datacenter](http://docs.rightscale.com/cm/dashboard/clouds/generic/datacenter_zones_concepts.html) that holds the volume (e.g. /api/clouds/6/datacenters/6IHONC8ANOUHI)

  * `parent_volume_snapshot_href` - The href of snapshot the volume was created of

## Attributes Reference

The following attributes are exported:

* `id` - The volume ID (e.g. /api/clouds/1/volumes/63NFHKF8B7RP4)

* `name` - The name of the volume

* `description` - The description of the volume

* `resource_uid` - The resource_uid of the volume (e.g. vol-045e33fd28a746c45)

* `links` - Hrefs of related API resources

* `size` - The volume size (in GB)

* `status` - The volume Status (e.g. available, in-use, ...)

* `updated_at` - Last update of the volume
