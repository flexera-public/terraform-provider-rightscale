package rightscale

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Provider returns the RightScale terraform provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"rightscale_api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("RIGHTSCALE_API_TOKEN", nil),
				Description: "API token used to authenticate the RightScale API requests",
			},
			"rightscale_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RIGHTSCALE_PROJECT_ID", nil),
				Description: "ID of RightScale project used to provision resources",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"rightscale_cloud":             dataSourceCloud(),
			"rightscale_credential":        dataSourceCredential(),
			"rightscale_datacenter":        dataSourceDatacenter(),
			"rightscale_deployment":        dataSourceDeployment(),
			"rightscale_image":             dataSourceImage(),
			"rightscale_instance":          dataSourceInstance(),
			"rightscale_instance_type":     dataSourceInstanceType(),
			"rightscale_multi_cloud_image": dataSourceMultiCloudImage(),
			"rightscale_network":           dataSourceNetwork(),
			"rightscale_network_gateway":   dataSourceNetworkGateway(),
			"rightscale_route_table":       dataSourceRouteTable(),
			"rightscale_security_group":    dataSourceSecurityGroup(),
			"rightscale_server":            dataSourceServer(),
			"rightscale_server_template":   dataSourceServerTemplate(),
			"rightscale_ssh_key":           dataSourceSSHKey(),
			"rightscale_subnet":            dataSourceSubnet(),
			"rightscale_volume":            dataSourceVolume(),
			"rightscale_volume_snapshot":   dataSourceVolumeSnapshot(),
			"rightscale_volume_type":       dataSourceVolumeType(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"rightscale_deployment":          resourceDeployment(),
			"rightscale_instance":            resourceInstance(),
			"rightscale_network":             resourceNetwork(),
			"rightscale_network_gateway":     resourceNetworkGateway(),
			"rightscale_route_table":         resourceRouteTable(),
			"rightscale_security_group":      resourceSecurityGroup(),
			"rightscale_security_group_rule": resourceSecurityGroupRule(),
			"rightscale_server":              resourceServer(),
			"rightscale_server_array":        resourceServerArray(),
			"rightscale_ssh_key":             resourceSSHKey(),
			"rightscale_subnet":              resourceSubnet(),
			"rightscale_cwf_process":         resourceCWFProcess(),
			"rightscale_credential":          resourceCredential(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(r *schema.ResourceData) (interface{}, error) {
	token := r.Get("rightscale_api_token").(string)
	if token == "" {
		return nil, errors.New("missing 'api_token'")
	}
	project, err := strconv.Atoi(r.Get("rightscale_project_id").(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse project ID: %s", err)
	}

	return rsc.New(token, project)
}
