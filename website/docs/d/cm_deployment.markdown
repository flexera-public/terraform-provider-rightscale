---
layout: "rightscale"
page_title: "Rightscale: deployment"
sidebar_current: "docs-rightscale-datasource-deployment"
description: |-
  Defines a deployment datasource to operate against. 
---

# rightscale_deployment

Use this data source to locate and extract info about an existing [deployment](http://docs.rightscale.com/cm/dashboard/manage/deployments/index.html) to pass to other rightscale resources.

## Example Usage: Get existing deployment href

```hcl
data "rightscale_deployment" "infrastructure" {
  filter {
    name = "Production Infrastructure US-East"
  }
}

output "prod-infra-us-east-href" {
  value = "${data.rightscale_deployment.infrastructure.id}"
}
```

## Argument Reference

The following arguments are supported:

* `view` - (Optional) Options include 'default,' 'inputs' or 'inputs_2_0.'  Defaults to 'default.'  Please see RightScale documentation for inputs for details on these different views. 

* `filter` - (Optional) Filter block to find matching deployment.

The `filter` block supports:

* `name` - (Optional) Credential name.  Pattern match. 

* `description` - (Optional) Description of credential.  Pattern match.

* `resource_group_href` - (Optional) Resource group href to filter on.

* `server_tag_scope` - (Optional) Tag routing scope to filter on.  Pattern match.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the credential.

* `description` - Description of the credential.

* `links` - Hrefs of related API resources.

* `locked` - Displays if the deployment is locked or not.

* `server_tag_scope` - Displays what the scope of tags are in the deployment. Options are "deployment" or "account."