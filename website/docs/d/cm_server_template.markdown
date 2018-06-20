---
layout: "rightscale"
page_title: "Rightscale: server_template"
sidebar_current: "docs-rightscale-datasource-server_template"
description: |-
  Defines a server template datasource to operate against.
---

# rightscale_server_template

Use this data source to get the Href or other attributes of a server template in your account for use in other resources.

Filter block is optional - ommitting it will result in the first available server template in a given cloud.

## Example Usage 1: Basic configuration of a server template data source

```hcl
 data "rightscale_server_template" "mysql" {
   filter {
     name = "Database Manager for MySQL"
     revision = 24
   }
 }

output "server template name" {
  value = "${data.rightscale_server_template.mysql.name}"
}

output "server template href" {
  value = "${data.rightscale_server_template.mysql.href}"
}
```

## Argument Reference

The following arguments are supported:

* `filter` (OPTIONAL) - The filter block supports:

  * `name` - The name of the server template
  
  * `revision` - The revision of the server template, use 0 to match latest non-committed version

  * `description` - The description of the server template
  
  * `lineage` - The lineage of the server template
  
  * `multi_cloud_image_href` - The href of the server template multicloud image resource

## Attributes Reference

The following attributes are exported:

* `name` - The name of the server template

* `description` - The description of the server template

* `lineage` - The lineage of the server template
  
* `revision` - The revision of the server template, use 0 to match latest non-committed version

* `links` - Hrefs of related API resources

* `href` - Href of the server template

