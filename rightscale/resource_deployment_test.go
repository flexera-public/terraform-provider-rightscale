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

func TestAccRightScaleDeployment_basic(t *testing.T) {
	t.Parallel()

	var (
		deploymentName = "terraform-test-" + testString + acctest.RandString(10)
		depl           cm15.Deployment
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeploymentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDeployment_basic(deploymentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentExists("rightscale_deployment.foobar", &depl),
					testAccCheckDeploymentDescription(&depl, description),
					testAccCheckDeploymentServerTagScope(&depl, serverTagScope),
					testAccCheckDeploymentDatasource("rightscale_deployment.foobar", "data.rightscale_deployment.foobar"),
				),
			},
		},
	})
}

func TestAccRightScaleDeployment_locked(t *testing.T) {
	t.Parallel()

	var (
		deploymentName = "terraform-test-" + testString + acctest.RandString(10)
		depl           cm15.Deployment
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeploymentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDeployment_locked(deploymentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentExists("rightscale_deployment.foobar", &depl),
					testAccCheckDeploymentLocked(&depl),
				),
			},
			// unlock so we can delete, also tests updates
			resource.TestStep{
				Config: testAccDeployment_unlocked(deploymentName),
			},
		},
	})
}

func testAccCheckDeploymentExists(n string, depl *cm15.Deployment) resource.TestCheckFunc {
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

func testAccCheckDeploymentDatasource(n, d string) resource.TestCheckFunc {
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
			"description",
			"locked",
			"sever_tag_scope",
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

func testAccCheckDeploymentDescription(depl *cm15.Deployment, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.Description != desc {
			return fmt.Errorf("got description %q, expected %q", depl.Description, desc)
		}
		return nil
	}

}

func testAccCheckDeploymentServerTagScope(depl *cm15.Deployment, scope string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if depl.ServerTagScope != scope {
			return fmt.Errorf("got server tag scope %q, expected %q", depl.ServerTagScope, scope)
		}
		return nil
	}
}

func testAccCheckDeploymentLocked(depl *cm15.Deployment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// // unlock once we have checked so we can delete the deployment.
		// defer depl.Locator(getCMClient(s)).Unlock()
		if !depl.Locked {
			return fmt.Errorf("expected deployment to be locked")
		}
		return nil
	}
}

func testAccCheckDeploymentDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_deployment" {
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

func testAccDeployment_basic(dep string) string {
	return fmt.Sprintf(`
resource "rightscale_deployment" "foobar" {
	name                = %q
	description         = %q
	server_tag_scope    = %q
}

data "rightscale_deployment" "foobar" {
	filter {
		name          = "${rightscale_deployment.foobar.name}"
		description   = "${rightscale_deployment.foobar.description}"
	}
}
`, dep, description, serverTagScope)
}

func testAccDeployment_locked(dep string) string {
	return fmt.Sprintf(`
resource "rightscale_deployment" "foobar" {
	name             = "%s"
	description      = "Terraform RightScale provider test deployment - locked"
	server_tag_scope = %q
	locked           = true
}
`, dep, serverTagScope)
}

func testAccDeployment_unlocked(dep string) string {
	return fmt.Sprintf(`
resource "rightscale_deployment" "foobar" {
	name             = "%s"
	description      = "Terraform RightScale provider test deployment - locked"
	server_tag_scope = %q
	locked           = false
}
`, dep, serverTagScope)
}
