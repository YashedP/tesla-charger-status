package tesla

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	defaultAcceptHeader = "application/json"
	defaultUserAgent    = "tesla-charger-status/1.0"
	retryCount          = 3
	retryWaitTime       = 100 * time.Millisecond
	retryMaxWaitTime    = 300 * time.Millisecond
)

type FleetClient struct {
	baseURL string
}

func NewFleetClient(baseURL string) *FleetClient {
	return &FleetClient{baseURL: strings.TrimRight(baseURL, "/")}
}

func (c *FleetClient) GetChargingState(ctx context.Context, httpClient *http.Client, vin string) (string, error) {
	if vin == "" {
		return "", fmt.Errorf("vin is required")
	}
	if httpClient == nil {
		return "", fmt.Errorf("http client is required")
	}

	endpoint := fmt.Sprintf("/api/1/vehicles/%s/vehicle_data", url.PathEscape(vin))
	payload := &vehicleDataPayload{}

	resp, err := c.newRequestClient(httpClient).
		R().
		SetContext(ctx).
		SetResult(payload).
		Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("perform tesla request: %w", err)
	}

	if resp.StatusCode() >= http.StatusMultipleChoices {
		return "", fmt.Errorf(
			"tesla api status=%d body=%q",
			resp.StatusCode(),
			strings.TrimSpace(string(resp.Body())),
		)
	}

	if payload.Error != "" {
		return "", fmt.Errorf("tesla error=%q description=%q", payload.Error, payload.ErrorDescription)
	}

	state := strings.TrimSpace(payload.Response.ChargeState.ChargingState)
	if state == "" {
		return "", fmt.Errorf("missing charge_state.charging_state in response")
	}

	return state, nil
}

func (c *FleetClient) newRequestClient(httpClient *http.Client) *resty.Client {
	client := resty.NewWithClient(httpClient)
	client.SetBaseURL(c.baseURL)
	client.SetHeader("Accept", defaultAcceptHeader)
	client.SetHeader("User-Agent", defaultUserAgent)
	client.SetRetryCount(retryCount)
	client.SetRetryWaitTime(retryWaitTime)
	client.SetRetryMaxWaitTime(retryMaxWaitTime)
	client.AddRetryCondition(func(response *resty.Response, err error) bool {
		if err != nil {
			return true
		}
		if response == nil {
			return false
		}
		status := response.StatusCode()
		return status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
	})
	return client
}

type vehicleDataPayload struct {
	Response struct {
		ChargeState struct {
			ChargingState string `json:"charging_state"`
		} `json:"charge_state"`
	} `json:"response"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
