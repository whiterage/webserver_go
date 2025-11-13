package worker

import (
	"context"
	"net/http"
	"time"

	"github.com/whiterage/webserver_go/pkg/models"
)

type HTTPChecker struct {
	client *http.Client
}

func NewHTTPChecker(timeout time.Duration) *HTTPChecker {
	return &HTTPChecker{client: &http.Client{Timeout: timeout}}
}

func (c *HTTPChecker) Check(ctx context.Context, url string) models.LinkStatus {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return models.LinkStatus{URL: url, Status: "error"}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return models.LinkStatus{URL: url, Status: "unavailable"}
	}
	resp.Body.Close()

	status := "available"
	if resp.StatusCode >= 400 {
		status = "unavailable"
	}

	return models.LinkStatus{
		URL:       url,
		Status:    status,
		CheckTime: time.Now().UTC(),
	}
}
