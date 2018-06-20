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

func TestAccRSVolumeSnapshotDatasource(t *testing.T) {
	t.Parallel()

	var (
		cloudHref    = os.Getenv("RIGHTSCALE_CLOUD_HREF")
		validObjHref = regexp.MustCompile(fmt.Sprintf("^%s/volume_snapshots/[0-9A-Z]?", cloudHref))
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
				Config: testAccRSVolumeSnapshotDatasource(cloudHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRSVolumeSnapshotExists("data.rightscale_volume_snapshot.a_volume_snapshot", &objHref),
					testAccCheckRSVolumeSnapshotHref(&objHref, validObjHref),
					testAccCheckRSVolumeSnapshotKeys("data.rightscale_volume_snapshot.a_volume_snapshot"),
				),
			},
		},
	})
}

func testAccRSVolumeSnapshotDatasource(ch string) string {
	return fmt.Sprintf(`
		data "rightscale_volume_snapshot" "a_volume_snapshot" {
			cloud_href = "%s"
			filter {
				name = "06" // Avoid too many results error
			}
		}
		`, ch)
}

func testAccCheckRSVolumeSnapshotHref(href *string, re *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := re.FindStringSubmatch(*href)
		if len(m) != 1 {
			return fmt.Errorf("invalid HRef received: %s, expected something that matches %v", *href, re)
		}
		return nil
	}

}

func testAccCheckRSVolumeSnapshotExists(n string, ch *string) resource.TestCheckFunc {
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

func testAccCheckRSVolumeSnapshotKeys(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attributes := []string{
			"id",
			"href",
			"name",
			"state",
			"resource_uid",
			"description",
			"size",
			"created_at",
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
