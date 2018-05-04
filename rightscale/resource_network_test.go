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

const (
	networkDescription = "Terraform RightScale provider test Network"
	networkName        = "TerraformProviderTest"
	networkCidrBlock   = "192.168.10.0/24"
)

func TestAccRightScaleNetwork(t *testing.T) {
	t.Parallel()

	var (
		NetworkName = "terraform-test-" + testString + acctest.RandString(10)
		depl        cm15.Network
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetwork(NetworkName, networkDescription, networkCidrBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("rightscale_network.test_network", &depl),
					testAccCheckNetworkDescription(&depl, networkDescription),
					testAccCheckNetworkCidrBlock(&depl, networkCidrBlock),
				),
			},
		},
	})
}

func testAccCheckNetworkExists(n string, depl *cm15.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().NetworkLocator(getHrefFromID(rs.Primary.ID))

		found, err := loc.Show()
		if err != nil {
			return err
		}

		*depl = *found

		return nil
	}
}

func testAccCheckNetworkDescription(depl *cm15.Network, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.Description != desc {
			return fmt.Errorf("got description %q, expected %q", depl.Description, desc)
		}
		return nil
	}

}

func testAccCheckNetworkCidrBlock(depl *cm15.Network, cidr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.CidrBlock != cidr {
			return fmt.Errorf("got cidr_block %q, expected %q", depl.CidrBlock, cidr)
		}
		return nil
	}
}

func testAccCheckNetworkDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_network" {
			continue
		}

		loc := c.NetworkLocator(getHrefFromID(rs.Primary.ID))
		depls, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of Network: %s", err)
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
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}

func testAccNetwork(name string, desc string, cidr string) string {
	return fmt.Sprintf(`
		resource "rightscale_network" "test_network" {
		   name = %q
		   description = %q
		   cidr_block = %q
		   cloud_href = "/api/clouds/6"
		 }
`, name, desc, cidr)
}
