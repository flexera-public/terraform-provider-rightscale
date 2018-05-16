---
layout: "rightscale"
page_title: "Rightscale: instance"
sidebar_current: "docs-rightscale-resource-instance"
description: |-
  Create and maintain a RightScale instance.
---

# rightscale_instance

Use this resource to create, update or destroy an instance.

## Example Usage : Basic configuration of an instance resource

```hcl
resource "rightscale_instance" "crunis_instance" {
  cloud_href = "/api/clouds/6"
  image_href = "/api/clouds/6/images/3TRNL47PJB97N"
  instance_type_href = "/api/clouds/6/instance_types/8SCHNH0JBHE1R"
  deployment_href = "/api/deployments/934588004"
  name = "My Instance"
}
```

## Argument Reference

The following arguments are supported:

 * `name` - (Required) The name of the instance.

 * `cloud_href` - (Required) The cloud_href the instance belongs to.

 * `image_href` - (Required) The href of the instance image.

 * `instance_type_href` - (Required) The href of the instance type.

 * `server_template_href` - (Optional) The href of the instance server template resource.

 * `inputs` - (Optional) Inputs associated with an instance when incarnated from a [server](https://github.com/rightscale/terraform-provider-rightscale/blob/master/rightscale/website/docs/r/cm_server.markdown) or [server_array](https://github.com/rightscale/terraform-provider-rightscale/blob/master/rightscale/website/docs/r/cm_server_array.markdown).

 * `associate_public_ip_address` - (Optional) Indicates if the instance will get a Public IP address.

 * `cloud_specific_attributes` - (Optional) Attributes specific to the cloud the instance belongs to.

* `datacenter_href` - (Optional) The href of the datacenter that holds the instance (e.g. /api/clouds/6/datacenters/6IHONC8ANOUHI).

* `deployment_href` - (Optional) The href of the deployment that contains the instance (e.g. /api/deployments/594684003).

* `ip_forwarding_enabled` - (Optional) Allows this Instance to send and receive network traffic when the source and destination IP addresses do not match the IP address of this Instance.

* `kernel_image_href` - (Optional) The href of the instance kernel image.

* `ramdisk_image_href` - (Optional) The href of the instance ramdisk image.

* `security_group_hrefs` - (Optional) The href of the instance security groups.

* `placement_group_href` - (Optional) The href of the [placement_group](http://docs.rightscale.com/cm/dashboard/clouds/aws/ec2_placement_groups.html) that contains the instance (e.g. /api/placement_groups/512SV3FUJA7OO).

* `ssh_key_href` - (Optional) The href of the SSH key to use.

* `subnet_hrefs` - (Optional) The hrefs of the instance subnet.

* `user_data` - (Optional) User data that RightScale automatically passes to your instance at boot time.

* `locked` - (Optional)  Whether instance is locked, a locked instance cannot be terminated or deleted.
