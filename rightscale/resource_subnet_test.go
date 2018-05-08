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
	subnetDescription = "Terraform RightScale provider test subnet"
	subnetCidrBlock   = "192.168.20.0/24"
)

func TestAccRightScalesubnet(t *testing.T) {
	t.Parallel()

	var (
		subnetName  = "terraform-test-" + testString + acctest.RandString(10)
		depl        cm15.Subnet
		cloudHref   = getTestCloudFromEnv()
		networkHref = getTestNetworkFromEnv()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSubnet(subnetName, subnetDescription, subnetCidrBlock, cloudHref, networkHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists("rightscale_subnet.test_subnet", &depl),
					testAccCheckSubnetDescription(&depl, subnetDescription),
					testAccCheckSubnetCidrBlock(&depl, subnetCidrBlock),
				),
			},
		},
	})
}

func testAccCheckSubnetExists(n string, depl *cm15.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().SubnetLocator(getHrefFromID(rs.Primary.ID))

		found, err := loc.Show()
		if err != nil {
			return err
		}

		*depl = *found

		return nil
	}
}

func testAccCheckSubnetDescription(depl *cm15.Subnet, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.Description != desc {
			return fmt.Errorf("got description %q, expected %q", depl.Description, desc)
		}
		return nil
	}

}

func testAccCheckSubnetCidrBlock(depl *cm15.Subnet, cidr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.CidrBlock != cidr {
			return fmt.Errorf("got cidr_block %q, expected %q", depl.CidrBlock, cidr)
		}
		return nil
	}
}

func testAccCheckSubnetDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_subnet" {
			continue
		}

		loc := c.SubnetLocator(getHrefFromID(rs.Primary.ID))
		depls, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of subnet: %s", err)
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
			return fmt.Errorf("subnet still exists")
		}
	}

	return nil
}

func testAccSubnet(name string, desc string, cidr string, cloud string, network string) string {
	return fmt.Sprintf(`
		resource "rightscale_subnet" "test_subnet" {
		   name = %q
		   description = %q
		   cidr_block = %q
		   cloud_href = %q
		   network_href = %q
		 }
`, name, desc, cidr, cloud, network)
}
