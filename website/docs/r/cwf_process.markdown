---
layout: "rightscale"
page_title: "Rightscale: CWF Process"
sidebar_current: "docs-rightscale-resource-cwf_process"
description: |-
  Create and maintain a rightscale CloudWorkFlow process.
---

# rightscale_cwf_process

Use this resource to create or destroy RightScale [CloudWorkFlow processes](http://docs.rightscale.com/ss/reference/rcl/).

Creating the CWF process will run it synchronously and return the output values (if any).

Destroying a CWF process will remove it from the CWF Console table of executed processes. It WON'T undo any action that it could have executed nor would stop any running process.

It is NOT possible to update a CWF process.

## Example Usage

This example CWF process will look for all servers whose name starts by "db-slave-" and will execute the specified RightScript on them,
returning the number of servers that have been reached.

```hcl
resource "rightscale_cwf_process" "run_executable_by_prefix" {

  parameters = [
     { "kind" = "string"
       "value" = "db-slave-" },
     { "kind" = "string"
       "value" = "/api/right_scripts/1018361003" }
     ]

  source = <<EOF
define main($server_prefix, $rightscript_href) return $servers_reached do
  @servers = rs_cm.servers.get(filter: ["name==" + $server_prefix])
  @servers.current_instance().run_executable(right_script_href: $rightscript_href)
  $servers_reached = size(@servers)
end
EOF

}

output "cwf_status" {
  value = "${rightscale_cwf_process.run_executable_by_prefix.status}"
}

output "cwf_servers_updated" {
  value = "${rightscale_cwf_process.run_executable_by_prefix.outputs["$servers_reached"]}"
}
```

## Argument Reference

The following arguments are supported:

* `source` - (Required) Source code to be executed, written in [RCL (RightScale CloudWorkFlow Language)](http://docs.rightscale.com/ss/reference/rcl/v2/index.html). The function should always be called main. Example:
```hcl
  source = <<EOF
define main($a, $b) return $result do
  $result = $a + $b
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