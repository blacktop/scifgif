package giphy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// A Client communicates with the Giphy API.
type Client struct {
	// APIKey is the key used for requests to the Giphy API
	APIKey string

	// Limit is the limit used for requests to the Giphy API
	Limit int

	// Rating is the rating used for requests to the Giphy API
	Rating string

	// BaseURL is the base url for Giphy API.
	BaseURL *url.URL

	// BasePath is the base path for the gifs endpoints
	BasePath string

	// User agent used for HTTP requests to Giphy API.
	UserAgent string

	// HTTP client used to communicate with the Giphy API.
	httpClient *http.Client
}

// Getopt reads environment variables.
// If not found will return a supplied default value
func Getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

// GetIntopt reads environment variables.
// If not found will return a supplied default value
func GetIntopt(name string, dfault int) int {
	if value, err := strconv.Atoi(os.Getenv(name)); err == nil {
		return value
	}
	return dfault
}

// NewClient returns a new Giphy API client.
// If no *http.Client were provided then http.DefaultClient is used.
func NewClient(httpClients ...*http.Client) *Client {
	var httpClient *http.Client

	if len(httpClients) > 0 && httpClients[0] != nil {
		httpClient = httpClients[0]
	} else {
		cloned := *http.DefaultClient
		httpClient = &cloned
	}

	c := &Client{
		APIKey: Getopt("GIPHY_API_KEY", "dc6zaTOxFJmzC"),
		Rating: Getopt("GIPHY_RATING", "PG"),
		Limit:  GetIntopt("GIPHY_LIMIT", 25),
		BaseURL: &url.URL{
			Scheme: Getopt("GIPHY_BASE_URL_SCHEME", "https"),
			Host:   Getopt("GIPHY_BASE_URL_HOST", "api.giphy.com"),
		},
		BasePath:   Getopt("GIPHY_BASE_PATH", "/v1"),
		UserAgent:  Getopt("GIPHY_USER_AGENT", "scifgif.go"),
		httpClient: httpClient,
	}

	return c
}

// NewRequest creates an API request.
func (c *Client) NewRequest(s string) (*http.Request, error) {
	rel, err := url.Parse(c.BasePath + s)
	if err != nil {
		return nil, err
	}

	q := rel.Query()
	q.Set("api_key", c.APIKey)
	q.Set("rating", c.Rating)
	rel.RawQuery = q.Encode()

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// decoded and stored in the value pointed to by v, or returned as an error if
// an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	// Make sure to close the connection after replying to this request
	req.Close = true

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	if err != nil {
		return nil, fmt.Errorf("error reading response from %s %s: %s", req.Method, req.URL.RequestURI(), err)
	}

	return resp, nil
}
