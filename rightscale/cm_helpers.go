package rightscale

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// cmUIDSet is a helper function that returns true or false if 'resource_uid' is set as part
// of a filter block on a resource.
func cmUIDSet(d *schema.ResourceData) bool {
	if f, ok := d.GetOk("filter"); ok {
		filtersList := f.([]interface{})
		for _, filterIF := range filtersList {
			v := filterIF.(map[string]interface{})
			for filter, value := range v {
				if filter == "resource_uid" && value != "" {
					return true
				}
			}
		}
		// filter is set, but resource_uid is not.
		return false
	}
	// filter is not set
	return false
}

// cmIndexRetry is a helper function that takes a resource and a time to retry.
// Working with resources created outside of rightscale requires a little time for the clouds
// to be polled and those resources found.  This function will keep making index calls
// until the resource is located or the timeout is exceeded.
// Ex: Build aws resource vpc, create rs server resource to consume subnet from that vpc.
func cmIndexRetry(client rsc.Client, loc *rsc.Locator, typ string, d *schema.ResourceData, t int) error {
	// verify we didn't set an insane retry time - 1200 seconds == 20 min
	if t > 1200 {
		return fmt.Errorf("[ERROR] A timeout of '%v' seconds is not supported (1200 seconds max)", t)
	}
	timeout := time.After((time.Duration(t)) * time.Second)
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-timeout:
			// timeout hit - return an error
			return fmt.Errorf("[ERROR] - Time of '%v' seconds exceeded and '%s' resource has not been located", t, typ)
		case <-tick:
			// attempt list call that includes cloud resource_uid
			res, err := client.List(loc, typ, cmFilters(d))
			// error from src - raise and return error
			if err != nil {
				return err
			}
			// no results found - retry
			if len(res) == 0 {
				log.Printf("[DEBUG] Index listing did not locate '%s' object with matching resource_uid and filters - retrying...", typ)
			} else {
				// results found - return
				log.Printf("[DEBUG] Success! - Index listing located '%s' object with matching resource_uid and filters", typ)
				return nil
			}
		}
	}
}

// cmFilters is a helper function that returns fields representing valid
// RightScale rcl filter parameters built from the resource data "filter"
// field.
func cmFilters(d *schema.ResourceData) rsc.Fields {
	var filters rsc.Fields
	var arrify []string
	if f, ok := d.GetOk("filter"); ok {
		filtersList := f.([]interface{})
		for _, filterIF := range filtersList {
			v := filterIF.(map[string]interface{})
			for filter, value := range v {
				if value != "" {
					z := fmt.Sprintf("%v==%v", filter, value)
					arrify = append(arrify, z)
				}
			}
		}
		filters = rsc.Fields{"filter": arrify}
	}
	return filters
}

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
