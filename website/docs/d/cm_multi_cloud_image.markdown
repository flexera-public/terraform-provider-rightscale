---
layout: "rightscale"
page_title: "Rightscale: multi_cloud_image"
sidebar_current: "docs-rightscale-datasource-multi_cloud_image"
description: |-
  Defines a multi cloud image datasource to operate against.
---

# rightscale_multi_cloud_image

Use this data source to get the ID or other attributes of an existing [multi cloud image](http://docs.rightscale.com/cm/dashboard/design/multicloud_images/) for use in other resources.

Filter block is optional - ommitting it will result in the first available multi cloud image in the account.

## Example Usage 1: Basic configuration of a multi cloud image data source

```hcl
data "rightscale_multi_cloud_image" "centos_64" {
   filter {
     name = "RightImage_CentOS_6.4_x64_v13.5"
     revision = 43
   }
 }

output "multi cloud image name" {
  value = "${data.rightscale_multi_cloud_image.centos_64.name}"
}

output "multi cloud image ID" {
  value = "${data.rightscale_multi_cloud_image.centos_64.id}"
}
```

## Argument Reference

The following arguments are supported:

* `server_template_href` (OPTIONAL) - The server_template_href the multi cloud image appears in

* `filter` (OPTIONAL) - The filter block supports:

  * `id` - The multi cloud image ID (e.g. /api/clouds/1/multi_cloud_images/4VODPN6TQ60RC)

  * `name` - The name of the multi cloud image

  * `description` - The description of the multi cloud image

  * `revision` - The revision of multi-cloud image, use 0 to match latest non-committed version


## Attributes Reference

The following attributes are exported:

* `name` - The name of the multi cloud image

* `description` - The description of the multi cloud image

* `revision` - The revision of multi-cloud image, use 0 to match latest non-committed version

* `links` - Hrefs of related API resources