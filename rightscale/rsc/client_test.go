package rsc

import (
	"encoding/json"
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

type mockServer struct {
	or        []string // rshosts previous value
	ohi       bool     // httpclient.Insecure previous value
	service   *httptest.Server
	projectID int
}

func (ms *mockServer) launch(t *testing.T, testCase string) Client {
	if os.Getenv("RSC_NOMOCK") != "" {
		ms.service = nil
		ms.projectID = validProjectID(t)
		c, err := New(validToken(t), ms.projectID)
		if err != nil {
			t.Fatal(err)
		}
		return c
	}

	ms.projectID = 62656
	retries := 3
	ms.service = httptest.NewServer(http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			var response string

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
				if len(request.URL.Query()) == 0 {
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
					response = fmt.Sprintf(response, ms.projectID, ms.projectID+3)
				} else {
					response = `
					{
						"actions": [],
						"message": "You have successfully logged into the RightScale API.",
						"links": [
							{
								"rel": "account",
								"href": "/api/accounts/%s"
							},
							{
								"rel": "user",
								"href": "/api/users/11111"
							}
						]
					}`
					response = fmt.Sprintf(response, ms.projectID)
				}
			case fmt.Sprintf("/cwf/v1/accounts/%d/processes", ms.projectID):
				switch testCase {
				case "runProcess", "getProcess":
					writer.Header().Set("Location", fmt.Sprintf("/accounts/%d/processes/5b06d1b51c028800360030f9", ms.projectID))
				case "createServer", "createServerNoOutputs":
					writer.Header().Set("Location", fmt.Sprintf("/accounts/%d/processes/5b082948a17cac6ee9ece729", ms.projectID))
				case "runRCLWithDefinitions":
					writer.Header().Set("Location", fmt.Sprintf("/accounts/%d/processes/5b11716f1c02882cf0fdaa84", ms.projectID))
				default:
					panic(fmt.Errorf("Unknown testCase: %s", testCase))
				}
			case fmt.Sprintf("/cwf/v1/accounts/%d/processes/5b06d1b51c028800360030f9", ms.projectID):
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
				if retries > 0 {
					response = fmt.Sprintf(response, "")
				} else {
					response = fmt.Sprintf(response, outputs)
				}
			case fmt.Sprintf("/cwf/v1/accounts/%d/processes/5b082948a17cac6ee9ece729", ms.projectID):
				response_running := `
				{
					"id": "5b082948a17cac6ee9ece729",
					"href": "/accounts/62656/processes/5b082948a17cac6ee9ece729",
					"name": "0nwzhxbpxdn2z",
					"tasks": [
						{
							"id": "5b082948a17cac6ee9ece728",
							"href": "/accounts/62656/tasks/5b082948a17cac6ee9ece728",
							"name": "/root",
							"progress": {
								"percent": 60,
								"summary": "Retrieving field 'state'",
								"expression": {
									"id": "5b0829971136b00001909e2c",
									"href": "/accounts//expressions/5b0829971136b00001909e2c",
									"source": "@server.state",
									"variables": [],
									"references": []
								}
							},
							"status": "activity",
							"created_at": "2018-05-25T15:18:32.054Z",
							"updated_at": "2018-05-25T15:18:32.054Z"
						}
					],
					"outputs": [],
					"references": [],
					"variables": [],
					"source": "define main() return $href, $fields do\n\t$href = \"\"\n\t@server = rs_cm.servers.empty()\n\tsub timeout: 1h do\n\t\t@res = {\"fields\":{\"server\":{\"deployment_href\":\"/api/deployments/936965004\",\"instance\":{\"associate_public_ip_address\":true,\"cloud_href\":\"/api/clouds/1\",\"image_href\":\"/api/clouds/1/images/E0HCVNHNAV8KK\",\"instance_type_href\":\"/api/clouds/1/instance_types/9K1AU4K4RCBU4\",\"ip_forwarding_enabled\":false,\"name\":\"terraform-test-instance-7hcxelcntc-29ley7ipse\",\"server_template_href\":\"/api/server_templates/402254004\",\"subnet_hrefs\":[\"/api/clouds/1/subnets/52NUHI2B8LVH1\"]},\"name\":\"terraform-test-server-7hcxelcntc-8k8ae2u7ia\"}},\"namespace\":\"rs_cm\",\"type\":\"servers\"}\n\t\tcall server_provision_tf(@res) retrieve @server\n\t\t$href   = @server.href\n\t\t$res    = to_object(@res)\n\tend\n$final_state = @server.state\n\tif $final_state == \"operational\"\n\t\t$res = to_object(@server)\n    $fields = to_json($res[\"details\"][0])\n    @server = rs_cm.get(href: @server.href)\n  else\n    $server_name = @server.name\n    raise \"Failed to provision server. Expected state 'operational' but got '\" + $final_state + \"' for server: \" + $server_name + \" at href: \" + $href\nend\nend\n# custom provision that does not auto-cleanup on error\n\tdefine server_provision_tf(@res) return @server do\n\t\t# use RS canned provision to create\n\t\tcall rs__cwf_simple_provision(@res) retrieve @server\n\t\t$object = to_object(@res)\n\t\t# use custom launch to avoid cleanup on error\n\t\tcall tf_server_wait_for_provision(@server) retrieve @server\n\tend\n\n\tdefine tf_server_wait_for_provision(@server) return @server do\n\t\t$server_name = to_s(@server.name)\n\t\tsub on_error: tf_server_handle_launch_failure(@server) do\n\t\t\t@server.launch()\n\t\tend\n\t\t$final_state = \"launching\"\n\t\t# use RS canned logic to capture launching server state\n\t\tsub on_error: rs__cwf_skip_any_error() do\n\t\t\tsleep_until @server.state =~ \"^(operational|stranded|stranded in booting|stopped|terminated|inactive|error)$\"\n\t\t\t$final_state = @server.state\n\t\tend\n\tend\n\n\t# spit out error from launch call\n\tdefine tf_server_handle_launch_failure(@server) do\n\t\t$server_name = @server.name\n\t\tif $_errors \u0026\u0026 $_errors[0] \u0026\u0026 $_errors[0][\"response\"]\n\t\t\traise \"Error trying to launch server (\" + $server_name + \"): \" + $_errors[0][\"response\"][\"body\"]\n\t\telse\n\t\t\traise \"Error trying to launch server (\" + $server_name + \")\"\n\t\tend\n\tend\n",
					"main": "define main() return $href, $fields do\n|   $href = \"\"\n|   @server = rs_cm.servers.empty()\n|   sub timeout: \"1h\" do\n|   |   @res = { \"fields\": { \"server\": { \"deployment_href\": \"/api/deployments/936965004\", \"instance\": { \"associate_public_ip_address\": true, \"cloud_href\": \"/api/clouds/1\", \"image_href\": \"/api/clouds/1/images/E0HCVNHNAV8KK\", \"instance_type_href\": \"/api/clouds/1/instance_types/9K1AU4K4RCBU4\", \"ip_forwarding_enabled\": false, \"name\": \"terraform-test-instance-7hcxelcntc-29ley7ipse\", \"server_template_href\": \"/api/server_templates/402254004\", \"subnet_hrefs\": [\"/api/clouds/1/subnets/52NUHI2B8LVH1\"] }, \"name\": \"terraform-test-server-7hcxelcntc-8k8ae2u7ia\" } }, \"namespace\": \"rs_cm\", \"type\": \"servers\" }\n|   |   call server_provision_tf(@res) retrieve @server\n|   |   $href = @server.href\n|   |   $res = to_object(@res)\n|   end\n|   $final_state = @server.state\n|   if $final_state == \"operational\"\n|   |   $res = to_object(@server)\n|   |   $fields = to_json($res[\"details\"][0])\n|   |   @server = rs_cm.get({ \"href\": @server.href })\n|   elsif true|   |   $server_name = @server.name\n|   |   raise \"Failed to provision server. Expected state 'operational' but got '\" + $final_state + \"' for server: \" + $server_name + \" at href: \" + $href\n|   end\nend",
					"parameters": [],
					"application": "cwfconsole",
					"created_by": {
						"email": "support@rightscale.com",
						"id": 0,
						"name": "Terraform"
					},
					"created_at": "2018-05-25T15:18:32.054Z",
					"updated_at": "2018-05-25T15:18:33.978Z",
					"status": "running",
					"links": {
						"tasks": {
							"href": "/accounts/62656/processes/5b082948a17cac6ee9ece729/tasks"
						}
					}
				}
				`
				response_completed_no_outputs := `
				{
					"id": "5b082948a17cac6ee9ece729",
					"href": "/accounts/62656/processes/5b082948a17cac6ee9ece729",
					"name": "0nwzhxbpxdn2z",
					"tasks": [
						{
							"id": "5b082948a17cac6ee9ece728",
							"href": "/accounts/62656/tasks/5b082948a17cac6ee9ece728",
							"name": "/root",
							"progress": {
								"percent": 60,
								"summary": "Retrieving field 'state'",
								"expression": {
									"id": "5b0829971136b00001909e2c",
									"href": "/accounts//expressions/5b0829971136b00001909e2c",
									"source": "@server.state",
									"variables": [],
									"references": []
								}
							},
							"status": "activity",
							"created_at": "2018-05-25T15:18:32.054Z",
							"updated_at": "2018-05-25T15:18:32.054Z"
						}
					],
					"outputs": [],
					"references": [],
					"variables": [],
					"source": "define main() return $href, $fields do\n\t$href = \"\"\n\t@server = rs_cm.servers.empty()\n\tsub timeout: 1h do\n\t\t@res = {\"fields\":{\"server\":{\"deployment_href\":\"/api/deployments/936965004\",\"instance\":{\"associate_public_ip_address\":true,\"cloud_href\":\"/api/clouds/1\",\"image_href\":\"/api/clouds/1/images/E0HCVNHNAV8KK\",\"instance_type_href\":\"/api/clouds/1/instance_types/9K1AU4K4RCBU4\",\"ip_forwarding_enabled\":false,\"name\":\"terraform-test-instance-7hcxelcntc-29ley7ipse\",\"server_template_href\":\"/api/server_templates/402254004\",\"subnet_hrefs\":[\"/api/clouds/1/subnets/52NUHI2B8LVH1\"]},\"name\":\"terraform-test-server-7hcxelcntc-8k8ae2u7ia\"}},\"namespace\":\"rs_cm\",\"type\":\"servers\"}\n\t\tcall server_provision_tf(@res) retrieve @server\n\t\t$href   = @server.href\n\t\t$res    = to_object(@res)\n\tend\n$final_state = @server.state\n\tif $final_state == \"operational\"\n\t\t$res = to_object(@server)\n    $fields = to_json($res[\"details\"][0])\n    @server = rs_cm.get(href: @server.href)\n  else\n    $server_name = @server.name\n    raise \"Failed to provision server. Expected state 'operational' but got '\" + $final_state + \"' for server: \" + $server_name + \" at href: \" + $href\nend\nend\n# custom provision that does not auto-cleanup on error\n\tdefine server_provision_tf(@res) return @server do\n\t\t# use RS canned provision to create\n\t\tcall rs__cwf_simple_provision(@res) retrieve @server\n\t\t$object = to_object(@res)\n\t\t# use custom launch to avoid cleanup on error\n\t\tcall tf_server_wait_for_provision(@server) retrieve @server\n\tend\n\n\tdefine tf_server_wait_for_provision(@server) return @server do\n\t\t$server_name = to_s(@server.name)\n\t\tsub on_error: tf_server_handle_launch_failure(@server) do\n\t\t\t@server.launch()\n\t\tend\n\t\t$final_state = \"launching\"\n\t\t# use RS canned logic to capture launching server state\n\t\tsub on_error: rs__cwf_skip_any_error() do\n\t\t\tsleep_until @server.state =~ \"^(operational|stranded|stranded in booting|stopped|terminated|inactive|error)$\"\n\t\t\t$final_state = @server.state\n\t\tend\n\tend\n\n\t# spit out error from launch call\n\tdefine tf_server_handle_launch_failure(@server) do\n\t\t$server_name = @server.name\n\t\tif $_errors \u0026\u0026 $_errors[0] \u0026\u0026 $_errors[0][\"response\"]\n\t\t\traise \"Error trying to launch server (\" + $server_name + \"): \" + $_errors[0][\"response\"][\"body\"]\n\t\telse\n\t\t\traise \"Error trying to launch server (\" + $server_name + \")\"\n\t\tend\n\tend\n",
					"main": "define main() return $href, $fields do\n|   $href = \"\"\n|   @server = rs_cm.servers.empty()\n|   sub timeout: \"1h\" do\n|   |   @res = { \"fields\": { \"server\": { \"deployment_href\": \"/api/deployments/936965004\", \"instance\": { \"associate_public_ip_address\": true, \"cloud_href\": \"/api/clouds/1\", \"image_href\": \"/api/clouds/1/images/E0HCVNHNAV8KK\", \"instance_type_href\": \"/api/clouds/1/instance_types/9K1AU4K4RCBU4\", \"ip_forwarding_enabled\": false, \"name\": \"terraform-test-instance-7hcxelcntc-29ley7ipse\", \"server_template_href\": \"/api/server_templates/402254004\", \"subnet_hrefs\": [\"/api/clouds/1/subnets/52NUHI2B8LVH1\"] }, \"name\": \"terraform-test-server-7hcxelcntc-8k8ae2u7ia\" } }, \"namespace\": \"rs_cm\", \"type\": \"servers\" }\n|   |   call server_provision_tf(@res) retrieve @server\n|   |   $href = @server.href\n|   |   $res = to_object(@res)\n|   end\n|   $final_state = @server.state\n|   if $final_state == \"operational\"\n|   |   $res = to_object(@server)\n|   |   $fields = to_json($res[\"details\"][0])\n|   |   @server = rs_cm.get({ \"href\": @server.href })\n|   elsif true|   |   $server_name = @server.name\n|   |   raise \"Failed to provision server. Expected state 'operational' but got '\" + $final_state + \"' for server: \" + $server_name + \" at href: \" + $href\n|   end\nend",
					"parameters": [],
					"application": "cwfconsole",
					"created_by": {
						"email": "support@rightscale.com",
						"id": 0,
						"name": "Terraform"
					},
					"created_at": "2018-05-25T15:18:32.054Z",
					"updated_at": "2018-05-25T15:18:33.978Z",
					"status": "completed",
					"links": {
						"tasks": {
							"href": "/accounts/62656/processes/5b082948a17cac6ee9ece729/tasks"
						}
					}
				}
				`
				response_completed := `
				{
					"id": "5b082948a17cac6ee9ece729",
					"href": "/accounts/62656/processes/5b082948a17cac6ee9ece729",
					"name": "0nwzhxbpxdn2z",
					"tasks": [
						{
							"id": "5b082948a17cac6ee9ece728",
							"href": "/accounts/62656/tasks/5b082948a17cac6ee9ece728",
							"name": "/root",
							"progress": {
								"percent": 100,
								"summary": ""
							},
							"status": "completed",
							"created_at": "2018-05-25T15:18:32.054Z",
							"updated_at": "2018-05-25T15:18:32.054Z",
							"finished_at": "2018-05-25T15:19:54.377Z"
						}
					],
					"outputs": [
						{
							"name": "$href",
							"value": {
								"kind": "string",
								"value": "/api/servers/1797452004"
							}
						},
						{
							"name": "$fields",
							"value": {
								"kind": "string",
								"value": "{\"created_at\":\"2018/05/25 15:18:36 +0000\",\"description\":null,\"links\":[{\"rel\":\"self\",\"href\":\"/api/servers/1797452004\"},{\"href\":\"/api/deployments/936965004\",\"rel\":\"deployment\"},{\"href\":\"/api/clouds/1/instances/FEHVIJTCERGGU\",\"rel\":\"next_instance\"},{\"href\":\"/api/servers/1797452004/alert_specs\",\"rel\":\"alert_specs\"},{\"href\":\"/api/servers/1797452004/alerts\",\"rel\":\"alerts\"}],\"name\":\"terraform-test-server-7hcxelcntc-8k8ae2u7ia\",\"next_instance\":{\"private_ip_addresses\":[],\"state\":\"inactive\",\"name\":\"terraform-test-server-7hcxelcntc-8k8ae2u7ia\",\"cloud_specific_attributes\":{},\"created_at\":\"2018/05/25 15:18:36 +0000\",\"ip_forwarding_enabled\":false,\"public_ip_addresses\":[],\"updated_at\":\"2018/05/25 15:18:36 +0000\",\"associate_public_ip_address\":true,\"links\":[{\"href\":\"/api/clouds/1/instances/FEHVIJTCERGGU\",\"rel\":\"self\"},{\"rel\":\"cloud\",\"href\":\"/api/clouds/1\"},{\"href\":\"/api/deployments/936965004\",\"rel\":\"deployment\"},{\"href\":\"/api/server_templates/402254004\",\"rel\":\"server_template\"},{\"inherited_source\":\"server_template\",\"rel\":\"multi_cloud_image\",\"href\":\"/api/multi_cloud_images/442090004\"},{\"href\":\"/api/servers/1797452004\",\"rel\":\"parent\"},{\"href\":\"/api/clouds/1/instances/FEHVIJTCERGGU/volume_attachments\",\"rel\":\"volume_attachments\"},{\"href\":\"/api/clouds/1/instances/FEHVIJTCERGGU/inputs\",\"rel\":\"inputs\"},{\"href\":\"/api/clouds/1/instances/FEHVIJTCERGGU/monitoring_metrics\",\"rel\":\"monitoring_metrics\"},{\"rel\":\"alerts\",\"href\":\"/api/clouds/1/instances/FEHVIJTCERGGU/alerts\"},{\"href\":\"/api/clouds/1/instances/FEHVIJTCERGGU/alert_specs\",\"rel\":\"alert_specs\"}],\"pricing_type\":\"fixed\",\"resource_uid\":\"e53da384-602e-11e8-97b9-0242ac110002\",\"actions\":[{\"rel\":\"launch\"}],\"locked\":false},\"state\":\"inactive\",\"updated_at\":\"2018/05/25 15:18:36 +0000\",\"actions\":[{\"rel\":\"launch\"},{\"rel\":\"clone\"}]}"
							}
						}
					],
					"references": [],
					"variables": [],
					"source": "define main() return $href, $fields do\n\t$href = \"\"\n\t@server = rs_cm.servers.empty()\n\tsub timeout: 1h do\n\t\t@res = {\"fields\":{\"server\":{\"deployment_href\":\"/api/deployments/936965004\",\"instance\":{\"associate_public_ip_address\":true,\"cloud_href\":\"/api/clouds/1\",\"image_href\":\"/api/clouds/1/images/E0HCVNHNAV8KK\",\"instance_type_href\":\"/api/clouds/1/instance_types/9K1AU4K4RCBU4\",\"ip_forwarding_enabled\":false,\"name\":\"terraform-test-instance-7hcxelcntc-29ley7ipse\",\"server_template_href\":\"/api/server_templates/402254004\",\"subnet_hrefs\":[\"/api/clouds/1/subnets/52NUHI2B8LVH1\"]},\"name\":\"terraform-test-server-7hcxelcntc-8k8ae2u7ia\"}},\"namespace\":\"rs_cm\",\"type\":\"servers\"}\n\t\tcall server_provision_tf(@res) retrieve @server\n\t\t$href   = @server.href\n\t\t$res    = to_object(@res)\n\tend\n$final_state = @server.state\n\tif $final_state == \"operational\"\n\t\t$res = to_object(@server)\n    $fields = to_json($res[\"details\"][0])\n    @server = rs_cm.get(href: @server.href)\n  else\n    $server_name = @server.name\n    raise \"Failed to provision server. Expected state 'operational' but got '\" + $final_state + \"' for server: \" + $server_name + \" at href: \" + $href\nend\nend\n# custom provision that does not auto-cleanup on error\n\tdefine server_provision_tf(@res) return @server do\n\t\t# use RS canned provision to create\n\t\tcall rs__cwf_simple_provision(@res) retrieve @server\n\t\t$object = to_object(@res)\n\t\t# use custom launch to avoid cleanup on error\n\t\tcall tf_server_wait_for_provision(@server) retrieve @server\n\tend\n\n\tdefine tf_server_wait_for_provision(@server) return @server do\n\t\t$server_name = to_s(@server.name)\n\t\tsub on_error: tf_server_handle_launch_failure(@server) do\n\t\t\t@server.launch()\n\t\tend\n\t\t$final_state = \"launching\"\n\t\t# use RS canned logic to capture launching server state\n\t\tsub on_error: rs__cwf_skip_any_error() do\n\t\t\tsleep_until @server.state =~ \"^(operational|stranded|stranded in booting|stopped|terminated|inactive|error)$\"\n\t\t\t$final_state = @server.state\n\t\tend\n\tend\n\n\t# spit out error from launch call\n\tdefine tf_server_handle_launch_failure(@server) do\n\t\t$server_name = @server.name\n\t\tif $_errors \u0026\u0026 $_errors[0] \u0026\u0026 $_errors[0][\"response\"]\n\t\t\traise \"Error trying to launch server (\" + $server_name + \"): \" + $_errors[0][\"response\"][\"body\"]\n\t\telse\n\t\t\traise \"Error trying to launch server (\" + $server_name + \")\"\n\t\tend\n\tend\n",
					"main": "define main() return $href, $fields do\n|   $href = \"\"\n|   @server = rs_cm.servers.empty()\n|   sub timeout: \"1h\" do\n|   |   @res = { \"fields\": { \"server\": { \"deployment_href\": \"/api/deployments/936965004\", \"instance\": { \"associate_public_ip_address\": true, \"cloud_href\": \"/api/clouds/1\", \"image_href\": \"/api/clouds/1/images/E0HCVNHNAV8KK\", \"instance_type_href\": \"/api/clouds/1/instance_types/9K1AU4K4RCBU4\", \"ip_forwarding_enabled\": false, \"name\": \"terraform-test-instance-7hcxelcntc-29ley7ipse\", \"server_template_href\": \"/api/server_templates/402254004\", \"subnet_hrefs\": [\"/api/clouds/1/subnets/52NUHI2B8LVH1\"] }, \"name\": \"terraform-test-server-7hcxelcntc-8k8ae2u7ia\" } }, \"namespace\": \"rs_cm\", \"type\": \"servers\" }\n|   |   call server_provision_tf(@res) retrieve @server\n|   |   $href = @server.href\n|   |   $res = to_object(@res)\n|   end\n|   $final_state = @server.state\n|   if $final_state == \"operational\"\n|   |   $res = to_object(@server)\n|   |   $fields = to_json($res[\"details\"][0])\n|   |   @server = rs_cm.get({ \"href\": @server.href })\n|   elsif true|   |   $server_name = @server.name\n|   |   raise \"Failed to provision server. Expected state 'operational' but got '\" + $final_state + \"' for server: \" + $server_name + \" at href: \" + $href\n|   end\nend",
					"parameters": [],
					"application": "cwfconsole",
					"created_by": {
						"email": "support@rightscale.com",
						"id": 0,
						"name": "Terraform"
					},
					"created_at": "2018-05-25T15:18:32.054Z",
					"updated_at": "2018-05-25T15:18:33.978Z",
					"finished_at": "2018-05-25T15:19:54.382Z",
					"status": "completed",
					"links": {
						"tasks": {
							"href": "/accounts/62656/processes/5b082948a17cac6ee9ece729/tasks"
						}
					}
				}`
				switch testCase {
				case "createServer":
					retries = retries - 1
					switch retries {
					case 0:
						response = response_completed
					case 1:
						response = response_completed_no_outputs
					default:
						response = response_running
					}
				case "createServerNoOutputs":
					response = response_completed_no_outputs
				}
			case fmt.Sprintf("/cwf/v1/accounts/%d/processes/5b11716f1c02882cf0fdaa84", ms.projectID):
				response = `
				{
					"id": "5b11716f1c02882cf0fdaa84",
					"href": "/accounts/62656/processes/5b11716f1c02882cf0fdaa84",
					"name": "0ot51etbjt34l",
					"tasks": [
						{
							"id": "5b11716f1c02882cf0fdaa83",
							"href": "/accounts/62656/tasks/5b11716f1c02882cf0fdaa83",
							"name": "/root",
							"progress": {
								"percent": 100,
								"summary": ""
							},
							"status": "completed",
							"created_at": "2018-06-01T16:16:47.848Z",
							"updated_at": "2018-06-01T16:16:47.848Z",
							"finished_at": "2018-06-01T16:16:53.012Z"
						}
					],
					"outputs": [
						{
							"name": "$res",
							"value": {
								"kind": "number",
								"value": 35
							}
						}
					],
					"references": [],
					"variables": [],
					"source": "define main() return $res do\n\tsub timeout: 1h do\n\t\t$res = 23 + 12\n\tend\nend",
					"main": "define main() return $res do\n|   sub timeout: \"1h\" do\n|   |   $res = 23 + 12\n|   end\nend",
					"parameters": [],
					"application": "cwfconsole",
					"created_by": {
						"email": "user@rightscale.com",
						"id": 0,
						"name": "John Terraformer via Terraform"
					},
					"created_at": "2018-06-01T16:16:47.848Z",
					"updated_at": "2018-06-01T16:16:50.925Z",
					"finished_at": "2018-06-01T16:16:53.016Z",
					"status": "completed",
					"links": {
						"tasks": {
							"href": "/accounts/62656/processes/5b11716f1c02882cf0fdaa84/tasks"
						}
					}
				}`
			case "/api/users/11111":
				response = `
				{
					"email": "user@example.com",
					"first_name": "John",
					"last_name": "Terraformer",
					"login_name": "octocat",
					"company": "RightScale",
					"phone": "00000000",
					"timezone_name": "UTC",
					"created_at": "2014/08/12 17:18:54 +0000",
					"updated_at": "2018/05/29 09:20:32 +0000",
					"links": [
						{
							"rel": "self",
							"href": "/api/users/111111"
						}
					],
					"actions": []
				}`
			default:
				panic(fmt.Errorf("Mock Server received unknown PATH: %s", request.URL.Path))
				return
			}
			if len(response) > 0 {
				h := writer.Header()
				h.Set("Content-Type", "application/json")
				writer.Write([]byte(response))
			}
		}))

	ms.or = rshosts
	ms.ohi = httpclient.Insecure
	httpclient.Insecure = true
	rshosts = []string{ms.service.URL}

	c, err := New("mytoken", 62656)
	if err != nil {
		ms.close(t)
		t.Fatal(err)
	}
	return c
}

func (ms *mockServer) close(t *testing.T) {
	if ms.service != nil {
		rshosts = ms.or
		httpclient.Insecure = ms.ohi
		ms.service.Close()
		ms.service = nil
	}
}

func TestRunProcess(t *testing.T) {
	var ms mockServer
	c := ms.launch(t, "runProcess")
	defer ms.close(t)

	source := `define main() return $res do
	$res = 11 + 31
end
`
	process, err := c.RunProcess(source, nil)
	if err != nil {
		t.Errorf("got error %q, expected none", err)
		return
	}
	if process.Outputs["$res"] != "42" {
		t.Errorf("got $res equal %s, expected 42", process.Outputs["$res"])
	}
}

func TestGet(t *testing.T) {
	client, _ := New(validToken(t), validProjectID(t))
	l := Locator{
		Href:         "/api/deployments/936965004",
		Namespace:    "rs_cm",
		Type:         "",
		ActionParams: nil,
	}
	_, _ = client.Get(&l)
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

func TestCreateServer(t *testing.T) {
	var ms mockServer
	client := ms.launch(t, "createServer")
	defer ms.close(t)

	fields := `
	{
		"server": {
			"deployment_href": "/api/deployments/936965004",
			"instance": {
				"associate_public_ip_address": true,
				"cloud_href": "/api/clouds/1",
				"datacenter_href": "",
				"deployment_href": "",
				"image_href": "/api/clouds/1/images/E0HCVNHNAV8KK",
				"instance_type_href": "/api/clouds/1/instance_types/9K1AU4K4RCBU4",
				"ip_forwarding_enabled": false,
				"kernel_image_href": "",
				"name": "terraform-test-instance-7hcxelcntc-29ley7ipse",
				"placement_group_href": "",
				"ramdisk_image_href": "",
				"security_group_hrefs": [],
				"server_template_href": "/api/server_templates/402254004",
				"ssh_key_href": "",
				"subnet_hrefs": [
					"/api/clouds/1/subnets/52NUHI2B8LVH1"
				],
				"user_data": ""
			},
			"name": "terraform-test-server-7hcxelcntc-8k8ae2u7ia"
		}
	}`

	var fs Fields
	err := json.Unmarshal([]byte(fields), &fs)
	if err != nil {
		t.Errorf("got error %q, expected none", err)
		return
	}
	_, err = client.CreateServer("rs_cm", "servers", fs)
	if err != nil {
		t.Errorf("got error %q, expected none", err)
		return
	}

	ms.close(t)
	client = ms.launch(t, "createServerNoOutputs")
	defer ms.close(t)
	_, err = client.CreateServer("rs_cm", "servers", fs)
	ee := "no Outputs received from your CWF process, check your return clause"
	if err.Error() != ee {
		t.Errorf("expected error (%s) not received", ee)
		return
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

func TestRunRCLWithDefinitions(t *testing.T) {
	var ms mockServer
	cl := ms.launch(t, "runRCLWithDefinitions")
	defer ms.close(t)

	var rcl = `
	$res = 23 + 12
	`

	res, err := cl.(*client).runRCLWithDefinitions(rcl, "", "$res")
	if err != nil {
		t.Fatal(err)
	}

	ev := "35"
	if res["$res"] != ev {
		t.Errorf("got result %s, expected %s", res["$res"], ev)
	}
}

func TestGetProcess(t *testing.T) {
	if os.Getenv("RSC_NOMOCK") != "" {
		// This only works if using MOCK SERVER
		return
	}
	var ms mockServer
	cl := ms.launch(t, "getProcess")
	defer ms.close(t)

	res, err := cl.GetProcess(fmt.Sprintf("/accounts/%d/processes/5b11716f1c02882cf0fdaa84", ms.projectID))
	if err != nil {
		t.Fatal(err)
	}

	ev := "completed"
	if res.Status != ev {
		t.Errorf("got status %s, expected %s", res.Status, ev)
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

func TestUser(t *testing.T) {
	var ms mockServer
	c := ms.launch(t, "TestUser")
	defer ms.close(t)

	u, err := c.GetUser()
	if err != nil {
		t.Fatal(err)
	}

	tu := userString(u)
	if os.Getenv("RSC_NOMOCK") == "" {
		if tu != "John Terraformer via Terraform" {
			t.Errorf("got wrong user (%s)", tu)
		}
	} else {
		if !strings.HasSuffix(tu, " via Terraform") {
			t.Errorf("user string doesn't end with 'via Terraform' (%s)", tu)
		}
		if len(strings.TrimSpace(strings.TrimSuffix(tu, "via Terraform"))) == 0 {
			t.Errorf("user name contains only spaces")
		}
	}

}
