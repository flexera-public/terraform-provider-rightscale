---
layout: "rightscale"
page_title: "Rightscale: server"
sidebar_current: "docs-rightscale-datasource-server"
description: |-
  Create and maintain a RightScale server 
---

# rightscale_server

Use this resource to create, update or destroy a RightScale server 

## Example Usage : Basic configuration of a server resource

```hcl
resource "rightscale_server" "web_server" {
  name = "web_server"
  deployment_href = "/api/deployments/1234"
  instance {
    cloud_href = "/api/clouds/1234"
    image_href = "/api/clouds/1234/images/1234"
    instance_type_href = "/api/clouds/1234/instance_types/1234"
    name = "web_instance"
    server_template_href = "/api/server_templates/1234"
  }
}
```

## Argument Reference

The following arguments are supported:

* `deployment_href` - (Required) The href of the deployment

* `description` - (Optional) A description of the server

* `instance` - (Required) See [rightscale_instance](https://github.com/rightscale/terraform-provider-rightscale/blob/master/rightscale/website/docs/r/cm_server.markdown)

* `name` - (Required) The name of the server

* `optimized` - (Optional) A flag indicating whether Instances of this Server should be optimized for high-performance volumes

* `cloud_href` - (Required) The ID of the cloud with the ssh key you want

## Attributes Reference

The following attributes are exported:

* `links` - Hrefs of related API resources
