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
	description    = "Terraform RightScale provider test deployment"
	serverTagScope = "deployment"
)

func TestAccRightScaleCMDeployment_basic(t *testing.T) {
	t.Parallel()

	var (
		deploymentName = "terraform-test-" + testString + acctest.RandString(10)
		depl           cm15.Deployment
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCMDeploymentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCMDeployment_basic(deploymentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCMDeploymentExists("rightscale_cm_deployment.foobar", &depl),
					testAccCheckCMDeploymentDescription(&depl, description),
					testAccCheckCMDeploymentServerTagScope(&depl, serverTagScope),
				),
			},
		},
	})
}

func TestAccRightScaleCMDeployment_locked(t *testing.T) {
	t.Parallel()

	var (
		deploymentName = "terraform-test-" + testString + acctest.RandString(10)
		depl           cm15.Deployment
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCMDeploymentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCMDeployment_locked(deploymentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCMDeploymentExists("rightscale_cm_deployment.foobar", &depl),
					testAccCheckCMDeploymentLocked(&depl),
				),
			},
			// unlock so we can delete, also tests updates
			resource.TestStep{
				Config: testAccCMDeployment_unlocked(deploymentName),
			},
		},
	})
}

func testAccCheckCMDeploymentExists(n string, depl *cm15.Deployment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().DeploymentLocator(getHrefFromID(rs.Primary.ID))

		found, err := loc.Show(nil)
		if err != nil {
			return err
		}

		*depl = *found

		return nil
	}
}

func testAccCheckCMDeploymentDescription(depl *cm15.Deployment, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.Description != desc {
			return fmt.Errorf("got description %q, expected %q", depl.Description, desc)
		}
		return nil
	}

}

func testAccCheckCMDeploymentServerTagScope(depl *cm15.Deployment, scope string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.ServerTagScope != scope {
			return fmt.Errorf("got server tag scope %q, expected %q", depl.ServerTagScope, scope)
		}
		return nil
	}
}

func testAccCheckCMDeploymentLocked(depl *cm15.Deployment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// // unlock once we have checked so we can delete the deployment.
		// defer depl.Locator(getCMClient(s)).Unlock()
		if !depl.Locked {
			return fmt.Errorf("expected deployment to be locked")
		}
		return nil
	}
}

func testAccCheckCMDeploymentDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_cm_deployment" {
			continue
		}

		loc := c.DeploymentLocator(getHrefFromID(rs.Primary.ID))
		depls, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of deployment: %s", err)
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
			return fmt.Errorf("deployment still exists")
		}
	}

	return nil
}

func testAccCMDeployment_basic(dep string) string {
	return fmt.Sprintf(`
resource "rightscale_cm_deployment" "foobar" {
	name                = %q
	description         = %q
	server_tag_scope    = %q
}
`, dep, description, serverTagScope)
}

func testAccCMDeployment_locked(dep string) string {
	return fmt.Sprintf(`
resource "rightscale_cm_deployment" "foobar" {
	name             = "%s"
	description      = "Terraform RightScale provider test deployment - locked"
	server_tag_scope = %q
	locked           = true
}
`, dep, serverTagScope)
}

func testAccCMDeployment_unlocked(dep string) string {
	return fmt.Sprintf(`
resource "rightscale_cm_deployment" "foobar" {
	name             = "%s"
	description      = "Terraform RightScale provider test deployment - locked"
	server_tag_scope = %q
	locked           = false
}
`, dep, serverTagScope)
}
