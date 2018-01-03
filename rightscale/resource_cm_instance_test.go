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

func TestAccRightScaleCMInstance_basic(t *testing.T) {
	t.Parallel()

	var (
		instanceName = "terraform-test-instance-" + testString + "-" + acctest.RandString(10)
		imageHref    = getTestImageFromEnv()
		typeHref     = getTestInstanceTypeFromEnv()
		cloudHref    = getTestCloudFromEnv()
		inst         cm15.Instance
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCMInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCMInstance_basic(instanceName, cloudHref, imageHref, typeHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCMInstanceExists("rightscale_cm_instance.test-instance", &inst),
				),
			},
		},
	})
}

func TestAccRightScaleCMInstance_userdata(t *testing.T) {
	t.Parallel()

	var (
		instanceName = "terraform-test-instance-" + testString + "-" + acctest.RandString(10)
		imageHref    = getTestImageFromEnv()
		typeHref     = getTestInstanceTypeFromEnv()
		cloudHref    = getTestCloudFromEnv()
		userData     = "UserData" + acctest.RandString(10)
		inst         cm15.Instance
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCMInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCMInstance_userdata(instanceName, cloudHref, imageHref, typeHref, userData),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCMInstanceExists("rightscale_cm_instance.test-instance", &inst),
					testAccCheckCMInstanceUserdata(userData, &inst),
				),
			},
		},
	})
}

func TestAccRightScaleCMInstance_locked(t *testing.T) {
	t.Parallel()

	var (
		instanceName = "terraform-test-instance-" + testString + "-" + acctest.RandString(10)
		imageHref    = getTestImageFromEnv()
		typeHref     = getTestInstanceTypeFromEnv()
		cloudHref    = getTestCloudFromEnv()
		inst         cm15.Instance
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCMInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCMInstance_basic(instanceName, cloudHref, imageHref, typeHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCMInstanceExists("rightscale_cm_instance.test-instance", &inst),
				),
			},
			resource.TestStep{
				Config: testAccCMInstance_locked(instanceName, cloudHref, imageHref, typeHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCMInstanceExists("rightscale_cm_instance.test-instance", &inst),
					testAccCheckCMInstanceLocked(&inst),
				),
			},
			resource.TestStep{
				Config: testAccCMInstance_unlocked(instanceName, cloudHref, imageHref, typeHref),
			},
		},
	})
}

func testAccCheckCMInstanceExists(n string, inst *cm15.Instance) resource.TestCheckFunc {
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

func testAccCheckCMInstanceDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_cm_instance" {
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

func testAccCheckCMInstanceLocked(inst *cm15.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// // unlock once we have checked so we can delete the instance.
		if !inst.Locked {
			return fmt.Errorf("expected instance to be locked")
		}
		return nil
	}
}

func testAccCheckCMInstanceUserdata(userData string, inst *cm15.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if inst.UserData != userData {
			return fmt.Errorf("Instance userdata does not match expectation. Got %q, expected %q", inst.UserData, userData)
		}

		return nil
	}
}

func testAccCMInstance_basic(name string, cloud_href string, image_href string, instance_type_href string) string {
	return fmt.Sprintf(`
resource "rightscale_cm_instance" "test-instance" {
	name                = %q
	cloud_href          = %q
	image_href          = %q
	instance_type_href  = %q
  associate_public_ip_address = true
}
`, name, cloud_href, image_href, instance_type_href)
}

func testAccCMInstance_locked(name string, cloud_href string, image_href string, instance_type_href string) string {
	return fmt.Sprintf(`
resource "rightscale_cm_instance" "test-instance" {
	name                = %q
	cloud_href          = %q
	image_href          = %q
	instance_type_href  = %q
  locked              = true
  associate_public_ip_address = true
}
`, name, cloud_href, image_href, instance_type_href)
}

func testAccCMInstance_unlocked(name string, cloud_href string, image_href string, instance_type_href string) string {
	return fmt.Sprintf(`
resource "rightscale_cm_instance" "test-instance" {
	name                = %q
	cloud_href          = %q
	image_href          = %q
	instance_type_href  = %q
  locked              = false
  associate_public_ip_address = true
}
`, name, cloud_href, image_href, instance_type_href)
}

func testAccCMInstance_userdata(name string, cloud_href string, image_href string, instance_type_href string, userdata string) string {
	return fmt.Sprintf(`
resource "rightscale_cm_instance" "test-instance" {
	name                = %q
	cloud_href          = %q
	image_href          = %q
	instance_type_href  = %q
  user_data           = %q
  associate_public_ip_address = true
}
`, name, cloud_href, image_href, instance_type_href, userdata)
}
