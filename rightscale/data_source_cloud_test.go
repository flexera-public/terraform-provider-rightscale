package rightscale

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/rightscale/rsc/rsapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRSCloudDatasource(t *testing.T) {
	t.Parallel()

	var (
		validCloudHref = regexp.MustCompile("^/api/clouds/[0-9]?")
		cloudHref      string
	)

	type cmClient interface {
		API() *rsapi.API
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRSCloudDatasource("amazon"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRSCloudExists("data.rightscale_cloud.a_cloud", &cloudHref),
					testAccCheckRSCloudHref(&cloudHref, validCloudHref),
					testAccCheckRSCloudKeys("data.rightscale_cloud.a_cloud"),
				),
			},
		},
	})
}

func testAccRSCloudDatasource(ct string) string {
	return fmt.Sprintf(`
data "rightscale_cloud" "a_cloud" {
	filter = {
		cloud_type = "%s"
	}
}
`, ct)
}

func testAccCheckRSCloudHref(href *string, re *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := re.FindStringSubmatch(*href)
		if len(m) != 1 {
			return fmt.Errorf("invalid HRef received: %s, expected something that matches %v", *href, re)
		}
		return nil
	}

}

func testAccCheckRSCloudExists(n string, ch *string) resource.TestCheckFunc {
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

func testAccCheckRSCloudKeys(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attributes := []string{
			"name",
			"description",
			"display_name",
			"cloud_type",
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
