package rsc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/rightscale/rsc/httpclient"
	"github.com/rightscale/rsc/rsapi"
)

// NOTE: the tests below make "real" API requests to the RightScale platform.
// They use credentials passed in environment variables to perform auth. The
// tests are skipped if the environment variables are missing.
//
// The tests use the following environment variables:
//
//     * RIGHTSCALE_API_TOKEN is the API token used to auth API requests made to RightScale.
//     * RIGHTSCALE_PROJECT_ID is the RightScale project used to run the tests.
//     * DEBUG causes additional output useful to troubleshoot issues.

func launchMockServer(t *testing.T) *httptest.Server {
	retries := 3
	return httptest.NewServer(http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			var (
				response  = ""
				projectID = validProjectID(t)
			)

			t.Logf("URL: %s", request.URL.Path)
			switch request.URL.Path {
			case "/api/oauth2":
				response = `
				{
					"access_token": "%s",
					"expires_in": 7200,
					"token_type": "bearer"
				}`
				response = fmt.Sprintf(response, acctest.RandString(370)) // access_token
			case "/api/sessions", "/api/sessions/accounts":
				response = `
				[{
					"created_at": "2008/08/01 17:09:31 +0000",
					"links": [
						{
							"href": "/api/accounts/%d",
							"rel": "self"
						},
						{
							"href": "/api/users/1",
							"rel": "owner"
						},
						{
							"href": "/api/clusters/3",
							"rel": "cluster"
						}
					],
					"name": "Account one",
					"updated_at": "2018/04/25 06:19:43 +0000"
				},
				{
					"created_at": "2008/08/01 17:01:31 +0000",
					"links": [
						{
							"href": "/api/accounts/%d",
							"rel": "self"
						},
						{
							"href": "/api/users/13",
							"rel": "owner"
						},
						{
							"href": "/api/clusters/24",
							"rel": "cluster"
						}
					],
					"name": "Another account",
					"updated_at": "2018/04/25 06:11:43 +0000"
				}]`
				response = fmt.Sprintf(response, projectID, projectID+3)
			case fmt.Sprintf("/cwf/v1/accounts/%d/processes", projectID):
				writer.Header().Set("Location", fmt.Sprintf("/accounts/%d/processes/5b06d799a17cac6ee9ebd62a", projectID))
			case fmt.Sprintf("/cwf/v1/accounts/%d/processes/5b06d799a17cac6ee9ebd62a", projectID):
				response = `
				{
					"id": "5b06d1b51c028800360030f9",
					"href": "/accounts/62656/processes/5b06d1b51c028800360030f9",
					"name": "07ppy1wzmcsk4",
					"tasks": [
						{
							"id": "5b06d1b51c028800360030f8",
							"href": "/accounts/62656/tasks/5b06d1b51c028800360030f8",
							"name": "/root",
							"progress": {
								"percent": 100,
								"summary": ""
							},
							"status": "completed",
							"created_at": "2018-05-24T14:52:37.500Z",
							"updated_at": "2018-05-24T14:52:37.500Z",
							"finished_at": "2018-05-24T14:52:41.157Z"
						}
					],
					"outputs": [
						%s
					],
					"references": [],
					"variables": [],
					"source": "define main() return $res do\n\t$res = 11 + 31\nend\n",
					"main": "define main() return $res do\n|   $res = 11 + 31\nend",
					"parameters": [],
					"application": "cwfconsole",
					"created_by": {
						"email": "support@example.com",
						"id": 0,
						"name": "Terraform"
					},
					"created_at": "2018-05-24T14:52:37.500Z",
					"updated_at": "2018-05-24T14:52:41.121Z",
					"finished_at": "2018-05-24T14:52:41.162Z",
					"status": "completed",
					"links": {
						"tasks": {
							"href": "/accounts/62656/processes/5b06d1b51c028800360030f9/tasks"
						}
					}
				}`
				outputs := `{
								"name": "$res",
								"value": {
									"kind": "number",
									"value": 42
								}
							}`
				retries = retries - 1
				if retries == 0 {
					response = fmt.Sprintf(response, outputs)
				} else {
					response = fmt.Sprintf(response, "")
				}
			default:
				// TODO write a log line warning of unknown PATH
				return
			}
			if len(response) > 0 {
				h := writer.Header()
				h.Set("Content-Type", "application/json")
				writer.Write([]byte(response))
			}
		}))
}
func TestRunProcess(t *testing.T) {
	service := launchMockServer(t)

	rb := rshosts
	hin := httpclient.Insecure
	defer func() {
		rshosts = rb
		httpclient.Insecure = hin
		service.Close()
	}()
	httpclient.Insecure = true
	rshosts = []string{service.URL}

	client, err := New(validToken(t), validProjectID(t))
	if err != nil {
		t.Errorf("got error %q, expected none", err)
	}

	source := `define main() return $res do
	$res = 11 + 31
end
`
	process, err := client.RunProcess(source, nil)
	if err != nil {
		t.Errorf("got error %q, expected none", err)
		return
	}
	if process.Outputs["$res"] != "42" {
		t.Errorf("got $res equal %s, expected 42", process.Outputs["$res"])
	}
}
func TestAuthenticate(t *testing.T) {
	token := validToken(t)
	project := validProjectID(t)
	cases := []struct {
		Name      string
		Token     string
		ProjectID int
		Error     string
	}{
		{"valid", token, project, ""},
		{"invalid-token", "foo", project, "failed to authenticate"},
		{"invalid-project-id", token, 0, "session does not give access to project 0"},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := New(c.Token, c.ProjectID)
			if c.Error == "" {
				if err != nil {
					t.Errorf("got error %q, expected none", err)
				}
			} else {
				if err == nil {
					t.Errorf("got no error, expected %q", c.Error)
				} else {
					if err.Error() != c.Error {
						t.Errorf("got error %q, expected %q", err.Error(), c.Error)
					}
				}
			}
		})
	}
}

func TestList(t *testing.T) {
	const (
		namespace = "rs_cm"
		typ       = "clouds"
	)
	var (
		token   = validToken(t)
		project = validProjectID(t)
	)
	cases := []struct {
		Name           string
		Namespace      string
		Type           string
		Href           string
		Link           string
		Filters        Fields
		ExpectedPrefix string
		ExpectedError  string
	}{
		{"clouds", namespace, typ, "", "", nil, "", ""},
		{"filtered", namespace, typ, "", "", Fields{"filter[]": "name==EC2"}, "EC2", ""},
		{"linked", namespace, typ, "/api/clouds/1", "datacenters", nil, "", ""},
		{"linked-and-filtered", namespace, typ, "/api/clouds/1", "datacenters", Fields{"filter[]": "name==us-east-1a"}, "us-east-1a", ""},
		{"no-namespace", "", typ, "", "", nil, "", "resource locator is invalid: namespace is missing"},
		{"no-type-no-href", "", "", "", "", nil, "", "resource locator is invalid: namespace is missing"},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			cl, err := New(token, project)
			if err != nil {
				t.Fatal(err)
			}
			loc := &Locator{Namespace: c.Namespace, Type: c.Type, Href: c.Href}

			clouds, err := cl.List(loc, c.Link, c.Filters)

			if err != nil {
				if c.ExpectedError == "" {
					t.Errorf("got error %q, expected none", err)
					return
				}
				if c.ExpectedError != err.Error() {
					t.Errorf("got error %q, expected %q", err, c.ExpectedError)
				}
				return
			}
			if c.ExpectedPrefix != "" {
				for i, cloud := range clouds {
					if !strings.HasPrefix(cloud.Fields["name"].(string), c.ExpectedPrefix) {
						t.Errorf("got name %q at index %d, expected prefix %q", cloud.Fields["name"], i, c.ExpectedPrefix)
					}
				}
			}
		})
	}
}

func TestCreate(t *testing.T) {
	const (
		namespace = "rs_cm"
		typ       = "deployment"
		deplDesc  = "Created by tests"
	)
	depl := "Terraform Provider Test Deployment " + acctest.RandString(4)
	token := validToken(t)
	project := validProjectID(t)
	cl, err := New(token, project)
	if err != nil {
		t.Fatal(err)
	}
	rs := cl.(*client).rs
	cleanDeployment(t, depl, rs)
	defer cleanDeployment(t, depl, rs)

	_, err = cl.Create(namespace, typ, Fields{"deployment": Fields{"name": depl, "description": deplDesc}})

	if err != nil {
		t.Errorf("got error %q, expected none", err)
		return
	}
	d := showDeployment(t, depl, rs)
	if d == nil {
		t.Errorf("deployment not created")
		return
	}
	if d["name"].(string) != depl {
		t.Errorf("got deployment with name %v, expected %q", d["name"], depl)
	}
	if d["description"].(string) != deplDesc {
		t.Errorf("got deployment with description %v, expected %q", d["description"], deplDesc)
	}
}

func TestDelete(t *testing.T) {
	const (
		namespace = "rs_cm"
		typ       = "deployment"
		deplDesc  = "Created by tests"
	)
	depl := "Terraform Provider Test Deployment " + acctest.RandString(4)
	token := validToken(t)
	project := validProjectID(t)
	cl, err := New(token, project)
	if err != nil {
		t.Fatal(err)
	}
	rs := cl.(*client).rs
	cleanDeployment(t, depl, rs)
	defer cleanDeployment(t, depl, rs)

	res, err := cl.Create(namespace, typ, Fields{"deployment": Fields{"name": depl, "description": deplDesc}})
	if err != nil {
		t.Fatal(err)
	}

	err = cl.Delete(res.Locator)

	if err != nil {
		t.Errorf("got error %q, expected none", err)
		return
	}
	d := showDeployment(t, depl, rs)
	if d != nil {
		t.Errorf("deployment not deleted")
		return
	}
}

func validToken(t *testing.T) string {
	tok := os.Getenv("RIGHTSCALE_API_TOKEN")
	if tok == "" {
		t.Skip("RIGHTSCALE_API_TOKEN environment variable not defined, skipping authentication test")
	}
	return tok
}

func validProjectID(t *testing.T) int {
	pid := os.Getenv("RIGHTSCALE_PROJECT_ID")
	if pid == "" {
		t.Skip("RIGHTSCALE_PROJECT_ID environment variable not defined")
	}
	projectID, err := strconv.Atoi(pid)
	if err != nil {
		t.Fatal(err)
	}
	return projectID
}

func showDeployment(t *testing.T, depl string, rs *rsapi.API) map[string]interface{} {
	req, err := rs.BuildHTTPRequest("GET", "/api/deployments", "1.5", rsapi.APIParams{"filter[]": "name==" + depl}, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := rs.PerformRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to retrieve deployment: index returned %q", resp.Status)
	}
	ms, err := rs.LoadResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	if len(ms.([]interface{})) == 0 {
		return nil
	}
	return ms.([]interface{})[0].(map[string]interface{})
}

func cleanDeployment(t *testing.T, depl string, rs *rsapi.API) {
	var id string
	{
		m := showDeployment(t, depl, rs)
		if m == nil {
			return
		}
		links := m["links"].([]interface{})
		var href string
		for _, l := range links {
			rel := l.(map[string]interface{})["rel"].(string)
			if rel != "self" {
				continue
			}
			href = l.(map[string]interface{})["href"].(string)
			break
		}
		idx := strings.LastIndex(href, "/")
		id = href[idx+1:]
		if id == "" {
			t.Fatalf("failed to retrieve deployment id, href: %q", href)
		}
	}
	req, err := rs.BuildHTTPRequest("DELETE", "/api/deployments/"+id, "1.5", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := rs.PerformRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("failed to delete deployment: destroy returned %q", resp.Status)
	}
}

func TestOnlyPopulated(t *testing.T) {
	var expectedResult Fields = map[string]interface{}{
		"keyString": "value1",
		"keyList":   []interface{}{0, 1, 2, 3},
		"keyMap": Fields{
			"subkey1": "subvalue1",
		},
		"keyInt": 0,
	}
	var testFields Fields = map[string]interface{}{
		"keyString":      "value1",
		"keyStringEmpty": "",
		"keyList":        []interface{}{0, 1, 2, 3},
		"keyListEmpty":   []interface{}{},
		"keyMap": Fields{
			"subkey1": "subvalue1",
			"subkey2": "",
		},
		"keyMapEmpty": Fields{},
		"keyInt":      0,
	}
	if !reflect.DeepEqual(testFields.onlyPopulated(), expectedResult) {
		t.Errorf("Result of onlyPopulated was incorrect, got: %v, expected: %v", testFields.onlyPopulated(), expectedResult)
	}
}

func TestAnalyzeSource(t *testing.T) {
	const (
		invalidSource = `
define main( do
	foo
end
`
		goodSourceNoReturn = `
define main() do
	bar
end
`
		goodSourceWithReturn = `
define main() return $out1, $out2 $suma do
	$out1 = 156.5534
	$out2 = 42421000
	$suma = $out1 + $out2
end
`
	)

	var sourcetests = []struct {
		name           string
		source         string
		valid          bool
		expectsOutputs bool
	}{
		{"invalidSource", invalidSource, false, false},
		{"goodSourceNoReturn", goodSourceNoReturn, true, false},
		{"goodSourceWithReturn", goodSourceWithReturn, true, true},
	}

	for _, tt := range sourcetests {
		t.Run(tt.source, func(t *testing.T) {
			expectsOutputs, err := analyzeSource(tt.source)
			if tt.valid && err != nil {
				t.Errorf("source `%s` should be valid (but got error `%v`)", tt.name, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("source `%s` should be invalid", tt.name)
			}
			if expectsOutputs != tt.expectsOutputs {
				t.Errorf("source `%s` got incorrect expectsOutputs value: `%t` (should be `%t`)", tt.name, expectsOutputs, tt.expectsOutputs)
			}
		})
	}
}
