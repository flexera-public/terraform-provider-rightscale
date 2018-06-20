package rightscale

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/rightscale/rsc/rsapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRSDatacenterDatasource(t *testing.T) {
	t.Parallel()

	var (
		cloudHref           = os.Getenv("RIGHTSCALE_CLOUD_HREF")
		validDatacenterHref = regexp.MustCompile(fmt.Sprintf("^%s/datacenters/[0-9A-Z]?", cloudHref))
		datacenterHref      string
	)

	type cmClient interface {
		API() *rsapi.API
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRSDatacenterDatasource(cloudHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRSDatacenterExists("data.rightscale_datacenter.a_dc", &datacenterHref),
					testAccCheckRSDatacenterHref(&datacenterHref, validDatacenterHref),
					testAccCheckRSDatacenterKeys("data.rightscale_datacenter.a_dc"),
				),
			},
		},
	})
}

func testAccRSDatacenterDatasource(ch string) string {
	return fmt.Sprintf(`
data "rightscale_datacenter" "a_dc" {
	cloud_href = "%s"
}
`, ch)
}

func testAccCheckRSDatacenterHref(href *string, re *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := re.FindStringSubmatch(*href)
		if len(m) != 1 {
			return fmt.Errorf("invalid HRef received: %s, expected something that matches %v", *href, re)
		}
		return nil
	}

}

func testAccCheckRSDatacenterExists(n string, dh *string) resource.TestCheckFunc {
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

func testAccCheckRSDatacenterKeys(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attributes := []string{
			"name",
			"description",
			"href",
			"resource_uid",
			"cloud_href",
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
