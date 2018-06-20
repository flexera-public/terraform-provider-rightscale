package rightscale

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/rightscale/rsc/rsapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRSServerTemplateDatasource(t *testing.T) {
	t.Parallel()

	var (
		validObjHref = regexp.MustCompile("^/api/server_templates/[0-9]?")
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
				Config: testAccRSServerTemplateDatasource("amazon"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRSServerTemplateExists("data.rightscale_server_template.a_server_template", &objHref),
					testAccCheckRSServerTemplateHref(&objHref, validObjHref),
					testAccCheckRSServerTemplateKeys("data.rightscale_server_template.a_server_template"),
				),
			},
		},
	})
}

func testAccRSServerTemplateDatasource(ct string) string {
	return `
data "rightscale_server_template" "a_server_template" {
}
`
}

func testAccCheckRSServerTemplateHref(href *string, re *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := re.FindStringSubmatch(*href)
		if len(m) != 1 {
			return fmt.Errorf("invalid HRef received: %s, expected something that matches %v", *href, re)
		}
		return nil
	}

}

func testAccCheckRSServerTemplateExists(n string, ch *string) resource.TestCheckFunc {
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

func testAccCheckRSServerTemplateKeys(n string) resource.TestCheckFunc {
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
			"revision",
			"lineage",
		}

		for _, attr := range attributes {
			if _, ok := rs.Primary.Attributes[attr]; !ok {
				return fmt.Errorf("Datasource doesn't contain attribute %s", attr)
			}
		}
		return nil
	}
}
