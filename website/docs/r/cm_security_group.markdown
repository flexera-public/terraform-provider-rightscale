---
layout: "rightscale"
page_title: "Rightscale: security_group"
sidebar_current: "docs-rightscale-resource-security_group"
description: |-
  Create and maintain a rightscale cloud management security_group.
---

# rightscale_security_group

Use this resource to create, update or destroy rightscale [security groups](http://reference.rightscale.com/api1.5/resources/ResourceSecurityGroups.html).

## Example Usage - Create a security group

```hcl
resource "rightscale_security_group" "us-oregon-devops-vpc-security-group" {
  name = "us-oregon-devops-vpc-sg"
  description = "AWS US Oregon vpc security group for devopery"
  cloud_href = "/api/clouds/6"
  network_href = "${rightscale_network.us-oregon-devops-vpc.href}"
}

output "us-oregon-devops-vpc-sg-href" {
  value = "${rightscale_security_group.us-oregon-devops-vpc-security-group.href}"
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` - (Required) Cloud you want to create the security group in.

* `network_href` - (Required) Href of the network to create the security group in.

* `name` - (Required) Security group name.

* `description` - (Optional) Security group description.

* `deployment_href` - (Optional) Href of the deployment that owns the security group.  If you wish to use a deployment object as top level ownership construct, perhaps allocating the new security group to a single deployment, then provide this href.

## Attributes Reference

The following attributes are exported:

* `href` - Href of the security group.

* `resource_uid` - Cloud resource_uid.

* `links` - Hrefs of related API resources.
