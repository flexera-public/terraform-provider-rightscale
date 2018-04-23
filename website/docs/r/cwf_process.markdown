---
layout: "rightscale"
page_title: "Rightscale: CWF Process"
sidebar_current: "docs-rightscale-resource-cwf_process"
description: |-
  Create and maintain a rightscale CloudWorkFlow process.
---

# rightscale_cwf_process

Use this resource to create, update or destroy rightscale [CloudWorkFlow processes](http://docs.rightscale.com/ss/reference/rcl/).

## Example Usage

```hcl
resource "rightscale_cwf_process" "a_process" {
  source = <<EOF
define main() return $total do
  $a = 156.5534
  $b = 42421000
  $total = $a + $b
  end
EOF
}

output "process_status" {
  value = "${rightscale_cwf_process.a_process.status}"
}

output "process_total" {
  value = "${rightscale_cwf_process.a_process.outputs["$total"]}"
}

output "process_outputs" {
  value = "${rightscale_cwf_process.a_process.outputs}"
}
```

You can also use the output of a module as a `source`, this way, you can have your library of CWF scripts stored in a module and then call the one you want to execute:

```hcl
module "lib_cwf" {
  # Local filesystem library
  source = "../cwf_lib"
  # ... or Github repository:
  # source = "github.com/org/cwf_lib"
}

resource "rightscale_cwf_process" "a_process" {
  source = "${module.my_cwf_lib.taskX}"
}

output "process_total" {
  value = "${rightscale_cwf_process.a_process.outputs["$total"]}"
}
```

## Argument Reference

The following arguments are supported:

* `source` - (Required) Source code to be executed, written in [RCL (RightScale CloudWorkFlow Language)](http://docs.rightscale.com/ss/reference/rcl/v2/index.html).

## Attributes Reference

The following attributes are exported:

* `status` - Process status, one of "completed", "failed", "canceled" or "aborted".

* `error` - Process execution error if any.

* `outputs` - Process outputs if any. This is a TypeMap, one particular output can be accessed via `outputs["$var"]` (see example above)