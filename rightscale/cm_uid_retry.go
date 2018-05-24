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

// cmUIDRetry is a helper function that takes a resource and a time to retry
// meant to be invoked when datasources use the 'resource_uid' filter.
// This implies a cloud resource that SHOULD showup in RightScale as cloud
// polling updates with a resource created outside of rightscale.
// Ex: Build aws resource vpc, create rs server resource to consume subnet from that vpc.
func cmUIDRetry(client rsc.Client, loc *rsc.Locator, typ string, d *schema.ResourceData, t int) error {
	// verify we didn't set an insane retry time - 1200 seconds == 20 min
	if t > 1200 {
		return fmt.Errorf("[ERROR] A timeout above '%v' seconds is not supported", t)
	}
	timeout := time.After((time.Duration(t)) * time.Second)
	tick := time.Tick(10 + time.Second)
	for {
		select {
		case <-timeout:
			// timeout hit - return an error
			return fmt.Errorf("[ERROR] - Time of '%v' seconds exceeded and resource has not been located", t)
		case <-tick:
			// attempt list call that includes cloud resource_uid
			res, err := client.List(loc, typ, cmFilters(d))
			// error from src - raise and return error
			if err != nil {
				return err
			}
			// no results found - retry
			if len(res) == 0 {
				log.Println("[INFO] Index listing did not locate object with matching resource_uid and filter tags - try again later")
			} else {
				// results found - set and return
				log.Println("[INFO] Success! - Index listing located object with matching resource_uid and filter tags")
				for k, v := range res[0].Fields {
					d.Set(k, v)
				}
				d.SetId(res[0].Locator.Href)
				return nil
			}
		}
	}
}
