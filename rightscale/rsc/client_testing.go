package rsc

type (
	// TestClient is a mock of the Client.
	TestClient struct {
		*TestClientExpectation
	}
)

// Make sure TestClient implements Client interface
var _ Client = (*TestClient)(nil)

// NewTestClient creates a Client which uses expectations set by the tests to
// implement the behavior.
func NewTestClient() *TestClient {
	return &TestClient{NewTestClientExpectation()}
}

// Create runs any preset expectation.
func (c *TestClient) Create(namespace, typ string, fields Fields) (*Resource, error) {
	if e := c.Expectation("Create"); e != nil {
		return e.(func(string, string, Fields) (*Resource, error))(namespace, typ, fields)
	}
	return nil, nil
}

// CreateServer runs any preset expectation.
func (c *TestClient) CreateServer(namespace, typ string, fields Fields) (*Resource, error) {
	if e := c.Expectation("CreateServer"); e != nil {
		return e.(func(string, string, Fields) (*Resource, error))(namespace, typ, fields)
	}
	return nil, nil
}

// List runs any preset expectation.
func (c *TestClient) List(l *Locator, link string, filters Fields) ([]*Resource, error) {
	if e := c.Expectation("List"); e != nil {
		return e.(func(*Locator, string, Fields) ([]*Resource, error))(l, link, filters)
	}
	return nil, nil
}

// Get runs any preset expectation.
func (c *TestClient) Get(l *Locator) (*Resource, error) {
	if e := c.Expectation("Get"); e != nil {
		return e.(func(*Locator) (*Resource, error))(l)
	}
	return nil, nil
}

// Update runs any preset expectation.
func (c *TestClient) Update(l *Locator, fields Fields) error {
	if e := c.Expectation("Update"); e != nil {
		return e.(func(*Locator, Fields) error)(l, fields)
	}
	return nil
}

// Delete runs any preset expectation.
func (c *TestClient) Delete(l *Locator) error {
	if e := c.Expectation("Delete"); e != nil {
		return e.(func(*Locator) error)(l)
	}
	return nil
}

// RunProcess runs any preset expectation.
func (c *TestClient) RunProcess(rcl string, parameters []*Parameter) (*Process, error) {
	if e := c.Expectation("RunProcess"); e != nil {
		return e.(func(string, []*Parameter) (*Process, error))(rcl, parameters)
	}
	return nil, nil
}

// GetProcess runs any preset expectation.
func (c *TestClient) GetProcess(href string) (*Process, error) {
	if e := c.Expectation("GetProcess"); e != nil {
		return e.(func(string) (*Process, error))(href)
	}
	return nil, nil
}

// DeleteProcess runs any preset expectation.
func (c *TestClient) DeleteProcess(href string) error {
	if e := c.Expectation("DeleteProcess"); e != nil {
		return e.(func(string) error)(href)
	}
	return nil
}

// GetUser returns the user's information (name, surname, email, company, ...)
// The user is the one that generated the RefreshToken provided to authenticate
// in RightScale
func (c *TestClient) GetUser() (map[string]interface{}, error) {
	if e := c.Expectation("GetUser"); e != nil {
		return e.(func() (map[string]interface{}, error))()
	}
	return nil, nil
}
