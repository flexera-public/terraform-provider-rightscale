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

func TestAccRightScaleServerArray_basic(t *testing.T) {
	t.Parallel()

	var (
		serverarrayName = "terraform-test-serverarray-" + testString + acctest.RandString(10)
		instanceName    = "terraform-test-instance-" + testString + "-" + acctest.RandString(10)
		state           = "enabled"
		threshold       = "75"
		minCount        = "2"
		maxCount        = "2"
		datacenterMax   = "4"
		datacenterHref  = getTestDatacenterFromEnv()
		imageHref       = getTestImageFromEnv()
		typeHref        = getTestInstanceTypeFromEnv()
		cloudHref       = getTestCloudFromEnv()
		templateHref    = getTestTemplateFromEnv()
		deploymentHref  = getTestDeploymentFromEnv()
		subnetHref      = getTestSubnetFromEnv()
		serverArray     cm15.ServerArray
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerArrayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServerArray_basic(serverarrayName, state, deploymentHref, threshold, minCount, maxCount, datacenterHref, datacenterMax, instanceName, cloudHref, imageHref, typeHref, templateHref, subnetHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerArrayExists("rightscale_server_array.test_server_array", &serverArray),
					testAccCheckServerArrayHas2Instances(&serverArray),
				),
			},
		},
	})
}

func testAccCheckServerArrayExists(n string, serverArray *cm15.ServerArray) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().ServerArrayLocator(getHrefFromID(rs.Primary.ID))

		found, err := loc.Show(rsapi.APIParams{"view": "default"})
		if err != nil {
			return err
		}

		*serverArray = *found

		return nil
	}
}

func testAccCheckServerArrayHas2Instances(serverArray *cm15.ServerArray) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if serverArray.InstancesCount != 2 {
			return fmt.Errorf("ServerArray contains %v servers (should contain 2)", serverArray.InstancesCount)
		}

		return nil
	}
}

func testAccCheckServerArrayDestroy(s *terraform.State) error {
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

func testAccServerArray_basic(serverarrayName string, state string, deploymentHref string, threshold string, minCount string, maxCount string, datacenterHref string, datacenterMax string, instanceName string, cloudHref string, imageHref string, typeHref string, templateHref string, subnetHref string) string {
	return fmt.Sprintf(`
resource "rightscale_server_array" "test_server_array" {
	array_type = "alert"

	datacenter_policy = [{
		datacenter_href = %q
		max             = %q
		weight          = 100
	}]

	elasticity_params = {
		alert_specific_params = {
		decision_threshold = %q
		}

		bounds = {
		min_count = %q
		max_count = %q
		}

		pacing = {
		resize_down_by = 1
		resize_up_by   = 1
		}
	}

	instance = {
		cloud_href           = %q
		image_href           = %q
		instance_type_href   = %q
		server_template_href = %q
		name                 = %q
		subnet_hrefs         = [%q]
		associate_public_ip_address = true
	}

	name            = %q
	state           = %q
	deployment_href = %q
	}
`, datacenterHref, datacenterMax, threshold, minCount, maxCount, cloudHref, imageHref, typeHref, templateHref, instanceName, subnetHref, serverarrayName, state, deploymentHref)
}
