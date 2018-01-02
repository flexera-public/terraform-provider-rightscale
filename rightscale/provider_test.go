package rightscale

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rightscale/rsc/cm15"
	"github.com/rightscale/rsc/rsapi"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testString string

var envVars map[string]string = map[string]string{
	"cred":         "RIGHTSCALE_API_TOKEN",
	"project":      "RIGHTSCALE_PROJECT_ID",
	"cloud":        "RIGHTSCALE_CLOUD_HREF",
	"instanceType": "RIGHTSCALE_INSTANCE_TYPE_HREF",
	"image":        "RIGHTSCALE_IMAGE_HREF",
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"rightscale": testAccProvider,
	}
	testString = acctest.RandString(10)
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	for _, envVar := range envVars {
		if v := envSearch(envVar); v == "" {
			t.Fatalf("%s must be set for acceptance tests", envVar)
		}
	}
}

// getCMClient returns a low level API 1.5 client.
func getCMClient() *cm15.API {
	type cmClient interface {
		API() *rsapi.API
	}

	c := testAccProvider.Meta().(cmClient)
	return &cm15.API{API: c.API()}
}

// testAccPreCheck ensures at least one of the project env variables is set.
func getTestProjectFromEnv() string {
	return envSearch(envVars["project"])
}

// testAccPreCheck ensures at least one of the credentials env variables is set.
func getTestCredsFromEnv() string {
	return envSearch(envVars["cred"])
}

// testAccPreCheck ensures at least one of the credentials env variables is set.
func getTestCloudFromEnv() string {
	return envSearch(envVars["cloud"])
}

// testAccPreCheck ensures at least one of the credentials env variables is set.
func getTestInstanceTypeFromEnv() string {
	return envSearch(envVars["instanceType"])
}

// testAccPreCheck ensures at least one of the credentials env variables is set.
func getTestImageFromEnv() string {
	return envSearch(envVars["image"])
}

func getHrefFromID(id string) string {
	return strings.Split(id, ":")[1]
}

func envSearch(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return ""
}
