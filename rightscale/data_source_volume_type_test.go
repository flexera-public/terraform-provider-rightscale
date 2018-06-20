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

func TestAccRSVolumeTypeDatasource(t *testing.T) {
	t.Parallel()

	var (
		cloudHref    = os.Getenv("RIGHTSCALE_CLOUD_HREF")
		validObjHref = regexp.MustCompile(fmt.Sprintf("^%s/volume_types/[0-9A-Z]?", cloudHref))
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
				Config: testAccRSVolumeTypeDatasource(cloudHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRSVolumeTypeExists("data.rightscale_volume_type.a_volume_type", &objHref),
					testAccCheckRSVolumeTypeHref(&objHref, validObjHref),
					testAccCheckRSVolumeTypeKeys("data.rightscale_volume_type.a_volume_type"),
				),
			},
		},
	})
}

func testAccRSVolumeTypeDatasource(ch string) string {
	return fmt.Sprintf(`
		data "rightscale_volume_type" "a_volume_type" {
			cloud_href = "%s"
		}
		`, ch)
}

func testAccCheckRSVolumeTypeHref(href *string, re *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := re.FindStringSubmatch(*href)
		if len(m) != 1 {
			return fmt.Errorf("invalid HRef received: %s, expected something that matches %v", *href, re)
		}
		return nil
	}

}

func testAccCheckRSVolumeTypeExists(n string, ch *string) resource.TestCheckFunc {
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

func testAccCheckRSVolumeTypeKeys(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attributes := []string{
			"name",
			"description",
			"resource_uid",
			"id",
			"href",
			"cloud_href",
			"created_at",
			"updated_at",
		}

		for _, attr := range attributes {
			if _, ok := rs.Primary.Attributes[attr]; !ok {
				return fmt.Errorf("Datasource doesn't contain attribute %s", attr)
			}
		}
		return nil
	}
}
