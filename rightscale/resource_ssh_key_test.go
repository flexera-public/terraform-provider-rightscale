package rightscale

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rightscale/rsc/cm15"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRightScaleSSHKey_basic(t *testing.T) {
	t.Parallel()

	var (
		sshKeyName = "terraform-test-ssh-key-" + testString + "-" + acctest.RandString(10)
		cloudHref  = getTestCloudFromEnv()
		sshKey     cm15.SshKey
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSSHKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSSHKey_basic(sshKeyName, cloudHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSHKeyExists("rightscale_ssh_key.ssh_key_test", &sshKey),
				),
			},
		},
	})
}

func testAccCheckSSHKeyExists(n string, sshKey *cm15.SshKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().SshKeyLocator(getHrefFromID(rs.Primary.ID))

		found, err := loc.Show(nil)
		if err != nil {
			return err
		}

		*sshKey = *found

		return nil
	}
}

func testAccCheckSSHKeyDatasource(n, d string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Check datasource output matches resource
		ds, ok := s.RootModule().Resources[d]
		if !ok {
			return fmt.Errorf("Not found: %s", d)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		credentialAttrToCheck := []string{
			"name",
			"resource_uid",
			"created_at",
		}

		for _, attr := range credentialAttrToCheck {
			fmt.Printf("k: %v, v: %v \n", attr, rsAttr[attr])
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

func testAccCheckSSHKeyDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_ssh_key" {
			continue
		}

		loc := c.SshKeyLocator(getHrefFromID(rs.Primary.ID))
		sshKeys, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of key: %s", err)
		}
		found := false
		self := strings.Split(rs.Primary.ID, ":")[1]
		for _, key := range sshKeys {
			if string(key.Locator(c).Href) == self {
				found = true
				break
			}
		}
		if found {
			return fmt.Errorf("ssh key still exists")
		}
	}

	return nil
}

func testAccSSHKey_basic(name string, cloud_href string) string {
	return fmt.Sprintf(`
resource "rightscale_ssh_key" "ssh_key_test" {
	name                = %q
	cloud_href          = %q
}

data "rightscale_ssh_key" "ssh_key_test" {
	cloud_href          = %q
	filter {
		name            = "${rightscale_ssh_key.ssh_key_test.name}"
	}
}
`, name, cloud_href, cloud_href)
}
