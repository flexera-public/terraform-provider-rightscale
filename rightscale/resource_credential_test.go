package rightscale

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRightScaleCredential_basic(t *testing.T) {
	t.Parallel()

	var (
		credentialName        = "terraform-test-credential-" + testString + "-" + acctest.RandString(10)
		credentialValue       = "thisIsATest_thisIsOnlyATest"
		credentialDescription = "A test cred created by the rs tf provider"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCredentialDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCredential_basic(credentialName, credentialValue, credentialDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCredential("rightscale_credential.credential_test", "data.rightscale_credential.credential_test"),
				),
			},
		},
	})
}

func testAccCheckCredential(r string, d string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		ds, ok := s.RootModule().Resources[d]
		if !ok {
			return fmt.Errorf("Not found: %s", d)
		}

		loc := getCMClient().CredentialLocator(getHrefFromID(rs.Primary.ID))

		_, err := loc.Show(nil)
		if err != nil {
			return err
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		credentialAttrToCheck := []string{
			"name",
			"description",
			"value",
			"created_at",
		}

		for _, attr := range credentialAttrToCheck {
			if dsAttr[attr] != rsAttr[attr] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr,
					dsAttr[attr],
					rsAttr[attr],
				)
			}
		}

		return nil
	}
}

func testAccCheckCredentialDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_credential" {
			continue
		}

		loc := c.CredentialLocator(getHrefFromID(rs.Primary.ID))
		credentials, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of credential: %s", err)
		}
		found := false
		self := strings.Split(rs.Primary.ID, ":")[1]
		for _, key := range credentials {
			if string(key.Locator(c).Href) == self {
				found = true
				break
			}
		}
		if found {
			return fmt.Errorf("credential still exists")
		}
	}

	return nil
}

func testAccCredential_basic(name string, value string, description string) string {
	return fmt.Sprintf(`
resource "rightscale_credential" "credential_test" {
	name              = %q
	value	          = %q
	description		  = %q
}

data "rightscale_credential" "credential_test" {
	filter {
		name          = "${rightscale_credential.credential_test.name}"
		description   = "${rightscale_credential.credential_test.description}"
	}
}
`, name, value, description)
}
