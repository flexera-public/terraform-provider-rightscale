---
layout: "rightscale"
page_title: "Rightscale: instance_type"
sidebar_current: "docs-rightscale-datasource-instance-type"
description: |-
  Defines a instance_type datasource to operate against.
---

# rightscale_instance_type

Use this data source to get the ID (rightscale href) of an instance type (eg "m4.large" vs "n1-standard" vs "DSv2") in a specific cloud for use in other resources.

## Example Usage - Get href for instance type "m4.large" in aws us-oregon cloud

```hcl
data "rightscale_cloud" "ec2_us_oregon" {
  filter {
    name = "EC2 us-west-2"
    cloud_type = "amazon"
  }
}

data "rightscale_instance_type" "m4_large" {
  cloud_href = "${data.rightscale_cloud.ec2_us_oregon.id}"
  filter {
    name = "m4.large"
  }
}
...
```

## Argument Reference

**Note - omitting the filter block IS valid and will return the first object available in the specific cloud and your account.  Probably not what you really want.**

The following arguments are supported:

* `cloud_href` - (Required) The ID of the cloud with the instance type you want.

* `filter` - (Optional) block supports:

  * `resource_uid` (Optional) Instance type resource uid.  If this filter option is set, additional retry logic will be enabled to wait up to 5 minutes for cloud resources to be polled and populated for use.

  * `name` - (Optional) Instance type name as displayed in cm platform.  Pattern match.

  * `description` - (Optional) Instance type description as displayed in cm platform.  Pattern match.

  * `cpu_architecture` - (Optional) Instance type CPU architecture as referenced in cm platform.  Generally "x64_64", etc.  Pattern match.

## Attributes Reference

The following attributes are exported:

* `resource_uid` - Instance type unique resource identifier as displayed in cm platform.

* `name` - Instance type name as displayed in cm platform.

* `description` - Instance type description as displayed in cm platform.

* `cpu_architecture` - Instance type CPU architecture as displayed in cm platform.

* `cpu_count` - Instance type CPU count as displayed in cm platform.

* `cpu_speed` - Instance type CPU speed as displayed in cm platform.

* `memory` - Instance type memory as displayed in cm platform.

* `links` - Hrefs of related API resources.