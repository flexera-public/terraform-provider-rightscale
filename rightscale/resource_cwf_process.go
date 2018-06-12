package rightscale

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

func resourceCWFProcess() *schema.Resource {
	return &schema.Resource{
		Read:   resourceCWFProcessRead,
		Delete: resourceCWFProcessDelete,
		Create: resourceCWFProcessCreate,

		Schema: map[string]*schema.Schema{
			"parameters": &schema.Schema{
				Description: "main definition parameter values in order of declaration",
				Type:        schema.TypeList,
				Elem:        cwfParameterResource(),
				Optional:    true,
				ForceNew:    true,
			},
			"source": &schema.Schema{
				Description: "process source code, must contain a definition called 'main'",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			// Read-only fields
			"outputs": {
				Type:        schema.TypeMap,
				Description: `process outputs if any`,
				Computed:    true,
			},
			"error": {
				Type:        schema.TypeString,
				Description: `process execution error if any`,
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: `process status, one of "completed", "failed", "canceled" or "aborted"`,
				Computed:    true,
			},
		},
	}
}

func cwfParameterResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kind": {
				Type:        schema.TypeString,
				Description: entityKindDesc,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"array", "boolean", "collection",
					"datetime", "declaration", "null",
					"number", "object", "string",
				}, false),
			},
			"value": {
				Type:        schema.TypeString,
				Description: entityValueDesc,
				Required:    true,
			},
		},
	}
}

func resourceCWFProcessCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)

	var params []*rsc.Parameter
	{
		p, ok := d.GetOk("parameters")
		if ok {
			params = make([]*rsc.Parameter, len(p.([]interface{})))
			for i, param := range p.([]interface{}) {
				pa, err := parseParam(param.(map[string]interface{}))
				if err != nil {
					return err
				}
				params[i] = pa
			}
		}
	}

	proc, err := client.RunProcess(d.Get("source").(string), params)
	if err != nil {
		return err
	}

	d.Set("outputs", proc.Outputs)
	d.Set("status", proc.Status)
	d.SetId(proc.Href)
	if proc.Error != nil {
		return fmt.Errorf(proc.Error.Error())
	}
	return nil
}

func resourceCWFProcessRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)

	p, err := client.GetProcess(d.Id())
	if err != nil {
		if err == rsc.ErrNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("outputs", p.Outputs)
	d.Set("status", p.Status)
	d.Set("error", p.Error)

	return nil
}

func resourceCWFProcessDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)

	return client.DeleteProcess(d.Id())
}

func parseParam(param map[string]interface{}) (*rsc.Parameter, error) {
	var (
		kind  string
		value string
	)
	{
		k, ok := param["kind"]
		if !ok {
			return nil, fmt.Errorf("invalid parameter definition: missing 'kind' key")
		}
		kind, ok = k.(string)
		if !ok {
			return nil, fmt.Errorf("invalid parameter definition: 'kind' is not a string, got %v", k)
		}
		v, ok := param["value"]
		if !ok {
			return nil, fmt.Errorf("invalid parameter definition: missing 'value' key")
		}
		value, ok = v.(string)
		if !ok {
			return nil, fmt.Errorf("invalid parameter definition: 'value' is not a string, got %v", v)
		}
	}
	var iv interface{}
	if kind == "string" {
		iv = value
	} else {
		err := json.Unmarshal([]byte(value), &iv)
		if err != nil {
			return nil, fmt.Errorf("parameter 'value' is not valid json: %s", err)
		}
	}
	if kind == "array" {
		if err := entifyArrayParameter(&iv); err != nil {
			return nil, fmt.Errorf("parameter definition is not valid: %s", err)
		}
	}

	if err := validateParameter("parameter", kind, iv); err != nil {
		return nil, fmt.Errorf("parameter definition is not valid: %s", err)
	}
	return &rsc.Parameter{Kind: rsc.ParameterKind(kind), Value: iv}, nil
}

func entifyArrayParameter(iv *interface{}) error {
	es, ok := (*iv).([]interface{})
	if !ok {
		return fmt.Errorf("entity with kind 'array' is not a slice of entities: %v", iv)
	}
	for i, ei := range es {
		e, err := toEntity(ei)
		if err != nil {
			return err
		}
		es[i] = e
	}
	return nil
}

func validateParameter(ctx string, k string, iv interface{}) error {
	var pref string
	if ctx != "" {
		pref = ctx + ": "
	}
	switch k {
	case "string":
		if _, ok := iv.(string); !ok {
			return fmt.Errorf("%sentity with kind 'string' is not a string: %v", pref, iv)
		}
	case "number":
		if _, ok := iv.(float64); !ok {
			return fmt.Errorf("%sentity with kind 'number' is not a number: %v", pref, iv)
		}
	case "bool":
		if _, ok := iv.(bool); !ok {
			return fmt.Errorf("%sentity with kind 'bool' is not a boolean value: %v", pref, iv)
		}
	case "null":
		if iv != nil {
			return fmt.Errorf("%sentity with kind 'null' must have a nil value, got %v", pref, iv)
		}
	case "array":
		// Necessary checks already run at entifyArrayParameters
	case "object":
		oi, ok := iv.(map[string]interface{})
		if !ok {
			return fmt.Errorf("entity with kind 'object' is not a map, got %v", iv)
		}
		for k, v := range oi {
			e, err := entityFromInterface(v)
			if err != nil {
				return fmt.Errorf("object value is not an entity: %s", err)
			}
			if err := validateParameter(ctx+fmt.Sprintf("[%s]", k), string(e.Kind), e.Value); err != nil {
				return fmt.Errorf("%s.%s: %s", ctx, k, err)
			}
		}

	case "declaration":
		m, ok := iv.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%sentity with kind 'declaration' does not define a declaration: %v", pref, iv)
		}

		n, ok := m["namespace"]
		if !ok {
			return fmt.Errorf("%scollection is missing 'namespace' key", pref)
		}
		if _, ok := n.(string); !ok {
			return fmt.Errorf("%scollection 'namespace' key must have a string value, got %v", pref, n)
		}

		t, ok := m["type"]
		if !ok {
			return fmt.Errorf("%scollection is missing 'type' key", pref)
		}
		if _, ok := t.(string); !ok {
			return fmt.Errorf("%scollection 'type' key must have a string value, got %v", pref, t)
		}

		ds, ok := m["fields"]
		if !ok {
			return fmt.Errorf("%scollection is missing 'fields' key", pref)
		}
		if _, ok := ds.(map[string]interface{}); !ok {
			return fmt.Errorf("%scollection 'fields' key must be an object, got %v", pref, n)
		}

	case "collection":
		m, ok := iv.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%sentity with kind 'collection' does not define a collection: %v", pref, iv)
		}

		n, ok := m["namespace"]
		if !ok {
			return fmt.Errorf("%scollection is missing 'namespace' key", pref)
		}
		if _, ok := n.(string); !ok {
			return fmt.Errorf("%scollection 'namespace' key must have a string value, got %v", pref, n)
		}

		t, ok := m["type"]
		if !ok {
			return fmt.Errorf("%scollection is missing 'type' key", pref)
		}
		if _, ok := t.(string); !ok {
			return fmt.Errorf("%scollection 'type' key must have a string value, got %v", pref, t)
		}

		hs, ok := m["hrefs"]
		if !ok {
			return fmt.Errorf("%scollection is missing 'hrefs' key", pref)
		}
		hrefs, ok := hs.([]interface{})
		if !ok {
			return fmt.Errorf("%scollection 'hrefs' key must be a slice of strings value, got %v", pref, n)
		}
		for _, href := range hrefs {
			if _, ok := href.(string); !ok {
				return fmt.Errorf("%scollection 'hrefs' key must be a slice of strings value, got %v", pref, n)
			}
		}

		ds, ok := m["details"]
		if !ok {
			return fmt.Errorf("%scollection is missing 'details' key", pref)
		}
		details, ok := ds.([]interface{})
		if !ok {
			return fmt.Errorf("%scollection 'details' key must be a slice of objects, got %v", pref, n)
		}
		for _, href := range details {
			if _, ok := href.(map[string]interface{}); !ok {
				return fmt.Errorf("%scollection 'details' key must be a slice of objects, got %v", pref, n)
			}
		}

		if len(hrefs) != len(details) {
			return fmt.Errorf("collection 'hrefs' and 'details' slices must have the same number of elements, 'hrefs' has %d elements while details has %d", len(hrefs), len(details))
		}

	default:
		return fmt.Errorf("%sinvalid entity kind %q", pref, k)
	}

	return nil
}

// entityFromInterface returns an entity from a interface value created by
// unmarshaling an entity JSON representation.
func entityFromInterface(iv interface{}) (*rsc.Parameter, error) {
	m, ok := iv.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value is not an entity: %v", iv)
	}
	k, ok := m["kind"]
	if !ok {
		return nil, fmt.Errorf("value is not an entity (missing 'kind' key): %v", iv)
	}
	kind, ok := k.(string)
	if !ok {
		return nil, fmt.Errorf("entity kind is not a string, got %v", kind)
	}
	value, ok := m["value"]
	if !ok {
		return nil, fmt.Errorf("value is not an entity (missing 'value' key): %v", iv)
	}
	if err := validateParameterKind(kind); err != nil {
		return nil, err
	}
	return &rsc.Parameter{Kind: rsc.ParameterKind(kind), Value: value}, nil
}

// validateParameterKind makes sure the given string is a valid ParameterKind.
func validateParameterKind(v string) error {
	switch rsc.ParameterKind(v) {
	case rsc.KindString, rsc.KindNumber, rsc.KindBool,
		rsc.KindNull, rsc.KindArray, rsc.KindObject, rsc.KindCollection,
		rsc.KindDeclaration:
		return nil
	}
	return fmt.Errorf("%q is not a valid entity kind", v)
}

// toEntity converts interface to entity (rsc.Parameter with Kind and Value)
func toEntity(entity interface{}) (*rsc.Parameter, error) {
	var value interface{}
	kind := rsc.KindString
	switch actual := entity.(type) {
	case string:
		value = actual
	case bool:
		kind = rsc.KindBool
		value = actual
	case float64:
		kind = rsc.KindNumber
		value = actual
	case []interface{}:
		kind = rsc.KindArray
		values := make([]*rsc.Parameter, len(actual))
		for i, v := range actual {
			e, err := toEntity(v)
			if err != nil {
				return nil, err
			}
			values[i] = e
		}
		value = values
	case map[string]interface{}:
		kind = rsc.KindObject
		values := make(map[string]*rsc.Parameter, len(actual))
		for k, v := range actual {
			e, err := toEntity(v)
			if err != nil {
				return nil, err
			}
			values[k] = e
		}
		value = values
	case nil:
		kind = rsc.KindNull
		value = nil
	default:
		return nil, fmt.Errorf("unknown type when converting to entity: %T", actual)
	}
	return &rsc.Parameter{
		Kind:  kind,
		Value: value,
	}, nil
}

const entityKindDesc = `Entities have a kind and a value.
  * The kind for a JSON value corresponds to the JSON type name ('string',
  'number', 'bool', 'null', 'array', 'object')
  * The kind for a resource declaration is 'declaration'
  * The kind for a resource collection is 'collection'`

const entityValueDesc = `Entities have a kind and a value.
  * The value for entities of kind 'string', 'number', 'bool', 'null', 'array'
    or 'object' is the JSON reprensentation of the underlying value.
  * The value for a resource collection is a JSON object with the following keys:
      - 'namespace' contains a string representing the collection namespace (e.g. 'rs')
      - 'type' contains a string representing the collection type (e.g. 'servers')
      - 'hrefs' contains an arry of strings representing the hrefs of the resources in the collection
      - 'details' contains an array of hashes representing the resource attributes
  * The value for a resource declaration is a JSON object with the following keys:
      - 'namespace' contains a string representing the declaration namespace (e.g. 'rs')
      - 'type' contains a string reprensenting the declaration type (e.g. 'servers')
      - 'fields' contains a hash reprensenting the data needed to create the resource
      - 'dependencies' contains an array that lists names of references this declaration depends on`
