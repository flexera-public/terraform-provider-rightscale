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

func TestAccRightScaleInstance_basic(t *testing.T) {
	t.Parallel()

	var (
		instanceName = "terraform-test-instance-" + testString + "-" + acctest.RandString(10)
		imageHref    = getTestImageFromEnv()
		typeHref     = getTestInstanceTypeFromEnv()
		cloudHref    = getTestCloudFromEnv()
		subnetHref   = getTestSubnetFromEnv()
		inst         cm15.Instance
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstance_basic(instanceName, cloudHref, subnetHref, imageHref, typeHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("rightscale_instance.test-instance", &inst),
				),
			},
		},
	})
}

func TestAccRightScaleInstance_userdata(t *testing.T) {
	t.Parallel()

	var (
		instanceName = "terraform-test-instance-" + testString + "-" + acctest.RandString(10)
		imageHref    = getTestImageFromEnv()
		typeHref     = getTestInstanceTypeFromEnv()
		cloudHref    = getTestCloudFromEnv()
		subnetHref   = getTestSubnetFromEnv()
		userData     = "UserData" + acctest.RandString(10)
		inst         cm15.Instance
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstance_userdata(instanceName, cloudHref, subnetHref, imageHref, typeHref, userData),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("rightscale_instance.test-instance", &inst),
					testAccCheckInstanceUserdata(userData, &inst),
				),
			},
		},
	})
}

func TestAccRightScaleInstance_locked(t *testing.T) {
	t.Parallel()

	var (
		instanceName = "terraform-test-instance-" + testString + "-" + acctest.RandString(10)
		imageHref    = getTestImageFromEnv()
		typeHref     = getTestInstanceTypeFromEnv()
		cloudHref    = getTestCloudFromEnv()
		subnetHref   = getTestSubnetFromEnv()
		inst         cm15.Instance
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstance_basic(instanceName, cloudHref, subnetHref, imageHref, typeHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("rightscale_instance.test-instance", &inst),
				),
			},
			resource.TestStep{
				Config: testAccInstance_locked(instanceName, cloudHref, subnetHref, imageHref, typeHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("rightscale_instance.test-instance", &inst),
					testAccCheckInstanceLocked(&inst),
				),
			},
			resource.TestStep{
				Config: testAccInstance_unlocked(instanceName, cloudHref, subnetHref, imageHref, typeHref),
			},
		},
	})
}

func testAccCheckInstanceExists(n string, inst *cm15.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().InstanceLocator(getHrefFromID(rs.Primary.ID))

		found, err := loc.Show(rsapi.APIParams{"view": "full"})
		if err != nil {
			return err
		}

		*inst = *found

		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_instance" {
			continue
		}

		loc := c.InstanceLocator(getHrefFromID(rs.Primary.ID))
		insts, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of instance: %s", err)
		}
		found := false
		self := strings.Split(rs.Primary.ID, ":")[1]
		for _, inst := range insts {
			if string(inst.Locator(c).Href) == self && inst.State != "terminated" {
				found = true
				break
			}
		}
		if found {
			return fmt.Errorf("instance still exists")
		}
	}

	return nil
}

func testAccCheckInstanceLocked(inst *cm15.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// // unlock once we have checked so we can delete the instance.
		if !inst.Locked {
			return fmt.Errorf("expected instance to be locked")
		}
		return nil
	}
}

func testAccCheckInstanceUserdata(userData string, inst *cm15.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if inst.UserData != userData {
			return fmt.Errorf("Instance userdata does not match expectation. Got %q, expected %q", inst.UserData, userData)
		}

		return nil
	}
}

func testAccInstance_basic(name string, cloud_href string, subnet_href string, image_href string, instance_type_href string) string {
	return fmt.Sprintf(`
resource "rightscale_instance" "test-instance" {
	name                = %q
	cloud_href          = %q
	subnet_hrefs        = [%q]
	image_href          = %q
	instance_type_href  = %q
  associate_public_ip_address = true
}
`, name, cloud_href, subnet_href, image_href, instance_type_href)
}

func testAccInstance_locked(name string, cloud_href string, subnet_href string, image_href string, instance_type_href string) string {
	return fmt.Sprintf(`
resource "rightscale_instance" "test-instance" {
	name                = %q
	cloud_href          = %q
	subnet_hrefs        = [%q]
	image_href          = %q
	instance_type_href  = %q
  locked              = true
  associate_public_ip_address = true
}
`, name, cloud_href, subnet_href, image_href, instance_type_href)
}

func testAccInstance_unlocked(name string, cloud_href string, subnet_href string, image_href string, instance_type_href string) string {
	return fmt.Sprintf(`
resource "rightscale_instance" "test-instance" {
	name                = %q
	cloud_href          = %q
	subnet_hrefs        = [%q]
	image_href          = %q
	instance_type_href  = %q
  locked              = false
  associate_public_ip_address = true
}
`, name, cloud_href, subnet_href, image_href, instance_type_href)
}

func testAccInstance_userdata(name string, cloud_href string, subnet_href string, image_href string, instance_type_href string, userdata string) string {
	return fmt.Sprintf(`
resource "rightscale_instance" "test-instance" {
	name                = %q
	cloud_href          = %q
	subnet_hrefs        = [%q]
	image_href          = %q
	instance_type_href  = %q
  user_data           = %q
  associate_public_ip_address = true
}
`, name, cloud_href, subnet_href, image_href, instance_type_href, userdata)
}
