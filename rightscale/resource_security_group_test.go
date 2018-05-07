package rightscale

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rightscale/rsc/cm15"
	"github.com/rightscale/rsc/rsapi"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	securityGroupDescription = "Terraform RightScale provider test Security Group"
)

func TestAccRightScaleSecurityGroup(t *testing.T) {
	t.Parallel()

	var (
		SecurityGroupName = "terraform-test-" + testString + acctest.RandString(10)
		depl              cm15.SecurityGroup
		// This test will execute against default network in this cloud
		cloudHref   = getTestCloudFromEnv()
		networkHref = getTestNetworkFromEnv()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroup(SecurityGroupName, securityGroupDescription, cloudHref, networkHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("rightscale_security_group.test_sg", &depl),
					testAccCheckSecurityGroupDescription(&depl, securityGroupDescription),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupExists(n string, depl *cm15.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().SecurityGroupLocator(getHrefFromID(rs.Primary.ID))

		var params rsapi.APIParams
		found, err := loc.Show(params)
		if err != nil {
			return err
		}

		*depl = *found

		return nil
	}
}

func testAccCheckSecurityGroupDescription(depl *cm15.SecurityGroup, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.Description != desc {
			return fmt.Errorf("got description %q, expected %q", depl.Description, desc)
		}
		return nil
	}

}

func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_security_group" {
			continue
		}

		loc := c.SecurityGroupLocator(getHrefFromID(rs.Primary.ID))
		depls, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of security group: %s", err)
		}
		found := false
		self := strings.Split(rs.Primary.ID, ":")[1]
		for _, depl := range depls {
			if string(depl.Locator(c).Href) == self {
				found = true
				break
			}
		}
		if found {
			return fmt.Errorf("Security Group still exists")
		}
	}

	return nil
}

func testAccSecurityGroup(name string, desc string, cloud string, network string) string {
	return fmt.Sprintf(`
		resource "rightscale_security_group" "test_sg" {
		   name = %q
		   description = %q
		   cloud_href = %q
		   network_href = %q
		 }
`, name, desc, cloud, network)
}
