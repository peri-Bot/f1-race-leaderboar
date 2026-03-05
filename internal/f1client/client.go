package f1client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/periBot/f1-race-leaderboard/internal/models"
)

const (
	baseURL        = "https://api.jolpi.ca/ergast/f1"
	defaultTimeout = 5 * time.Second
)

// Client fetches Formula 1 data from the Jolpica Ergast-compatible API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// New creates a new F1 API client.
func New() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    baseURL,
	}
}

// NewWithHTTPClient creates a client with a custom http.Client (for testing).
func NewWithHTTPClient(hc *http.Client) *Client {
	return &Client{
		httpClient: hc,
		baseURL:    baseURL,
	}
}

// Endpoint constants used as cache keys and route identifiers.
const (
	EndpointDriverStandings      = "driver_standings"
	EndpointConstructorStandings = "constructor_standings"
	EndpointLastResults          = "last_results"
	EndpointSchedule             = "schedule"
)

// endpointPaths maps logical endpoint names to API URL paths.
var endpointPaths = map[string]string{
	EndpointDriverStandings:      "/current/driverStandings.json",
	EndpointConstructorStandings: "/current/constructorStandings.json",
	EndpointLastResults:          "/current/last/results.json",
	EndpointSchedule:             "/current.json",
}

// Fetch retrieves F1 data for the given endpoint and returns the raw JSON.
func (c *Client) Fetch(endpoint string) ([]byte, error) {
	path, ok := endpointPaths[endpoint]
	if !ok {
		return nil, fmt.Errorf("unknown endpoint: %s", endpoint)
	}

	url := c.baseURL + path

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return body, nil
}

// FetchParsed retrieves and parses the F1 API response into the typed model.
func (c *Client) FetchParsed(endpoint string) (*models.APIResponse, error) {
	body, err := c.Fetch(endpoint)
	if err != nil {
		return nil, err
	}

	var resp models.APIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response JSON: %w", err)
	}

	return &resp, nil
}
