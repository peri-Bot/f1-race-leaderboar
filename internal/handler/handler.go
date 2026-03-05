package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/periBot/f1-race-leaderboard/internal/cache"
	"github.com/periBot/f1-race-leaderboard/internal/f1client"
)

// routeToEndpoint maps API Gateway route paths to F1 API endpoint keys.
var routeToEndpoint = map[string]string{
	"/standings":              f1client.EndpointDriverStandings,
	"/standings/drivers":      f1client.EndpointDriverStandings,
	"/standings/constructors": f1client.EndpointConstructorStandings,
	"/results":                f1client.EndpointLastResults,
	"/schedule":               f1client.EndpointSchedule,
}

// Handler processes API Gateway proxy requests for F1 data.
type Handler struct {
	cache    *cache.Cache
	f1Client *f1client.Client
}

// New creates a new Handler.
func New(c *cache.Cache, fc *f1client.Client) *Handler {
	return &Handler{
		cache:    c,
		f1Client: fc,
	}
}

// HandleRequest processes an API Gateway V2 HTTP request.
func (h *Handler) HandleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	path := req.RawPath
	if path == "" {
		path = req.RequestContext.HTTP.Path
	}

	endpoint, ok := routeToEndpoint[path]
	if !ok {
		return h.jsonError(http.StatusNotFound, "Unknown route: "+path), nil
	}

	// 1. Check cache
	data, hit, err := h.cache.Get(ctx, endpoint)
	if err != nil {
		log.Printf("WARN: cache get error for %s: %v", endpoint, err)
		// fall through to fetch from API
	}

	if hit {
		log.Printf("Cache HIT for %s", endpoint)
		return h.jsonOK(data, true), nil
	}

	// 2. Cache miss — fetch from external API
	log.Printf("Cache MISS for %s — fetching from API", endpoint)
	data, err = h.f1Client.Fetch(endpoint)
	if err != nil {
		log.Printf("ERROR: F1 API fetch failed for %s: %v", endpoint, err)
		return h.jsonError(http.StatusBadGateway, "Failed to fetch F1 data"), nil
	}

	// 3. Store in cache (best-effort)
	if putErr := h.cache.Put(ctx, endpoint, data); putErr != nil {
		log.Printf("WARN: cache put error for %s: %v", endpoint, putErr)
	}

	return h.jsonOK(data, false), nil
}

// jsonOK returns a 200 response with the raw JSON data.
func (h *Handler) jsonOK(data []byte, cached bool) events.APIGatewayV2HTTPResponse {
	source := "api"
	if cached {
		source = "cache"
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"X-Data-Source":               source,
		},
		Body: string(data),
	}
}

// jsonError returns an error response with a JSON body.
func (h *Handler) jsonError(status int, message string) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]string{
		"error": message,
	})

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}
}

// HealthCheck returns a simple health response (for the root path).
func (h *Handler) HealthCheck() events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]interface{}{
		"status": "ok",
		"routes": []string{
			"GET /standings",
			"GET /standings/drivers",
			"GET /standings/constructors",
			"GET /results",
			"GET /schedule",
		},
	})

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		Body: string(body),
	}
}

// HandleRequestWithHealth wraps HandleRequest to also serve a health check on /.
func (h *Handler) HandleRequestWithHealth(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	path := req.RawPath
	if path == "" {
		path = req.RequestContext.HTTP.Path
	}

	if path == "/" || path == "" {
		return h.HealthCheck(), nil
	}

	return h.HandleRequest(ctx, req)
}

// String implements fmt.Stringer for debugging.
func (h *Handler) String() string {
	return fmt.Sprintf("Handler{routes: %d}", len(routeToEndpoint))
}
