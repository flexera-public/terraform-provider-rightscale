---
layout: "rightscale"
page_title: "Rightscale: server_array"
sidebar_current: "docs-rightscale-datasource-server_array"
description: |-
  Create and maintain a RightScale server_array
---

# rightscale_server_array

Use this resource to create, update or destroy a RightScale server_array

## Example Usage : Basic configuration of a server_array resource

```hcl
resource "rightscale_server_array" "frontend_servers_array" {
	array_type = "alert"

	datacenter_policy = [{
		datacenter_href = "/api/clouds/1234/datacenters/DEOLL9UREJ7TA"
		max             = 4
		weight          = 100
	}]

	elasticity_params = {
		alert_specific_params = {
		decision_threshold = 75
		}

		bounds = {
		min_count = 1
		max_count = 4
		}

		pacing = {
		resize_down_by = 1
		resize_up_by   = 1
		}
	}

	instance = {
		cloud_href           = "/api/clouds/1234"
		image_href           = "/api/clouds/1234/images/1234"
		instance_type_href   = "/api/clouds/1234/instance_types/1234"
		server_template_href = "/api/server_templates/1234"
		name                 = "Frontend"
		subnet_hrefs         = ["/api/clouds/1/subnets/52NUHI2B8LVH1"]
	}

	name            = "FrontEnd Servers Array"
	state           = "enabled"
	deployment_href = "/api/deployments/1234"
	}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the server_array

* `description` - (Optional) Description of the server_array

* `state` - (Required) he status of the server array. If enabled, the server array is enabled for scaling actions. One of "enabled" or "disabled"

* `deployment_href` - (Required) ID of deployment in which to create server_array

* `array_type` - (Required) The type of server_array. One of "alert" or "queue"

* `optimized` - (Optional) A flag indicating whether Instances of this ServerArray should be optimized for high-performance volumes (e.g. Volumes supporting a specified number of IOPS). Not supported in all Clouds.

* `datacenter_policy` - (Required) This is an array of datacenter policies. Each one must contain:

  * `datacenter_href` - (Required) The href of the server_array's datacenter / zone.

  * `max` - (Required) Maximum numbers of servers that can be allocated in this datacenter (0 for unlimited).

  * `weight` - (Required) Instance allocation (should total 100% accross datacenter_policies).

* `elasticity_params` - (Required) 

  * `bounds` - (Required)

    * `min_count` - (Required) The minimum number of servers that must be operational at all times in the server array.

    * `max_count` - (Required) The maximum number of servers that can be operational at the same time in the server array.

  * `pacing` - (Required)

      * `resize_down_by` - (Required) The number of servers to scale down by.

      * `resize_up_by` - (Required) The number of servers to scale up by.

      * `resize_calm_time` - (Optional) The time (in minutes) on how long you want to wait before you repeat another action.

  * `alert_specific_params` - (Required if alert array_type specified)

    * `decision_threshold` - (Required) The percentage of servers that must agree in order to trigger an alert before an action is taken.

    * `voters_tag_predicate` - (Optional) The Voters Tag that RightScale will use in order to determine when to scale up/down.

  * `queue_specific_params` - (Required if queue alert_type specified)

    * `collect_audit_entries` - (Optional) The audit SQS queue that will store audit entries.

    * `item_age` - (Required)

      * `algorithm` - (Optional) The algorithm that defines how an item's age will be determined, either by the average age or max (oldest) age.

      * `max_age` - (Optional) The threshold (in seconds) before a resize action occurs on the server array.

      * `regexp` - (Optional) The regexp that helps the system determine an item's \"age\" in the queue. Example: created_at: (\\d\\d\\d\\d-\\d\\d-\\d\\d \\d\\d:\\d\\d:\\d\\d UTC)

      * `queue_size` - (Required) Defines the ratio of worker instances per items in the queue. Example: If there are 50 items in the queue and \"Items per instance\" is set to 10, the server array will resize to 5 worker instances (50/10). Default = 1

* `schedule` - (Optional)

  * `day` - (Required) Specifies the day when an alert-based array resizes. One of "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday".

  * `max_count` - (Required) The maximum number of servers that must be operational at all times in the server array. NOTE: Any changes that are made to the min/max count in the server array schedule will overwrite the array's default min/max count settings.

  * `min_count` - (Required) The minimum number of servers that must be operational at all times in the server array. NOTE: Any changes that are made to the min/max count in the server array schedule will overwrite the array's default min/max count settings.

  * `time` - (Required) Specifies the time when an alert-based array resizes.

* `instance` - (Required) See [rightscale_instance](https://github.com/rightscale/terraform-provider-rightscale/blob/master/rightscale/website/docs/r/cm_instance.markdown)


## Attribute Reference

* `links` - Hrefs of related API resources