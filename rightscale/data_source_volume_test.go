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

func TestAccRSVolumeDatasource(t *testing.T) {
	t.Parallel()

	var (
		cloudHref    = os.Getenv("RIGHTSCALE_CLOUD_HREF")
		validObjHref = regexp.MustCompile(fmt.Sprintf("^%s/volumes/[0-9A-Z]?", cloudHref))
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
				Config: testAccRSVolumeDatasource(cloudHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRSVolumeExists("data.rightscale_volume.a_volume", &objHref),
					testAccCheckRSVolumeHref(&objHref, validObjHref),
					testAccCheckRSVolumeKeys("data.rightscale_volume.a_volume"),
				),
			},
		},
	})
}

func testAccRSVolumeDatasource(ch string) string {
	return fmt.Sprintf(`
		data "rightscale_volume" "a_volume" {
			cloud_href = "%s"
		}
		`, ch)
}

func testAccCheckRSVolumeHref(href *string, re *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := re.FindStringSubmatch(*href)
		if len(m) != 1 {
			return fmt.Errorf("invalid HRef received: %s, expected something that matches %v", *href, re)
		}
		return nil
	}

}

func testAccCheckRSVolumeExists(n string, ch *string) resource.TestCheckFunc {
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

func testAccCheckRSVolumeKeys(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attributes := []string{
			"id",
			"href",
			"name",
			"updated_at",
			"status",
			"created_at",
			"description",
			"resource_uid",
			"cloud_href",
		}

		for _, attr := range attributes {
			if _, ok := rs.Primary.Attributes[attr]; !ok {
				return fmt.Errorf("Datasource doesn't contain attribute %s", attr)
			}
		}
		return nil
	}
}
