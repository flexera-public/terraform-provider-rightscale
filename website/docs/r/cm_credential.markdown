---
layout: "rightscale"
page_title: "Rightscale: credential"
sidebar_current: "docs-rightscale-resource-credential"
description: |-
  Create and maintain a RightScale credential resource.
---

# rightscale_credential

Use this resource to create, read, update or destroy RightScale credential objects.

## Example Usage

```hcl
resource "rightscale_credential" "database_password" {
  name = "DATABASE_PASSWORD"
  value = "rightscale11"
  description = "Top Secret database password"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the credential.

* `value` - (Required) Value of the credential.

* `description` - (Optional) Description of the credential.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the credential.

* `description` - Description of the credential.

* `value` - Value of the credential. 

* `links` - Hrefs of related API resources.

* `created_at` - Datestamp of credential creation.

* `updated_at` - Datestamp of when credential was updated last.