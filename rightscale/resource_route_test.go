package rightscale

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rightscale/rsc/cm15"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	routeFinalDestination = "0.0.0.0/0"
	routeNextHopType      = "network_gateway"
)

func TestAccRightScaleRoute(t *testing.T) {
	t.Parallel()

	var (
		routeDescription = "terraform-test-route" + testString + acctest.RandString(10)
		depl             cm15.Route
		networkGateway   = getTestNetworkGatewayFromEnv()
		routeTable       = getTestRouteTableFromEnv()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRoute(routeDescription, routeFinalDestination, routeNextHopType, networkGateway, routeTable),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteExists("rightscale_route.test_route", &depl),
					testAccCheckRouteDescription(&depl, routeDescription),
				),
			},
		},
	})
}

func testAccCheckRouteExists(n string, depl *cm15.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().RouteLocator(getHrefFromID(rs.Primary.ID))

		found, err := loc.Show()
		if err != nil {
			return err
		}

		*depl = *found

		return nil
	}
}

func testAccCheckRouteDescription(depl *cm15.Route, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.Description != desc {
			return fmt.Errorf("got description %q, expected %q", depl.Description, desc)
		}
		return nil
	}

}

func testAccCheckRouteDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_route" {
			continue
		}

		loc := c.RouteLocator(getHrefFromID(rs.Primary.ID))
		depls, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of route: %s", err)
		}
		found := false
		self := strings.Split(rs.Primary.ID, ":")[1]
		for _, depl := range depls {
			if string(depl.Locator(c).Href) == self {
				found = true
				break
			}
		}
		if found {
			return fmt.Errorf("route still exists")
		}
	}

	return nil
}

func testAccRoute(desc string, dest string, nextHopType string, gateway string, routeTable string) string {
	return fmt.Sprintf(`
		resource "rightscale_route" "test_route" {
		   description = %q
		   destination_cidr_block = %q
		   next_hop_type = %q
		   next_hop_href = %q
		   route_table_href = %q
		 }
`, desc, dest, nextHopType, gateway, routeTable)
}
