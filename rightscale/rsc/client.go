package rsc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	log15 "github.com/inconshreveable/log15"
	"github.com/rightscale/rsc/httpclient"
	rsclog "github.com/rightscale/rsc/log"
	"github.com/rightscale/rsc/rsapi"
)

var (
	// ErrNotFound is returned by Get and Update if no resource with the
	// given locator could be found.
	ErrNotFound = fmt.Errorf("resource not found")

	// rshosts lists the RightScale API hostnames.
	rshosts = []string{"us-3.rightscale.com", "us-4.rightscale.com"}
)

type (
	// Client is the RightScale API client.
	Client interface {
		// Create creates a new resource given a namespace, type name
		// and field values.
		Create(string, string, Fields) (*Resource, error)
		// CreateServer creates a new server resource given a namespace, type name
		// and field values.
		CreateServer(string, string, Fields) (*Resource, error)
		// List lists resources given a root resource locator, an
		// optional link to nested resources and optional filters. If
		// the root locator contains a Href then a link must be provided
		// and List returns the resources retrieved by following the
		// link. If the root locator does not contain a href then no
		// link can be provided and List returns the (top level)
		// resources of the given locator type. In both cases the result
		// may be filtered using the last argument.
		List(*Locator, string, Fields) ([]*Resource, error)
		// Get retrieves a resource given its locator. Get returns
		// ErrNotFound if no resource could be found for the given
		// locator.
		Get(*Locator) (*Resource, error)
		// Update updates the fields of a resource. Update returns
		// ErrNotFound if no resource could be found for the given
		// locator.
		Update(*Locator, Fields) error
		// Delete deletes an existing resource. Delete does nothing if
		// the resource cannot be found.
		Delete(*Locator) error
		// RunProcess runs a Cloud Workflow process initialized with the
		// given source code. The process runs synchronously and
		// RunProcess returns after it completes or fails. The RCL code
		// must define a definition called 'main' that accepts the given
		// parameter values.
		RunProcess(source string, parameters []*Parameter) (*Process, error)
		// GetProcess retrieves the process with the given href.
		GetProcess(href string) (*Process, error)
		// DeleteProcess deletes the process with the given href.
		DeleteProcess(href string) error
		// GetUser returns the user's information (name, surname, email, company, ...)
		GetUser() (map[string]interface{}, error)
	}

	// Resource represents a resource managed by the RightScale platform.
	Resource struct {
		// Locator is the resource locator.
		Locator *Locator
		// Fields lists the resource fields.
		Fields Fields
	}

	// Locator contains the information needed to manage a cloud resource
	// via the RightScale APIs.
	Locator struct {
		// Href is the resource path as defined by the underlying
		// service provider.
		Href string
		// Namespace identifies the service that exposes the resource.
		// The value can be one of the RightScale built-in namespaces:
		// "rs_cm", "rs_ss" or "rs_ca" or the name of a RightScale
		// plugin.
		Namespace string
		// Type is the name of the resource type scoped by the
		// namespace, e.g. "servers".
		Type string
		// ActionParams allows for the passing of arbitrary extra parameters
		// to the RightScale APIs.  These extra parameters are generally scoped
		// to a given namespace and resource.
		ActionParams map[string]string
	}

	// Fields represent arbitrary resource fields as consumed by the
	// underlying API.
	Fields map[string]interface{}

	// Parameter describes a RCL definition parameter value or output.
	Parameter struct {
		// Kind is the kind of parameter.
		Kind ParameterKind `json:"kind"`

		// Value is the parameter value. The mapping of parameter kind
		// to Go type is as follows:
		//
		//    Parameter Kind     | Go type
		//    ----------------+-----------------------------------------------
		//    KindString      | string
		//    KindNumber      | [u]int, [u]int32, [u]int64, float32 or float64
		//    KindBool        | bool
		//    KindDateTime    | string (RFC3339 time value)
		//    KindDuration    | string (RCL duration, e.g. "1h1s")
		//    KindNull        | nil
		//    KindArray       | []*Parameter
		//    KindObject      | map[string]interface{}
		//    KindCollection  | map[string]interface{}
		//    KindDeclaration | map[string]interface{}
		//
		// The map values for Parameter strict with kind:
		//
		//    - KindObject must be Parameter structs.
		//    - KindCollection must be map[string]interface{} with keys
		//      'namespace', 'type', 'hrefs' and 'details'.
		//    - KindDeclaration must be map[string]interface{} with keys
		//      'namespace', 'type' and 'fields'.
		//
		Value interface{} `json:"value"`
	}

	// ParameterKind is the RCL definition parameter kind enum.
	ParameterKind string

	// Process represents a Cloud Workflow process.
	Process struct {
		// Href is the process API resource href.
		Href string
		// Outputs lists the process outputs.
		Outputs map[string]interface{}
		// Status is the process status, one of "completed", "failed",
		// "canceled" or "aborted".
		Status string
		// Error is a synthesized error constructed when the process
		// fails. It may be ErrNotFound in case the process failed due
		// to a "ResourceNotFound" RightScale API response.
		Error error
	}

	// client is the Client interface implementation.
	client struct {
		// APIToken is the token used to authenticate RightScale API
		// requests.
		APIToken string
		// ProjectID is the id of the RightScale project (a.k.a.
		// account)
		ProjectID int

		rs *rsapi.API

		user map[string]interface{}
	}
)

const (
	KindString      ParameterKind = "string"
	KindNumber      ParameterKind = "number"
	KindBool        ParameterKind = "bool"
	KindNull        ParameterKind = "null"
	KindArray       ParameterKind = "array"
	KindObject      ParameterKind = "object"
	KindCollection  ParameterKind = "collection"
	KindDeclaration ParameterKind = "declaration"
)

// New attempts to auth against all the RightScale hosts and initializes the
//  RightScale client on success.
func New(token string, projectID int) (Client, error) {
	if strings.ToUpper(os.Getenv("TF_LOG")) == "TRACE" {
		// Shows network dumps
		httpclient.DumpFormat = httpclient.Debug // Add '| httpclient.Verbose' to see auth headers
		// Links rsc's log15 with TF's log
		rsclog.Logger.SetHandler(log15.FuncHandler(func(r *log15.Record) error {
			log.Printf("%s", string(log15.LogfmtFormat().Format(r)))
			return nil
		}))
	}
	auth := rsapi.NewOAuthAuthenticator(token, projectID)
	for _, host := range rshosts {
		rs := rsapi.New(host, auth)
		if err := rs.CanAuthenticate(); err == nil {
			req, err := rs.BuildHTTPRequest("GET", "/api/sessions/accounts", "1.5", nil, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to make session accounts request: %s", err)
			}
			resp, err := rs.PerformRequest(req)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve accounts: %s", err)
			}
			as, err := rs.LoadResponse(resp)
			if err != nil {
				return nil, fmt.Errorf("failed to load session accounts: %s", err)
			}
			if err := checkProject(as.([]interface{}), projectID); err != nil {
				return nil, err
			}
			return &client{
				APIToken:  token,
				ProjectID: projectID,
				rs:        rs,
			}, nil
		}
	}
	return nil, fmt.Errorf("failed to authenticate")
}

// List retrieves the list of resources pointed to by l optionally filtering the
// results with the given filters. The supported filters differ depending on
// the underlying resource, refer to the RightScale API 1.5 docs for details on
// the CM resources.
//
// List returns an empty slice if there is no resource for the given locator and
// filters.
func (rsc *client) List(l *Locator, link string, filters Fields) ([]*Resource, error) {
	if l.Namespace == "" {
		return nil, fmt.Errorf("resource locator is invalid: namespace is missing")
	}
	// params can be filters, views, etc.
	var params string
	{
		options := make(map[string]interface{})

		// possible api resource params eg 'view'
		if len(l.ActionParams) > 0 {
			for k, v := range l.ActionParams {
				options[k] = v
			}
		}

		// possible api resource filters eg 'name'
		if len(filters) > 0 {
			for k, v := range filters {
				options[k] = v
			}
		}

		// marshal and convert to string
		f, err := json.Marshal(options)
		if err != nil {
			return nil, fmt.Errorf("invalid list parameters: %s", err)
		}
		params = string(f)
	}

	var prefix string
	{
		if l.Href != "" {
			if link == "" {
				return nil, fmt.Errorf("cannot list nested resources: missing link")
			}
			prefix = fmt.Sprintf("@res = %s.get(href: %q).%s(%s)",
				l.Namespace, l.Href, link, params)
		} else {
			if l.Type == "" {
				return nil, fmt.Errorf("resource locator is invalid: type is missing")
			}
			prefix = fmt.Sprintf("@res = %s.%s.get(%s)",
				l.Namespace, l.Type, params)
		}
	}

	rcl := prefix + `
	$res    = to_object(@res)
	$hrefs  = to_json($res["hrefs"])
	$fields = to_json($res["details"])
	$type   = $res["type"]`

	outputs, err := rsc.runRCL(rcl, "$hrefs", "$fields", "$type")
	if err != nil {
		return nil, err
	}

	var (
		hrefs   []string
		details []map[string]interface{}
	)
	{
		err := json.Unmarshal([]byte(outputs["$hrefs"].(string)), &hrefs)
		if err != nil {
			return nil, fmt.Errorf("invalid list hrefs: %s", err)
		}
		err = json.Unmarshal([]byte(outputs["$fields"].(string)), &details)
		if err != nil {
			return nil, fmt.Errorf("invalid list fields: %s", err)
		}
	}
	typ := outputs["$type"].(string)
	res := make([]*Resource, len(hrefs))
	for i, href := range hrefs {
		loc := Locator{Namespace: l.Namespace, Type: typ, Href: href}
		res[i] = &Resource{Locator: &loc, Fields: details[i]}
	}

	return res, nil
}

// Get retrieves the details of the resource pointed to by l.
// The field 'Type' of the given Locator may be omitted.
//
// Get returns nil if there is no resource for the given locator.
func (rsc *client) Get(l *Locator) (*Resource, error) {
	if l.Namespace == "" {
		return nil, fmt.Errorf("resource locator is invalid: namespace is missing")
	}
	if l.Href == "" {
		return nil, fmt.Errorf("resource locator is invalid: href is missing")
	}

	// params can be views, etc.
	var params string

	{
		options := make(map[string]interface{})

		// for get Locator has Href - add href to options map
		options["href"] = l.Href

		// possible api resource params eg 'view' - add if we have any
		if len(l.ActionParams) > 0 {
			for k, v := range l.ActionParams {
				options[k] = v
			}
		}

		// marshal and convert to string
		f, err := json.Marshal(options)
		if err != nil {
			return nil, fmt.Errorf("invalid get parameters: %s", err)
		}
		params = string(f)
	}

	// construct rcl for get call
	prefix := fmt.Sprintf(`@res = %s.get(%s)`, l.Namespace, params)
	rcl := prefix + `
	$res = to_object(@res)
	$fields = to_json($res["details"][0])
	$type = $res["type"]`

	outputs, err := rsc.runRCL(rcl, "$fields", "$type")
	if err != nil {
		return nil, err
	}

	var fields Fields
	fs := outputs["$fields"].(string)
	err = json.Unmarshal([]byte(fs), &fields)
	if err != nil {
		return nil, err
	}
	typ := outputs["$type"].(string)
	loc := Locator{Namespace: l.Namespace, Type: typ, Href: l.Href}

	return &Resource{Locator: &loc, Fields: fields}, nil
}

// Update overwrites the fields of the resource.
// The field 'Type' of the resource Locator may be omitted.
func (rsc *client) Update(l *Locator, fields Fields) error {
	if l.Namespace == "" {
		return fmt.Errorf("resource locator is invalid: namespace is missing")
	}
	if l.Href == "" {
		return fmt.Errorf("resource locator is invalid: href is missing")
	}

	js, err := json.Marshal(fields.onlyPopulated())
	if err != nil {
		return err
	}
	rcl := fmt.Sprintf(`
	@resource = %s.get(href: "%s")
	@resource.update(%s)
	`, l.Namespace, l.Href, js)

	_, err = rsc.runRCL(rcl)
	return err
}

// Create creates the given resource. The "Href" field of the resource locator
// should not be set on input, it is set in the result.
func (rsc *client) Create(namespace, typ string, fields Fields) (*Resource, error) {
	m := map[string]interface{}{
		"namespace": namespace,
		"type":      typ,
		"fields":    fields.onlyPopulated(),
	}
	js, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	rcl := fmt.Sprintf(`
	@res = %s
	provision(@res)
	$href   = @res.href
	$res    = to_object(@res)
	$fields = to_json($res["details"][0])
	`, js)

	outputs, err := rsc.runRCL(rcl, "$href", "$fields")
	if err != nil {
		return nil, err
	}

	var ofields Fields
	err = json.Unmarshal([]byte(outputs["$fields"].(string)), &ofields)
	if err != nil {
		return nil, err
	}

	loc := Locator{
		Namespace: namespace,
		Type:      typ,
		Href:      outputs["$href"].(string),
	}
	return &Resource{Locator: &loc, Fields: ofields}, nil
}

// CreateServer creates the given resource. The "Href" field of the resource locator
// should not be set on input, it is set in the result.
func (rsc *client) CreateServer(namespace, typ string, fields Fields) (*Resource, error) {
	m := map[string]interface{}{
		"namespace": namespace,
		"type":      typ,
		"fields":    fields.onlyPopulated(),
	}
	js, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	serverSourceRcl := fmt.Sprintf(`define main() return $href, $fields do
		$href = ""
		@server = rs_cm.servers.empty()
		sub timeout: 1h do
			$res = %s
			@res = $res
			call server_create_tf(@res) retrieve @server
			if $res["fields"]["server"]["tags"]
				$tags = $res["fields"]["server"]["tags"]
				$href = [@server.href]
				call resource_add_tags($href, $tags)
			end
			call tf_server_wait_for_provision(@server) retrieve @server
			$res = to_object(@res)
			$final_state = @server.state
			if $final_state == "operational"
				$href = @server.href
				$res = to_object(@server)
				$fields = to_json($res["details"][0])
				@server = rs_cm.get(href: @server.href)
			else
				$server_name = @server.name
				$href = @server.href
				raise "Failed to provision server. Expected state 'operational' but got '" + $final_state + "' for server: " + $server_name + " at href: " + $href
			end
		end
	end

	define resource_add_tags($href, $tags) do
		rs_cm.tags.multi_add(resource_hrefs: $href, tags: $tags)
	end

	# create the server object
	define server_create_tf(@res) return @server do
		# use RS canned provision to create
		call rs__cwf_simple_provision(@res) retrieve @server
	end

	# custom launch object that does not auto-delete on errors
	define tf_server_wait_for_provision(@server) return @server do
		$server_name = to_s(@server.name)
		sub on_error: tf_server_handle_launch_failure(@server) do
			@server.launch()
		end
		$final_state = "launching"
		# use RS canned logic to capture launching server state
		sub on_error: rs__cwf_skip_any_error() do
			sleep_until @server.state =~ "^(operational|stranded|stranded in booting|stopped|terminated|inactive|error)$"
			$final_state = @server.state
		end
	end

	# spit out error from launch call
	define tf_server_handle_launch_failure(@server) do
		$server_name = @server.name
		if $_errors && $_errors[0] && $_errors[0]["response"]
			raise "Error trying to launch server (" + $server_name + "): " + $_errors[0]["response"]["body"]
		else
			raise "Error trying to launch server (" + $server_name + ")"
		end
	end`, js)

	ts := time.Now().Add(-time.Second * 15)
	p, err := rsc.RunProcess(serverSourceRcl, nil)
	if err != nil {
		return nil, err
	}

	outputs := p.Outputs
	loc := Locator{
		Namespace: namespace,
		Type:      typ,
		Href:      outputs["$href"].(string),
	}

	if p.Status != "completed" {
		te := time.Now().Add(time.Second * 15)
		e := fmt.Errorf(
			`unexpected process status %q. Error: %s.

Check your account audit entries for more details with:
rsc --refreshToken <refreshToken> --pp --account %d --host %s cm15 index /api/audit_entries  'start_date=%s' 'end_date=%s' 'limit=1000'`,
			p.Status,
			p.Error,
			rsc.ProjectID,
			rsc.rs.Host,
			ts.Format("2006/01/02 15:04:05 -0700"),
			te.Format("2006/01/02 15:04:05 -0700"),
		)
		return &Resource{Locator: &loc, Fields: nil}, e
	}

	var ofields Fields
	err = json.Unmarshal([]byte(outputs["$fields"].(string)), &ofields)
	if err != nil {
		return nil, err
	}

	return &Resource{Locator: &loc, Fields: ofields}, nil
}

// runRCLWithDefinitions provides a convenient method for running the given RCL code
// synchronously including with any definitions. It returns the outputs with the given variable or reference
// names.  It executes from a simply constructed 'main' so this may not be sufficient depending on the resource.
func (rsc *client) runRCLWithDefinitions(rcl string, defs string, outputs ...string) (map[string]interface{}, error) {
	source := "define main() "
	if len(outputs) > 0 {
		source += "return " + strings.Join(outputs, ", ") + " "
	}
	rcl = strings.Trim(rcl, "\n\t")
	rcl = strings.Replace(rcl, "\t", "\t\t", -1)
	source += "do\n\tsub timeout: 1h do\n\t\t" + rcl + "\n\tend\nend"
	if len(defs) > 0 {
		source += "\n" + defs
	}
	p, err := rsc.RunProcess(source, nil)
	if err != nil {
		return nil, err
	}
	if p.Status != "completed" {
		return p.Outputs, fmt.Errorf("unexpected process status %q. Error: %s", p.Status, p.Error)
	}
	return p.Outputs, nil
}

// Delete deletes the given resource.
// Only the Href field or res needs to be initialized.
func (rsc *client) Delete(l *Locator) error {
	if l.Namespace == "" {
		return fmt.Errorf("resource locator is invalid: namespace is missing")
	}
	if l.Href == "" {
		return fmt.Errorf("resource locator is invalid: href is missing")
	}
	rcl := fmt.Sprintf(`
	@res = %s.get(href: %q)
	delete(@res)
	`,
		l.Namespace, l.Href)
	_, err := rsc.runRCL(rcl)
	return err
}

// RunProcess runs the given RCL code synchronously and returns the process outputs.
// The code must contain a 'main' definition that accepts the given parameter values.
func (rsc *client) RunProcess(source string, params []*Parameter) (*Process, error) {
	expectsOutputs, err := analyzeSource(source)
	if err != nil {
		return nil, err
	}

	var (
		projectID   = strconv.Itoa(rsc.ProjectID)
		processHref string
		processID   string
	)
	{
		u, _ := rsc.GetUser()
		payload := rsapi.APIParams{
			"source":      source,
			"main":        "main",
			"rcl_version": "2",
			"parameters":  params,
			"application": "cwfconsole",
			"created_by": map[string]interface{}{
				"id":    0,
				"name":  userString(u),
				"email": u["email"],
			},
			"refresh_token": rsc.APIToken,
		}
		res, err := rsc.requestCWF("post", "/cwf/v1/accounts/"+projectID+"/processes", nil, payload)
		if err != nil {
			return nil, err
		}
		processHref = res.(map[string]interface{})["Location"].(string)
		parts := strings.Split(processHref, "/")
		processID = parts[len(parts)-1]
	}

	var (
		process *Process
	)

	// print link to CWF console if TF_LOG=TRACE is set, mainly useful for tests
	if strings.ToUpper(os.Getenv("TF_LOG")) == "TRACE" {
		host := strings.Replace(rsc.rs.Host, "us-", "selfservice-", 1)
		fmt.Printf("RUNNING: https://%s/designer/processes/%s\n%s\n", host, processID, source)
		then := time.Now()
		defer func() {
			if process != nil {
				fmt.Printf("==> %s (in %v)\n\n", process.Status, time.Now().Sub(then))
				if err == nil {
					err = process.Error
				}
			}
			if err != nil {
				fmt.Printf("** ERROR: %s\n\n", err)
			}
		}()
	}

	timeout := time.NewTimer(1 * time.Hour)
	ticker := time.NewTicker(2 * time.Second)
	expectsOutputsTimeout := 5
	defer ticker.Stop()
	defer timeout.Stop()

	for {
		select {
		case <-ticker.C:
			var res map[string]interface{}
			{
				var r interface{}
				r, err = rsc.requestCWF("get", "/cwf/v1"+processHref, rsapi.APIParams{"view": "expanded"}, nil)
				if err == nil {
					res = r.(map[string]interface{})
				}
			}
			if err != nil {
				return nil, err
			}

			switch res["status"].(string) {
			case "not_started", "running":
				continue

			case "completed":
				process = &Process{
					Href:    processHref,
					Status:  res["status"].(string),
					Outputs: processOutputs(res),
				}

				// Keep waiting if outputs aren't yet present
				if expectsOutputs && len(process.Outputs) == 0 {
					expectsOutputsTimeout--
					if expectsOutputsTimeout == 0 {
						err = fmt.Errorf("no Outputs received from your CWF process, check your return clause")
						return nil, err
					}
					continue
				}

				return process, nil

			// capture outputs on failed just in case we had artifacts - we want to record those IDs so we don't create orphans
			case "failed":
				process = &Process{
					Href:    processHref,
					Status:  res["status"].(string),
					Outputs: processOutputs(res),
					Error:   processErrors(res),
				}

				// Keep waiting if outputs aren't yet present
				if expectsOutputs && len(process.Outputs) == 0 {
					expectsOutputsTimeout--
					if expectsOutputsTimeout == 0 {
						// Don't generate error because sometimes when failing we don't have outputs
						return process, nil
					}
					continue
				}

				return process, nil

			default:
				process = &Process{
					Href:   processHref,
					Status: res["status"].(string),
					Error:  processErrors(res),
				}
				return process, nil
			}
		case <-timeout.C:
			err = fmt.Errorf("timed out after one hour")
			return nil, err
		}
	}
}

// GetProcess returns the process with the given href.
func (rsc *client) GetProcess(href string) (*Process, error) {
	var (
		res map[string]interface{}
		err error
	)
	{
		var r interface{}
		r, err = rsc.requestCWF("get", "/cwf/v1"+href, rsapi.APIParams{"view": "expanded"}, nil)
		if err == nil {
			res = r.(map[string]interface{})
		}
	}
	if err != nil {
		return nil, err
	}
	return &Process{
		Href:    href,
		Status:  res["status"].(string),
		Outputs: processOutputs(res),
		Error:   processErrors(res),
	}, nil
}

// DeleteProcess deletes the process with the given href.
func (rsc *client) DeleteProcess(href string) error {
	_, err := rsc.requestCWF("delete", "/cwf/v1/"+href, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns the user's information (name, surname, email, company, ...)
// The user is the one that generated the RefreshToken provided to authenticate
// in RightScale
func (rsc *client) GetUser() (user map[string]interface{}, err error) {
	if rsc.user == nil {
		ui := getCurrentUserID(rsc.rs)
		if ui == "" {
			err = fmt.Errorf("Couldn't retrieve information of user from credentials")
			return nil, err
		}
		rsc.user = getUserInfo(rsc.rs, ui)
	}
	return rsc.user, nil
}

// API returns the low level RightScale API. This is not exposed by the public
// interface and is mainly intended for use by tests.
func (rsc *client) API() *rsapi.API {
	return rsc.rs
}

// runRCL provides a convenient method for running the given RCL code
// synchronously. It returns the outputs with the given variable or reference
// names. The code must not include any definition (defaults), use runRCLWithDefinitions
// or or RunProcess for including definitions.
func (rsc *client) runRCL(rcl string, outputs ...string) (map[string]interface{}, error) {
	source := "define main() "
	if len(outputs) > 0 {
		source += "return " + strings.Join(outputs, ", ") + " "
	}
	rcl = strings.Trim(rcl, "\n\t")
	rcl = strings.Replace(rcl, "\t", "\t\t", -1)
	source += "do\n\tsub timeout: 1h do\n\t\t" + rcl + "\n\tend\nend"
	ts := time.Now().Add(-time.Second * 15)
	p, err := rsc.RunProcess(source, nil)
	if err != nil {
		return nil, err
	}
	if p.Status != "completed" {
		te := time.Now().Add(time.Second * 15)
		e := fmt.Errorf(
			`unexpected process status %q. Error: %s.

Check your account audit entries for more details with:
rsc --refreshToken <refreshToken> --pp --account %d --host %s cm15 index /api/audit_entries  'start_date=%s' 'end_date=%s' 'limit=1000'`,
			p.Status,
			p.Error,
			rsc.ProjectID,
			rsc.rs.Host,
			ts.Format("2006/01/02 15:04:05 -0700"),
			te.Format("2006/01/02 15:04:05 -0700"),
		)
		return nil, e
	}
	return p.Outputs, nil
}

// requestCWF makes a request to the RightScale Cloud Workflow API.
func (rsc *client) requestCWF(method, url string, params, payload rsapi.APIParams) (interface{}, error) {
	req, err := rsc.rs.BuildHTTPRequest(strings.ToUpper(method), url, "1.0", params, payload)
	if err != nil {
		return nil, err
	}
	req.Host = strings.Replace(rsc.rs.Host, "us-", "cloud-workflow", 1)

	res, err := rsc.rs.PerformRequest(req)
	if err != nil {
		return nil, err
	}
	// simplistic handling should be enough for this one API
	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected response status code %q", res.Status)
	}

	resp, err := rsc.rs.LoadResponse(res)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// checkProject verifies that the given project ID is one of the projects listed in the
// session/accounts response 'as'. It returns nil if so, an error otherwise.
func checkProject(as []interface{}, projectID int) error {
	for _, a := range as {
		links := a.(map[string]interface{})["links"].([]interface{})
		for _, l := range links {
			href := l.(map[string]interface{})["href"].(string)
			idx := strings.LastIndex(href, "/")
			id, err := strconv.Atoi(href[idx+1:])
			if err != nil {
				return fmt.Errorf("invalid project ID %q", href[idx+1:])
			}
			if id == projectID {
				return nil
			}
		}
	}
	return fmt.Errorf("session does not give access to project %d", projectID)
}

// processOutputs is a helper function that extracts the process outputs from an
// unmarshalled "GET process" response.
func processOutputs(res interface{}) map[string]interface{} {
	outs := res.(map[string]interface{})["outputs"].([]interface{})
	outputs := make(map[string]interface{}, len(outs))
	for _, out := range outs {
		om := out.(map[string]interface{})
		v := om["value"].(map[string]interface{})["value"]
		s, ok := v.(string)
		if !ok {
			m, _ := json.Marshal(v)
			s = string(m)
		}
		outputs[om["name"].(string)] = s
	}
	return outputs
}

// processErrors is a helper function that extracts the process errors if any
// from an unmarshalled "GET process" response.
func processErrors(res interface{}) error {
	var msgs []string
	{
		ts := res.(map[string]interface{})["tasks"].([]interface{})
		for _, task := range ts {
			if err, ok := task.(map[string]interface{})["error"]; ok {
				msg := err.(map[string]interface{})["message"].(string)
				if err := rclError(msg); err != nil {
					return err
				}
				msgs = append(msgs, msg)
			}
		}
	}
	if len(msgs) == 0 {
		return nil
	}
	return errors.New(strings.Join(msgs, ", "))
}

// rclError analyzes the error message returned by runCWF and maps it to one of
// the error variables defined in this package. Right now this only looks for
// not found errors. It uses a heuristic that looks for the text
// "ResourceNotFound" as returned by the RightScale API 1.5 or the status code
// 404.
func rclError(err string) error {
	if strings.Contains(err, "ResourceNotFound") || strings.Contains(err, "status code '404'") {
		return ErrNotFound
	}
	return nil
}

// onlyPopulated takes inFields and strips out all of the unpopulated parameters.
// Unpopulated parameters sent to CWF result in empty parameter errors.
func (inFields Fields) onlyPopulated() Fields {
	outFields := Fields{}
	for k, v := range inFields {
		switch v.(type) {
		case []interface{}:
			if len(v.([]interface{})) > 0 {
				outFields[k] = v.([]interface{})
			}
		case map[string]interface{}:
			if len(v.(map[string]interface{})) > 0 {
				outFields[k] = Fields(v.(map[string]interface{})).onlyPopulated()
			}
		case Fields:
			if len(v.(Fields)) > 0 {
				outFields[k] = v.(Fields).onlyPopulated()
			}
		case string:
			if v.(string) != "" {
				outFields[k] = v
			}
		default:
			outFields[k] = v
		}
	}
	return outFields
}

// analyzeSource checks that the defition of the source is valid, returning an error if it's not
// If the definition is valid, expectsOuputs boolean indicates if the defition includes the "return" keyword,
// which indicates that output values are expected.
var validDefinition = regexp.MustCompile("^[[:blank:]\n]*define[[:blank:]]*[\\w_\\.]+[[:blank:]]*\\([@$\\w _,]*\\)[[:blank:]]+(return.*[[:blank:]]+)?do")

func analyzeSource(source string) (expectsOutputs bool, err error) {
	m := validDefinition.FindStringSubmatch(source)
	if err != nil {
		return false, fmt.Errorf("error parsing rightscale_cwf_process source definition: %s", err)
	}

	if len(m) != 2 {
		return false, fmt.Errorf("invalid rightscale_cwf_process source definition")
	}

	// if return capture group matched, expectOutputs = true
	return m[1] != "", nil
}

// returns the user's ID via /api/sessions {view: whoami} call
func getCurrentUserID(rs *rsapi.API) string {
	req, err := rs.BuildHTTPRequest("GET", "/api/sessions", "1.5", rsapi.APIParams{"view": "whoami"}, nil)
	if err != nil {
		panic(err)
	}

	resp, err := rs.PerformRequest(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to retrieve user: index returned %q", resp.Status))
	}
	ms, err := rs.LoadResponse(resp)
	if err != nil {
		panic(err)
	}
	links := ms.(map[string]interface{})["links"]
	for _, el := range links.([]interface{}) {
		var kind, value string
		for k, v := range el.(map[string]interface{}) {
			if k == "rel" {
				kind = v.(string)
			}
			if k == "href" {
				value = v.(string)
			}
		}
		if kind == "user" {
			parts := strings.Split(value, "/")
			return parts[len(parts)-1]
		}
	}
	return ""
}

// retrieves user information providing the user ID via /api/users/<ID>
func getUserInfo(rs *rsapi.API, uid string) map[string]interface{} {
	req, err := rs.BuildHTTPRequest("GET", fmt.Sprintf("/api/users/%s", uid), "1.5", nil, nil)
	if err != nil {
		panic(err)
	}

	resp, err := rs.PerformRequest(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to retrieve user: index returned %q", resp.Status))
	}
	ms, err := rs.LoadResponse(resp)
	if err != nil {
		panic(err)
	}
	return ms.(map[string]interface{})
}

// generates a string from the user's map[string]interface{}
func userString(u map[string]interface{}) string {
	return fmt.Sprintf("%s %s via Terraform", u["first_name"], u["last_name"])
}
