---
layout: "rightscale"
page_title: "Rightscale: volume_type"
sidebar_current: "docs-rightscale-datasource-volume-type"
description: |-
  Defines a volume type datasource to operate against.
---

# rightscale_volume_type

Use this data source to get the ID or other attributes of an existing volume type (as defined by a given cloud) for use in other resources.

Filter block is optional - ommitting it will result in the first available volume_type in a given cloud.

## Example Usage: Basic configuration of a volume type data source

```hcl
data "rightscale_volume_type" "aws_us_east_ebs_gp2" {
  cloud_href = "/api/clouds/1"

  filter {
    name = "gp2"
  }
}

```

## Argument Reference

The following arguments are supported:

* `cloud_href` (REQUIRED) - The cloud_href the volume type belongs to

* `filter` (OPTIONAL) - The filter block supports:

  * `name` - The name of the volume type as reported by the rightscale platform

  * `resource_uid` - The resource_uid of the volume_type.  If this filter option is set, additional retry logic will be enabled to wait up to 5 minutes for cloud resources to be polled and populated for use.

## Attributes Reference

The following attributes are exported:

* `id` - The volume type ID (e.g. /api/clouds/1/volume_types/B37A8VOCJIODH)

* `name` - The name of the volume type.

* `description` - The description of the volume type.

* `resource_uid` - The resource_uid of the volume type. (e.g. gp2)

* `links` - Hrefs of related API resources

* `size` - The volume size (in GB) if applicable (depends on cloud)

* `created_at` - Creation date of the volume type

* `updated_at` - Last update of the volume type

* `href` - Href of the volume type