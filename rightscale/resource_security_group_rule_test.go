package rightscale

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/rightscale/rsc/cm15"
	"github.com/rightscale/rsc/rsapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRightScaleSecurityRuleGroup(t *testing.T) {
	t.Parallel()

	var (
		depl          cm15.SecurityGroupRule
		securityGroup = getTestSecurityGroupFromEnv()
		sgCidrIps     = fmt.Sprintf("192.168.%v.%v/32", acctest.RandIntRange(1, 254), acctest.RandIntRange(1, 254))
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupRule(securityGroup, sgCidrIps),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("rightscale_security_group_rule.test_sg_rule", &depl),
					testAccCheckSecurityGroupRuleCIDR(&depl, sgCidrIps),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupRuleExists(n string, depl *cm15.SecurityGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().SecurityGroupRuleLocator(getHrefFromID(rs.Primary.ID))

		var params rsapi.APIParams
		found, err := loc.Show(params)
		if err != nil {
			return err
		}

		*depl = *found

		return nil
	}
}

func testAccCheckSecurityGroupRuleCIDR(depl *cm15.SecurityGroupRule, cidrIps string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.CidrIps != cidrIps {
			return fmt.Errorf("got cidr_ips %q, expected %q", depl.CidrIps, cidrIps)
		}
		return nil
	}

}

func testAccCheckSecurityGroupRuleDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_security_group_rule" {
			continue
		}

		loc := c.SecurityGroupRuleLocator(getHrefFromID(rs.Primary.ID))
		depls, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of security group rule: %s", err)
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
			return fmt.Errorf("Security Group Rule still exists")
		}
	}

	return nil
}

func testAccSecurityGroupRule(sgHref string, cidrIps string) string {
	return fmt.Sprintf(`
		resource "rightscale_security_group_rule" "test_sg_rule" {
		   security_group_href = %q
		   direction = "ingress"
		   protocol = "tcp"
		   source_type = "cidr_ips"
		   cidr_ips = %q
		   protocol_details {
			   start_port = "22"
			   end_port = "22"
		   }
		 }
`, sgHref, cidrIps)
}
