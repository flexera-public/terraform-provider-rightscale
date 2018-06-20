package rightscale

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/rightscale/rsc/rsapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRSRouteTableatasource(t *testing.T) {
	t.Parallel()

	var (
		validObjHref = regexp.MustCompile("^/api/route_tables/[0-9]?")
		objHref      string
	)

	type cmClient interface {
		API() *rsapi.API
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRSRouteTableDatasource("amazon"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRSRouteTableExists("data.rightscale_route_table.a_route_table", &objHref),
					testAccCheckRSRouteTableHref(&objHref, validObjHref),
					testAccCheckRSRouteTableKeys("data.rightscale_route_table.a_route_table"),
				),
			},
		},
	})
}

func testAccRSRouteTableDatasource(ct string) string {
	return `
data "rightscale_route_table" "a_route_table" {
}
`
}

func testAccCheckRSRouteTableHref(href *string, re *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := re.FindStringSubmatch(*href)
		if len(m) != 1 {
			return fmt.Errorf("invalid HRef received: %s, expected something that matches %v", *href, re)
		}
		return nil
	}

}

func testAccCheckRSRouteTableExists(n string, ch *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		*ch = getHrefFromID(rs.Primary.ID)
		return nil
	}
}

func testAccCheckRSRouteTableKeys(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attributes := []string{
			"id",
			"href",
			"name",
			"description",
			"resource_uid",
		}

		for _, attr := range attributes {
			if _, ok := rs.Primary.Attributes[attr]; !ok {
				return fmt.Errorf("Datasource doesn't contain attribute %s", attr)
			}
		}
		return nil
	}
}
