---
layout: "rightscale"
page_title: "Rightscale: image"
sidebar_current: "docs-rightscale-datasource-image"
description: |-
  Defines an image datasource to operate against.
---

# rightscale_image

Use this data source to get the ID (rightscale href) of a registered image in a specific cloud for use in other resources.  Sets default filter scope to own account, but allows for public searching if specified in filter block.

Beware that searching a very popular cloud (say aws us-east) based on name with 'visibility = "public"' is gonna be slow...

## Example Usage #1 - Finding specific AMI in own account based on resource_uid

```hcl
data "rightscale_cloud" "ec2_us_oregon" {
  filter {
    name = "EC2 us-west-2"
    cloud_type = "amazon"
  }
}

data "rightscale_image" "my_sweet_ami" {
  cloud_href = "${data.rightscale_cloud.ec2_us_oregon.id}"
  filter {
    resource_uid = "ami-abcdefg"
  }
}
...
```

## Example Usage #2 - Finding public image in cloud based on filters on name, description, etc.

Warning: The more images a cloud has public, the longer this filter call will take.  Consider multiple filters to narrow the scope.

```hcl
data "rightscale_cloud" "ec2_us_oregon" {
  filter {
    name = "EC2 us-west-2"
    cloud_type = "amazon"
  }
}

data "rightscale_image" "my_sweet_ami" {
  cloud_href = "${data.rightscale_cloud.ec2_us_oregon.id}"
  visibility = "public"
  filter {
    name = "My Super Great AMI"
    os_platform = "linux"
    description = "AMI Image built from CI that does great things"
  }
}
...
```

## Argument Reference

**Note - omitting the filter block IS valid and will return the first private image object available in the specific cloud and your account.  Probably not what you really want.**

The following arguments are supported:

* `cloud_href` - (Required) The ID of the cloud with the image you want.

* `filter` - (Optional) block supports:

  * `visibility` (Optional) Image visibility as displayed in cm platform.  Options are "private" or "public."  Defaults to "private."

  * `resource_uid` (Optional) Image resource_uid as displayed in cm platform.

  * `name` - (Optional) Image name as displayed in cm platform.  Pattern match.

  * `description` - (Optional) Image description as displayed in cm platform.  Pattern match.

  * `image_type` - (Optional) Image type as referenced in cm platform. This will be either "machine", "machine_azure", "ramdisk" or "kernel".

  * `os_platform` - (Optional) Image OS platform as referenced in cm platform.  This will either be "windows" or "linux."

  * `cpu_architecture` - (Optional) Image CPU architecture as referenced in cm platform.  Generally "x64_64", etc.  Pattern match.

## Attributes Reference

The following attributes are exported:

* `visibility` - Image visibility as displayed in cm platform.

* `resource_uid` - Image unique resource identifier as displayed in cm platform.

* `name` - Image name as displayed in cm platform.

* `description` - Image description as displayed in cm platform.

* `cpu_architecture` - Image CPU architecture as referenced in cm platform.

* `os_platform` - Image OS platform as referenced in cm platform.

* `root_device_storage` - Image root device storage as reported in cm platform.  Eg "volume" vs "instance", etc.

* `image_type` - Image type as referenced in cm platform.

* `virtualization_type` - Image virtualization type as referenced in cm platform. Eg "hvm" etc.

* `links` - Hrefs of related API resources.