package rightscale

import (
	"fmt"

	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// cmInputs is a helper function that returns fields representing valid
// RightScale rcl input parameters built from the resource data "inputs"
// field.
func cmInputs(f []interface{}) (rsc.Fields, error) {
	var inputs rsc.Fields
	mapify := make(map[string]string)
	for _, i := range f {
		v, ok := i.(map[string]interface{})
		if !ok {
			return inputs, fmt.Errorf("inputsList does not appear to be properly handled as a string: %v", ok)
		}
		for k, v2 := range v {
			mapify[k] = v2.(string)
		}
	}
	inputs = rsc.Fields{"inputs": mapify}
	return inputs, nil
}
