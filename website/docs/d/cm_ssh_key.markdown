---
layout: "rightscale"
page_title: "Rightscale: ssh_key"
sidebar_current: "docs-rightscale-datasource-ssh-key"
description: |-
  Defines an ssh key datasource to operate against.
---

# rightscale_ssh_key

Use this data source to locate and extract info about an existing [ssh_key](http://reference.rightscale.com/api1.5/resources/ResourceSshKeys.html) to pass to other rightscale resources.  Define the 'sensitive' view to access the private key material.

## Example Usage 1: Basic Usage

```hcl
data "rightscale_ssh_key" "infra-ssh-key" {
  filter {
    name = "infra"
  }
  cloud_href = "${data.rightscale_cloud.ec2_us_oregon.href}"
}

data "rightscale_cloud" "ec2_us_oregon" {
  filter {
    name = "EC2 us-west-2"
    cloud_type = "amazon"
  }
}
```

## Example Usage 2: Private key material from created resource

```hcl
resource "rightscale_ssh_key" "resource_ssh_key" {
  name = "rs-tf-ssh-key"
  cloud_href = "${data.rightscale_cloud.ec2_us_oregon.href}"
}

data "rightscale_ssh_key" "read_resource_ssh_key" {
  filter {
    name = "${rightscale_ssh_key.resource_ssh_key.name}"
  }
  cloud_href = "${data.rightscale_cloud.ec2_us_oregon.href}"
  view = "sensitive"
  depends_on = ["rightscale_ssh_key.resource_ssh_key"]
}

output "read-private-key-material" {
  value = "${data.rightscale_ssh_key.read_resource_ssh_key.material}"
}

data "rightscale_cloud" "ec2_us_oregon" {
  filter {
    name = "EC2 us-west-2"
    cloud_type = "amazon"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cloud_href` - (Required) The Href of the cloud with the ssh key you want.

* `view` - (Optional) Set this to 'sensitive' to request the api return 'sensitive' information (in this case the private key material) with the request. Assumes rs account privs sufficient to do this operation.

The `filter` block supports:

* `name` - (Optional) SSH key name.  Pattern match.

* `resource_uid` - (Optional) resource_uid of the SSH key.  If this filter option is set, additional retry logic will be enabled to wait up to 5 minutes for cloud resources to be polled and populated for use.

## Attributes Reference

The following attributes are exported:

* `name` - Official cloud name as displayed in cm platform.

* `resource_uid` - resource_uid of the SSH key.

* `links` - Hrefs of related API resources.

* `material` - (Contextual) Available only if 'sensitive' view is set.

* `href` - Href of the SSH key.