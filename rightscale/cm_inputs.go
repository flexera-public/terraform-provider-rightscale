package rightscale

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// cmInputs is a helper function that returns fields representing valid
// RightScale rcl input parameters built from the resource data "inputs"
// field.
//func cmInputs(d *schema.ResourceData) (rsc.Fields, error) {
func cmInputs(f *schema.Set) (rsc.Fields, error) {
	var inputs rsc.Fields
	mapify := make(map[string]string)
	//f, ok := d.GetOk("inputs")
	//if !ok {
	//	return inputs, nil
	//}
	il := f.List()
	for _, i := range il {
		v, ok := i.(string)
		if !ok {
			return inputs, fmt.Errorf("inputsList does not appear to be properly handled as a string: %v", ok)
		}
		z := strings.Split(v, "=")
		if len(z) != 2 {
			return inputs, fmt.Errorf("input format should be formatted as 'KEY=type:value', split on '=' of '%s' should result in array count of 2, got count of: %v", z, len(z))
		}
		mapify[z[0]] = z[1]
	}
	inputs = rsc.Fields{"inputs": mapify}
	return inputs, nil
}
