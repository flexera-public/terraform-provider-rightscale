package rightscale

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rightscale/rsc/cm15"
	"github.com/rightscale/rsc/rsapi"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	routeTableDescription = "Terraform RightScale provider test route table"
)

func TestAccRightScaleRouteTable(t *testing.T) {
	t.Parallel()

	var (
		routeTableName = "terraform-test-" + testString + acctest.RandString(10)
		depl           cm15.RouteTable
		cloudHref      = getTestCloudFromEnv()
		networkHref    = getTestNetworkFromEnv()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteTable(routeTableName, routeTableDescription, cloudHref, networkHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists("rightscale_route_table.test_route_table", &depl),
					testAccCheckRouteTableDescription(&depl, routeTableDescription),
				),
			},
		},
	})
}

func testAccCheckRouteTableExists(n string, depl *cm15.RouteTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().RouteTableLocator(getHrefFromID(rs.Primary.ID))

		var params rsapi.APIParams
		found, err := loc.Show(params)
		if err != nil {
			return err
		}

		*depl = *found

		return nil
	}
}

func testAccCheckRouteTableDescription(depl *cm15.RouteTable, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.Description != desc {
			return fmt.Errorf("got description %q, expected %q", depl.Description, desc)
		}
		return nil
	}

}

func testAccCheckRouteTableDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_route_table" {
			continue
		}

		loc := c.RouteTableLocator(getHrefFromID(rs.Primary.ID))
		depls, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of route table: %s", err)
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
			return fmt.Errorf("route table still exists")
		}
	}

	return nil
}

func testAccRouteTable(name string, desc string, cloud string, network string) string {
	return fmt.Sprintf(`
		resource "rightscale_route_table" "test_route_table" {
		   name = %q
		   description = %q
		   cloud_href = %q
		   network_href = %q
		 }
`, name, desc, cloud, network)
}
