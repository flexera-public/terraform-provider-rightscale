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

func TestAccRightScaleServer_basic(t *testing.T) {
	t.Parallel()

	var (
		instanceName   = "terraform-test-instance-" + testString + "-" + acctest.RandString(10)
		serverName     = "terraform-test-server-" + testString + "-" + acctest.RandString(10)
		imageHref      = getTestImageFromEnv()
		typeHref       = getTestInstanceTypeFromEnv()
		cloudHref      = getTestCloudFromEnv()
		templateHref   = getTestTemplateFromEnv()
		deploymentHref = getTestDeploymentFromEnv()
		subnetHref     = getTestSubnetFromEnv()
		server         cm15.Server
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServer_basic(serverName, instanceName, cloudHref, imageHref, typeHref, deploymentHref, templateHref, subnetHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists("rightscale_server.test-server", &server),
					testAccCheckServerDatasource("rightscale_server.test-server", "data.rightscale_server.test-server"),
				),
			},
		},
	})
}

func TestAccRightScaleServer_inputs(t *testing.T) {
	t.Parallel()

	var (
		instanceName   = "terraform-test-instance-" + testString + "-" + acctest.RandString(10)
		serverName     = "terraform-test-server-" + testString + "-" + acctest.RandString(10)
		imageHref      = getTestImageFromEnv()
		typeHref       = getTestInstanceTypeFromEnv()
		cloudHref      = getTestCloudFromEnv()
		templateHref   = getTestTemplateFromEnv()
		deploymentHref = getTestDeploymentFromEnv()
		subnetHref     = getTestSubnetFromEnv()
		server         cm15.Server
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServer_inputs(serverName, instanceName, cloudHref, imageHref, typeHref, deploymentHref, templateHref, subnetHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists("rightscale_server.test-server-inputs", &server),
					testAccCheckServerInputs("rightscale_server.test-server-inputs", serverName, &server),
				),
			},
		},
	})
}

func testAccCheckServerExists(n string, server *cm15.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().ServerLocator(getHrefFromID(rs.Primary.ID))

		found, err := loc.Show(rsapi.APIParams{"view": "default"})
		if err != nil {
			return err
		}

		*server = *found

		return nil
	}
}

func testAccCheckServerDatasource(n, d string) resource.TestCheckFunc {
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
			"optimized",
			"cloud_href",
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

func testAccCheckServerDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_server" {
			continue
		}

		loc := c.ServerLocator(getHrefFromID(rs.Primary.ID))
		servers, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of server: %s", err)
		}
		found := false
		self := strings.Split(rs.Primary.ID, ":")[1]
		for _, server := range servers {
			if string(server.Locator(c).Href) == self && server.State != "terminated" {
				found = true
				break
			}
		}
		if found {
			return fmt.Errorf("server still exists")
		}
	}

	return nil
}

func testAccCheckServerInputs(n string, serverName string, server *cm15.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// take *Server object and execute a show with a specific view to populate the *Instance fields
		c := getCMClient()
		if server == nil {
			return fmt.Errorf("server cm15.Server object is nil: %s", n)
		}
		serverLocator := server.Locator(c)
		if serverLocator == nil {
			return fmt.Errorf("failed to extract server locator from cm15.Server object: %s", n)
		}
		serverInstanceDetails, err := serverLocator.Show(rsapi.APIParams{"view": "instance_detail"})
		if err != nil {
			return fmt.Errorf("failed show call for server with instance_detail view: %s", err)
		}
		nextInstance := serverInstanceDetails.NextInstance
		if nextInstance == nil {
			return fmt.Errorf("failed to extract instance locator from next_instance on cm15.Server object: %s", n)
		}
		currentInstanceLoc := nextInstance.Locator(c)
		if currentInstanceLoc == nil {
			return fmt.Errorf("failed to extract instance locator from currentInstance on cm15.Instance object: %s", n)
		}
		// execute a show with specific view on *Instance object to return array of hashes of inputs
		instance, err := currentInstanceLoc.Show(rsapi.APIParams{"view": "full_inputs_2_0"})
		if err != nil {
			return fmt.Errorf("failed show call for instance: %s", err)
		}
		// iterate over hash and compare with expected value to see if we have a match
		expectedMatch := fmt.Sprintf("text:%v", serverName)
		match := false
		for _, input := range instance.Inputs {
			if input["name"] == "SERVER_HOSTNAME" {
				if input["value"] == expectedMatch {
					match = true
					break
				}
			}
		}
		if !match {
			return fmt.Errorf("unable to return match to verify SERVER_HOSTNAME inputs")
		}
		return nil
	}
}

func testAccServer_basic(serverName string, instanceName string, cloudHref string, imageHref string, typeHref string, deploymentHref string, templateHref string, subnetHref string) string {
	return fmt.Sprintf(`
resource "rightscale_server" "test-server" {
  name                   = %q
  deployment_href        = %q
  instance {
    cloud_href           = %q
    image_href           = %q
    instance_type_href   = %q
    name                 = %q
	server_template_href = %q
	subnet_hrefs         = [%q]
  }
}

data "rightscale_server" "test-server" {
	filter {
		name          = "${rightscale_server.test-server.name}"
	}
}
  `, serverName, deploymentHref, cloudHref, imageHref, typeHref, instanceName, templateHref, subnetHref)
}

func testAccServer_inputs(serverName string, instanceName string, cloudHref string, imageHref string, typeHref string, deploymentHref string, templateHref string, subnetHref string) string {
	return fmt.Sprintf(`
resource "rightscale_server" "test-server-inputs" {
  name                   = %q
  deployment_href        = %q
  instance {
    cloud_href           = %q
    image_href           = %q
    instance_type_href   = %q
		name                 = %q
		inputs {
			SERVER_HOSTNAME = "text:%s"
		}
	server_template_href = %q
	subnet_hrefs         = [%q]
  }
}
  `, serverName, deploymentHref, cloudHref, imageHref, typeHref, instanceName, serverName, templateHref, subnetHref)
}
