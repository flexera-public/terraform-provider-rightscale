---
layout: "rightscale"
page_title: "Rightscale: credential"
sidebar_current: "docs-rightscale-datasource-credential"
description: |-
  Defines a credential datasource to operate against. 
---

# rightscale_credential

Use this data source to locate and extract the data of an existing credential to pass to other rightscale resources, or to access the values.  Viewing values of credentials assumes requisite account permission levels.

## Example Usage: Access credential value

```hcl
data "rightscale_credential" "account_aws_access_key_id" {
  filter {
    name = "AWS_ACCESS_KEY_ID"
  }
}

output "my-aws-access-key-id" {
  value = "${data.rightscale_credential.account_aws_access_key_id.value}"
}
```

## Argument Reference

The following arguments are supported:

* `view` - (Optional) Set this to 'default' to NOT request credential value with api response.  This allows use of existing credential with other rightscale provider resources (extracting href and handing to other resources). Offereed in case user lacks rs account privs sufficient to view credential values. 

The `filter` block supports:

* `name` - (Optional) Credential name.  Pattern match. 

* `description` - (Optional) Description of credential.  Pattern match.

## Attributes Reference

The following attributes are exported:

* `name` - Name of the credential.

* `description` - Description of the credential.

* `value` - (Contextual) Available unless if 'default' view is set.  Value of the credential.

* `links` - Hrefs of related API resources.

* `created_at` - Datestamp of credential creation.

* `updated_at` - Datestamp of when credential was updated last.

* `href` - Href of the credential.