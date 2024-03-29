package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	errInvalidMethod error = fmt.Errorf("invalid method")
	// HTTPClient is the client Direct commands use
	HTTPClient = &http.Client{
		Timeout: 30 * time.Second,
	}
)

// Direct is a command where each component (method, headers, etc) are separate
type Direct struct {
	Client HTTPClienter `json:"-"`
	Method string       `json:"method,omitempty"`
	Body   string       `json:"body,omitempty"`
	Web
}

// HTTPClienter allows for mocking the client in tests
type HTTPClienter interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	// GET is a GET request
	GET = "GET"
	// POST is a POST request
	POST = "POST"
	// PUT is a PUT request
	PUT = "PUT"
	// TODO: more methods?
)

// Valid validates a Direct command
func (d *Direct) Valid() error {
	switch strings.ToUpper(d.Method) {
	case GET, POST, PUT:
	default:
		return errInvalidMethod
	}

	return d.Web.Valid()
}

// DeParameterize replaces parameters in the url, headers, and body
func (d *Direct) DeParameterize(params []string) (err error) {
	for i, p := range params {
		d.Body = replaceParameter(d.Body, i+1, p)
	}
	if err = d.Web.DeParameterize(params); err != nil {
		return err
	}
	if strings.Contains(d.Body, parameter) {
		return errUnhandledParams
	}
	return
}

// Run executes a command directly
func (d *Direct) Run() ([]byte, error) {
	var bdy io.Reader
	if d.Body != "" {
		bdy = strings.NewReader(d.Body)
	}

	req, err := http.NewRequest(strings.ToUpper(d.Method), d.URL, bdy)
	if err != nil {
		return nil, err
	}

	// Set our User-Agent by default. Will be overridden with a user's custom User-Agent
	// TODO: don't hardcode this
	req.Header.Set("User-Agent", "Data Connector for Google Sheets (https://github.com/brentadamson/dataconnector)")
	for _, header := range d.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
