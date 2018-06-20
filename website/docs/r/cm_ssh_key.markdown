---
layout: "rightscale"
page_title: "Rightscale: ssh_key"
sidebar_current: "docs-rightscale-resource-ssh_key"
description: |-
  Create and maintain an ssh key resource in a given cloud.
---

# rightscale_ssh_key

Use this resource to create, update or destroy RightScale [SSH keys](http://reference.rightscale.com/api1.5/resources/ResourceSshKeys.html).

## Example Usage

```hcl
resource "rightscale_ssh_key" "infra-ssh-key" {
  name = "infra"
  cloud_href = ${data.rightscale_cloud.ec2_us_east_1.href}
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` - (Required) The href of the cloud with the ssh key you want.

* `name` - (Required) SSH Key name.

## Attributes Reference

The following attributes are exported:

* `resource_uid` - Cloud resource_uid.

* `links` - Hrefs of related API resources.