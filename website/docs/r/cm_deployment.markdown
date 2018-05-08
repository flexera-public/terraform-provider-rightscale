---
layout: "rightscale"
page_title: "Rightscale: deployment"
sidebar_current: "docs-rightscale-resource-deployment"
description: |-
  Create and maintain a RightScale deployment.
---

# rightscale_deployment

Use this resource to create, update or destroy RightScale [deployments](http://docs.rightscale.com/cm/dashboard/manage/deployments/index.html) in cloud management.

## Example Usage

```hcl
resource "rightscale_deployment" "production_sydney_deployment" {
  name = "production_sydney"
  description = "Production Operations in Sydney for Red Team"
}

output "sydney_prod_deployment_href" {
  value = "${rightscale_deployment.production_sydney_deployment.href}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Deployment name.

* `description` - (Optional) Deployment description.

* `resource_group_href` - (Optional) Href of the Windows Azure Resource Group attached to the deployment.

* `locked` - (Optional) Set to true to lock the deployment.

* `server_tag_scope` - (Optional) Routing scope for tags for servers in the deployment.  Options are 'account' or 'deployment,' defaults to 'deployment.'

## Attributes Reference

The following attributes are exported:

* `href` - Href of the deployment.

* `links` - Hrefs of related API resources.