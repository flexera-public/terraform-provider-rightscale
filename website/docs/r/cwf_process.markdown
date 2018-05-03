---
layout: "rightscale"
page_title: "Rightscale: CWF Process"
sidebar_current: "docs-rightscale-resource-cwf_process"
description: |-
  Create and maintain a rightscale CloudWorkFlow process.
---

# rightscale_cwf_process

Use this resource to create or destroy RightScale [CloudWorkFlow processes](http://docs.rightscale.com/ss/reference/rcl/).

Creating the CWF process runs it synchronously and returns the output values (if any). If the CWF process fails, the Terraform script fails too.

Destroying the resource deletes the corresponding CWF process. Destroying a running process causes it to end in error.

It is NOT possible to update a CWF process.

## Example Usage

This example CWF process looks for all servers whose names start with "db-slave-" and executes the specified RightScript on them,
returning the number of servers that have been affected.

```hcl
resource "rightscale_cwf_process" "run_executable_by_prefix" {

  parameters = [
     { "kind" = "string"
       "value" = "db-slave-" },
     { "kind" = "string"
       "value" = "/api/right_scripts/1018361003" }
     ]

  source = <<EOF
define main($instance_prefix, $rightscript_href) return $instances_affected do
  @instances = rs_cm.instances.get(filter: ["name==" + $instance_prefix, "state==operational"])
  @instances.run_executable(right_script_href: $rightscript_href)
  $instances_affected = size(@instances)
end
EOF

}

output "cwf_status" {
  value = "${rightscale_cwf_process.run_executable_by_prefix.status}"
}

output "cwf_servers_affected" {
  value = "${rightscale_cwf_process.run_executable_by_prefix.outputs["$instances_affected"]}"
}
```

## Argument Reference

The following arguments are supported:

* `source` - (Required) Source code to be executed, written in [RCL (RightScale CloudWorkFlow Language)](http://docs.rightscale.com/ss/reference/rcl/v2/index.html). Several functions can be defined but the entry function should be called `main`. Example:
```hcl
  source = <<EOF
define adder($n1, $n2) return $res do
  $res = $n1 + $n2
end
define main($a, $b) return $result do
  call adder($a, $b) retrieve $tmp
  $result = "The total is " + $tmp
end
EOF
```

* `parameters` - Parameters for the RCL function. It consists of an array of values corresponding to the values being passed to the function defined in the "source" field in order of declaration. The values are defined as string maps with the "kind" and "value" keys. "kind" contains the type of the value being passed, could be one of "array", "boolean", "collection", "datetime", "declaration", "null", "number", "object", "string". The "value" key contains the value. For example:
```hcl
  parameters = [
     { "kind" = "string"
       "value" = "db-slave-" },
     { "kind" = "number"
       "value" = "42" }
     ]
```

## Attributes Reference

The following attributes are exported:

* `status` - Process status, one of "completed", "failed", "canceled" or "aborted".

* `error` - Process execution error if any.

* `outputs` - Process outputs if any. This is a TypeMap, one particular output can be accessed via `outputs["$var"]`, see "Example Usage" section.