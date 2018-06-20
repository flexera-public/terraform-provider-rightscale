---
layout: "rightscale"
page_title: "Rightscale: server"
sidebar_current: "docs-rightscale-datasource-server"
description: |-
  Defines an server datasource to operate against. 
---

# rightscale_server

Use this data source to locate and extract info about an existing [server](http://reference.rightscale.com/api1.5/resources/ResourceServers.html) to pass to other rightscale resources.

## Example Usage 1: Basic configuration of a server data source

```hcl
data "rightscale_server" "web_server" {
  filter {
    name = "web"
  }
}
```
## Example Usage 2: Security group using a server's name

```hcl
data "rightscale_server" "web_server" {
  filter {
    name = "web"
  }
}

resource "rightscale_security_group" "sg_web_out" {
  name = "${data.rigthscale_server.web_server.name}"
  cloud_href = "/api/clouds/1234"
  deployment_href = "/api/deployments/1234
  description = "Web server security group"
  network_href = "/api/clouds/1234/networks/1234"
}
```

## Argument Reference

The following arguments are supported:

The `filter` block supports:

* `deployment_href` - The href of the deployment

* `name` - The name of the server

* `cloud_href` - The Href of the cloud with the ssh key you want

## Attributes Reference

The following attributes are exported:

* `description` - A description of the server

* `instance` - See [rightscale_instance](https://github.com/rightscale/terraform-provider-rightscale/blob/master/rightscale/website/docs/r/cm_server.markdown)

* `optimized` - A flag indicating whether instances of this server should be optimized for high-performance volumes

* `links` - Hrefs of related API resources

* `href` - Href of the server
