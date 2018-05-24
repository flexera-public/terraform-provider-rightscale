---
layout: "rightscale"
page_title: "Rightscale: volume_snapshot"
sidebar_current: "docs-rightscale-datasource-volume_snapshot"
description: |-
  Defines a volume snapshot datasource to operate against.
---

# rightscale_volume_snapshot

Use this data source to get the ID or other attributes of an existing volume snapshot for use in other resources.

Filter block is optional - ommitting it will result in the first available volume snapshot in a given cloud.

## Example Usage 1: Basic configuration of a volume snapshot data source

```hcl
data "rightscale_volume_snapshot" "mysql_master" {
   filter {
     name = "mysql_master"
   }
   cloud_href = "/api/clouds/1"
 }

output "snapshot name" {
  value = "${data.rightscale_volume_snapshot.mysql_master.name}"
}

output "snapshot ID" {
  value = "${data.rightscale_volume_snapshot.mysql_master.id}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` (REQUIRED) - The cloud_href the volume snapshot belongs to

* `filter` (OPTIONAL) - The filter block supports:

  * `id` - The volume snapshot ID (e.g. /api/clouds/1/volume_snapshots/4VODPN6TQ60RC)

  * `name` - The name of the volume snapshot

  * `description` - The description of the volume snapshot

  * `state` - The state of the volume snapshot (e.g.: available, pending, ...)

  * `parent_volume_href` - The ID of the parent resource

  * `resource_uid` - The resource_uid of the volume snapshot.  If this filter option is set, additional retry logic will be enabled to wait up to 5 minutes for cloud resources to be polled and populated for use.

  * `deployment_href` - The href of the [deployment](http://docs.rightscale.com/cm/dashboard/manage/deployments/) that contains the volume snapshot (e.g. /api/deployments/594684003)

## Attributes Reference

The following attributes are exported:

* `description` - The description of the volume snapshot

* `name` - The name of the volume snapshot

* `size` - The size of the volume snapshot

* `state` - The state of the volume snapshot (e.g.: available, pending, ...)

* `resource_uid` - The resource_uid of the volume snapshot (e.g. /api/clouds/1/volume_snapshots/4VODPN6TQ60RC)

* `links` - Hrefs of related API resources

* `created_ at` - Time of creation of the volume snapshot

* `updated_at` - Last update of the volume snapshot
