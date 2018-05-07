package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

func resourceServerArray() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceCreateFunc("rs_cm", "server_arrays", serverArrayWriteFields),
		Update: resourceUpdateFunc(serverArrayWriteFields),

		Schema: map[string]*schema.Schema{
			"array_type": &schema.Schema{
				Description:  "The array type for the Server Array.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"alert", "queue"}, false),
			},
			"datacenter_policy": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_href": &schema.Schema{
							Description: "The href of the Datacenter / Zone.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"max": &schema.Schema{
							Description: "Max instances (0 for unlimited).",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"weight": &schema.Schema{
							Description: "Instance allocation (should total 100%).",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
			},
			"deployment_href": &schema.Schema{
				Description: "ID of deployment in which to create server_array",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": &schema.Schema{
				Description: "description of server_array",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"elasticity_params": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_specific_params": &schema.Schema{
							Type:        schema.TypeList,
							Description: "Alert based server array params, required if 'array_type' is 'alert'",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"decision_threshold": &schema.Schema{
										Description: "The percentage of servers that must agree in order to trigger an alert before an action is taken.",
										Type:        schema.TypeInt,
										Optional:    true,
									},
									"voters_tag_predicate": &schema.Schema{
										Description: "The Voters Tag that RightScale will use in order to determine when to scale up/down.",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
						"bounds": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_count": &schema.Schema{
										Description: "The maximum number of servers that can be operational at the same time in the server array.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"min_count": &schema.Schema{
										Description: "The minimum number of servers that must be operational at all times in the server array.",
										Type:        schema.TypeInt,
										Required:    true,
									},
								},
							},
						},
						"pacing": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resize_calm_time": &schema.Schema{
										Description: "The time (in minutes) on how long you want to wait before you repeat another action.",
										Type:        schema.TypeInt,
										Optional:    true,
									},
									"resize_down_by": &schema.Schema{
										Description: "The number of servers to scale down by.",
										Type:        schema.TypeInt,
										Optional:    true,
									},
									"resize_up_by": &schema.Schema{
										Description: "The number of servers to scale up by.",
										Type:        schema.TypeInt,
										Optional:    true,
									},
								},
							},
						},
						"queue_specific_params": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"collect_audit_entries": &schema.Schema{
										Description: "The audit SQS queue that will store audit entries.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"item_age": &schema.Schema{
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"algorithm": &schema.Schema{
													Description:  "The algorithm that defines how an item's age will be determined, either by the average age or max (oldest) age.",
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"max_10", "avg_10"}, false),
												},
												"max_age": &schema.Schema{
													Description: "The threshold (in seconds) before a resize action occurs on the server array.",
													Type:        schema.TypeInt,
													Optional:    true,
												},
												"regexp": &schema.Schema{
													Description: "The regexp that helps the system determine an item's \"age\" in the queue. Example: created_at: (\\d\\d\\d\\d-\\d\\d-\\d\\d \\d\\d:\\d\\d:\\d\\d UTC)",
													Type:        schema.TypeString,
													Optional:    true,
												},
											},
										},
									},
									"queue_size": &schema.Schema{
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"items_per_instance": &schema.Schema{
													Description: "Defines the ratio of worker instances per items in the queue. Example: If there are 50 items in the queue and \"Items per instance\" is set to 10, the server array will resize to 5 worker instances (50/10). Default = 1",
													Type:        schema.TypeInt,
													Optional:    true,
												},
											},
										},
									},
								},
							},
						},
						"schedule": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": &schema.Schema{
										Description:  "Specifies the day when an alert-based array resizes.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, false),
									},
									"max_count": &schema.Schema{
										Description: "The maximum number of servers that must be operational at all times in the server array. NOTE: Any changes that are made to the min/max count in the server array schedule will overwrite the array's default min/max count settings.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"min_count": &schema.Schema{
										Description: "The minimum number of servers that must be operational at all times in the server array. NOTE: Any changes that are made to the min/max count in the server array schedule will overwrite the array's default min/max count settings.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"time": &schema.Schema{
										Description: "Specifies the time when an alert-based array resizes.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"instance": &schema.Schema{
				Description: "server array instance details",
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        resourceInstance(),
			},
			"name": &schema.Schema{
				Description: "name of server array",
				Type:        schema.TypeString,
				Required:    true,
			},
			"optimized": &schema.Schema{
				Description: "A flag indicating whether Instances of this ServerArray should be optimized for high-performance volumes (e.g. Volumes supporting a specified number of IOPS). Not supported in all Clouds.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"state": &schema.Schema{
				Description:  "The status of the server array. If active, the server array is enabled for scaling actions.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
			},

			// Read-only fields
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
		},
	}
}

func serverArrayWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{
		"array_type": d.Get("array_type"),
		"name":       d.Get("name"),
		"state":      d.Get("state"),
	}
	if i, ok := d.GetOk("instance"); ok {
		fields["instance"] = instanceWriteFieldsFromMap(i.([]interface{})[0].(map[string]interface{}))
	}
	if dp, ok := d.GetOk("datacenter_policy"); ok {
		dfs := make([]rsc.Fields, len(dp.([]interface{})))
		for i, df := range dp.([]interface{}) {
			dfs[i] = datacenterPolicyWriteFields(df.(map[string]interface{}))
		}
		fields["datacenter_policy"] = dfs
	}
	if dp, ok := d.GetOk("elasticity_params"); ok {
		fields["elasticity_params"] = elasticityParamsWriteFields(dp.([]interface{})[0].(map[string]interface{}))
	}

	for _, f := range []string{
		"deployment_href", "description", "optimized",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return rsc.Fields{"server_array": fields}
}

func datacenterPolicyWriteFields(d map[string]interface{}) rsc.Fields {
	fields := rsc.Fields{}
	for _, f := range []string{
		"datacenter_href", "max", "weight",
	} {
		fields[f] = d[f]
	}
	return fields
}

func elasticityParamsWriteFields(d map[string]interface{}) rsc.Fields {
	fields := rsc.Fields{}
	if ap, ok := d["alert_specific_params"]; ok && len(ap.([]interface{})) > 0 {
		fields["alert_specific_params"] = alertSpecificParamsWriteFields(ap.([]interface{})[0].(map[string]interface{}))
	}
	if b, ok := d["bounds"]; ok && len(b.([]interface{})) > 0 {
		if b.([]interface{})[0] != nil {
			fields["bounds"] = boundsWriteFields(b.([]interface{})[0].(map[string]interface{}))
		} else {
			fields["bounds"] = rsc.Fields{}
		}
	}
	if p, ok := d["pacing"]; ok && len(p.([]interface{})) > 0 {
		if p.([]interface{})[0] != nil {
			fields["pacing"] = pacingWriteFields(p.([]interface{})[0].(map[string]interface{}))
		} else {
			fields["pacing"] = rsc.Fields{}
		}
	}
	if q, ok := d["queue_specific_params"]; ok && len(q.([]interface{})) > 0 {
		fields["queue_specific_params"] = queueSpecificParamsWriteFields(q.([]interface{})[0].(*schema.ResourceData))
	}
	if s, ok := d["schedule"]; ok && len(s.([]interface{})) > 0 {
		sf := make([]rsc.Fields, len(s.([]interface{})))
		for i, sched := range s.([]interface{}) {
			sf[i] = scheduleWriteFields(sched.(map[string]interface{}))
		}
		fields["schedule"] = sf
	}
	return fields
}

func alertSpecificParamsWriteFields(d map[string]interface{}) rsc.Fields {
	fields := rsc.Fields{}
	for _, f := range []string{
		"decision_threshold", "voters_tag_predicate",
	} {
		if v, ok := d[f]; ok {
			fields[f] = v
		}
	}
	return fields
}

func boundsWriteFields(d map[string]interface{}) rsc.Fields {
	fields := rsc.Fields{}
	for _, f := range []string{
		"min_count", "max_count",
	} {
		if v, ok := d[f]; ok {
			fields[f] = v
		}
	}
	return fields
}

func pacingWriteFields(d map[string]interface{}) rsc.Fields {
	fields := rsc.Fields{}
	for _, f := range []string{
		"resize_calm_time", "resize_down_by", "resize_up_by",
	} {
		if v, ok := d[f]; ok {
			fields[f] = v
		}
	}
	return fields
}

func queueSpecificParamsWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{}
	if v, ok := d.GetOk("collect_audit_entries"); ok {
		fields["collect_audit_entries"] = v
	}
	if a, ok := d.GetOk("item_age"); ok {
		fields["item_age"] = itemAgeWriteFields(a.([]interface{})[0].(*schema.ResourceData))
	}
	if a, ok := d.GetOk("queue_size"); ok {
		fields["queue_size"] = queueSizeWriteFields(a.([]interface{})[0].(*schema.ResourceData))
	}
	return fields
}

func itemAgeWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{}
	for _, f := range []string{
		"algorithm", "max_age", "regexp",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return fields
}

func queueSizeWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{}
	if v, ok := d.GetOk("items_per_instance"); ok {
		fields["items_per_instance"] = v
	}
	return fields
}

func scheduleWriteFields(d map[string]interface{}) rsc.Fields {
	fields := rsc.Fields{}
	for _, f := range []string{
		"day", "max_count", "min_count", "time",
	} {
		fields[f] = d[f]
	}
	return fields
}
