package rightscale

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

func TestAccRightScaleCWFProcess_basic(t *testing.T) {
	t.Parallel()

	const src = `
define main() return $out do
    $out = 42
end
`
	var process rsc.Process

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCWFProcessDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCWFProcess_basic(src),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCWFProcessExists("rightscale_cwf_process.foobar", &process),
					testAccCheckCWFProcessOutput(&process, []string{"$out"}, []interface{}{string("42")}),
					testAccCheckCWFProcessStatus(&process, "completed"),
				),
			},
		},
	})
}

func TestAccRightScaleCWFProcess_params(t *testing.T) {
	t.Parallel()

	const src = `
define main($p) return $out do
    $out = $p
end
`
	var process rsc.Process

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCWFProcessDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCWFProcess_params(src, []string{"foobar"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCWFProcessExists("rightscale_cwf_process.foobar", &process),
					testAccCheckCWFProcessOutput(&process, []string{"$out"}, []interface{}{"foobar"}),
					testAccCheckCWFProcessStatus(&process, "completed"),
				),
			},
		},
	})
}

func TestAccRightScaleCWFProcess_error(t *testing.T) {
	t.Parallel()

	const src = `
define main() return $out do
    $out = 42
    raise "test error"
end
`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCWFProcessDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccCWFProcess_basic(src),
				ExpectError: regexp.MustCompile("test error"),
			},
		},
	})
}

func testAccCheckCWFProcessExists(n string, p *rsc.Process) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		c := testAccProvider.Meta().(rsc.Client)
		proc, err := c.GetProcess(rs.Primary.ID)
		if err != nil {
			return err
		}

		*p = *proc

		return nil
	}
}

func testAccCheckCWFProcessOutput(p *rsc.Process, names []string, vals []interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(p.Outputs) != len(names) {
			return fmt.Errorf("got %d output(s), expected %d", len(p.Outputs), len(names))
		}
		for i, e := range names {
			out, ok := p.Outputs[e]
			if !ok {
				return fmt.Errorf("output %q is missing", e)
			}
			if out != vals[i] {
				return fmt.Errorf("output %q has value %T: %#v, expected %T: %#v", e, out, out, vals[i], vals[i])
			}
		}
		return nil
	}
}

func testAccCheckCWFProcessStatus(p *rsc.Process, status string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if p.Status != status {
			return fmt.Errorf("got status %q, expected %q", p.Status, status)
		}
		return nil
	}
}

func testAccCheckCWFProcessError(p *rsc.Process, err string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !strings.Contains(p.Error.Error(), err) {
			return fmt.Errorf("got error %q, expected to contain %q", p.Error, err)
		}
		return nil
	}
}

func testAccCWFProcess_basic(src string) string {
	return fmt.Sprintf(`
resource "rightscale_cwf_process" "foobar" {
	source =  <<EOF
%s
EOF
}
`, src)
}

func testAccCWFProcess_params(src string, pvals []string) string {
	vs := make([]string, len(pvals))
	for i, pv := range pvals {
		vs[i] = fmt.Sprintf(`{
    kind = "string"
    value = %q
}`, pv)
	}
	return fmt.Sprintf(`
resource "rightscale_cwf_process" "foobar" {
	source     = <<EOF
%s
EOF
	parameters = [%s]
}
`, src, strings.Join(vs, ", "))
}

func testAccCheckCWFProcessDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(rsc.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_cwf_process" {
			continue
		}

		_, err := c.GetProcess(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("process still exists")
		}
		if err != rsc.ErrNotFound {
			return fmt.Errorf("failed to get process: %s", err)
		}
	}
	return nil
}
