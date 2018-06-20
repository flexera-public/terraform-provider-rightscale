package rightscale

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/rightscale/rsc/rsapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRSMCIDatasource(t *testing.T) {
	t.Parallel()

	var (
		validMCIHref = regexp.MustCompile("^/api/multi_cloud_images/[0-9A-Z]?")
		MCIHref      string
	)

	type cmClient interface {
		API() *rsapi.API
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRSMCIDatasource(MCIHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRSMCIExists("data.rightscale_multi_cloud_image.a_mci", &MCIHref),
					testAccCheckRSMCIHref(&MCIHref, validMCIHref),
					testAccCheckRSMCIKeys("data.rightscale_multi_cloud_image.a_mci"),
				),
			},
		},
	})
}

func testAccRSMCIDatasource(ch string) string {
	return `
data "rightscale_multi_cloud_image" "a_mci" {
}
`
}

func testAccCheckRSMCIHref(href *string, re *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := re.FindStringSubmatch(*href)
		if len(m) != 1 {
			return fmt.Errorf("invalid HRef received: %s, expected something that matches %v", *href, re)
		}
		return nil
	}

}

func testAccCheckRSMCIExists(n string, dh *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		*dh = getHrefFromID(rs.Primary.ID)
		return nil
	}
}

func testAccCheckRSMCIKeys(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attributes := []string{
			"name",
			"description",
			"revision",
			"href",
			"id",
		}

		for _, attr := range attributes {
			if _, ok := rs.Primary.Attributes[attr]; !ok {
				return fmt.Errorf("Datasource doesn't contain attribute %s", attr)
			}
		}
		return nil
	}
}
