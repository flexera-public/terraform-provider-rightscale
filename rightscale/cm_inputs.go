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
func cmInputs(d *schema.ResourceData) (rsc.Fields, error) {
	var inputs rsc.Fields
	var mapify map[string]string
	mapify = make(map[string]string)
	if f, ok := d.GetOk("inputs"); ok {
		inputsList := f.(*schema.Set).List()
		for _, inputIF := range inputsList {
			v, ok := inputIF.(string)
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
	}
	return inputs, nil
}
