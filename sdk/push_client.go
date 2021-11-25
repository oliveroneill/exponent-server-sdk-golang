package expo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	// DefaultHost is the default Expo host
	DefaultHost = "https://exp.host"
	// DefaultBaseAPIURL is the default path for API requests
	DefaultBaseAPIURL = "/--/api/v2"
)

// DefaultHTTPClient is the default *http.Client for making API requests
var DefaultHTTPClient = &http.Client{}

// PushClient is an object used for making push notification requests
type PushClient struct {
	host        string
	apiURL      string
	accessToken string
	httpClient  *http.Client
	url         string
}

// ClientConfig specifies params that can optionally be specified for alternate
// Expo config and path setup when sending API requests
type ClientConfig struct {
	Host        string
	APIURL      string
	AccessToken string
	HTTPClient  *http.Client
	URL         string
}

// NewPushClient creates a new Exponent push client
// See full API docs at https://docs.getexponent.com/versions/v13.0.0/guides/push-notifications.html#http-2-api
func NewPushClient(config *ClientConfig) *PushClient {
	c := new(PushClient)
	host := DefaultHost
	apiURL := DefaultBaseAPIURL
	httpClient := DefaultHTTPClient
	accessToken := ""
	url := ""

	if config != nil {
		if config.Host != "" {
			host = config.Host
		}
		if config.APIURL != "" {
			apiURL = config.APIURL
		}
		if config.AccessToken != "" {
			accessToken = config.AccessToken
		}
		if config.HTTPClient != nil {
			httpClient = config.HTTPClient
		}
		if config.URL != "" {
			url = config.URL
		}
	}
	c.host = host
	c.apiURL = apiURL
	c.httpClient = httpClient
	c.accessToken = accessToken
	c.url = url
	return c
}

// Publish sends a single push notification
// @param push_message: A PushMessage object
// @return an array of PushResponse objects which contains the results.
// @return error if any requests failed
func (c *PushClient) Publish(message *PushMessage) (PushResponse, error) {
	responses, err := c.PublishMultiple([]PushMessage{*message})
	if err != nil {
		return PushResponse{}, err
	}
	return responses[0], nil
}

// PublishMultiple sends multiple push notifications at once
// @param push_messages: An array of PushMessage objects.
// @return an array of PushResponse objects which contains the results.
// @return error if the request failed
func (c *PushClient) PublishMultiple(messages []PushMessage) ([]PushResponse, error) {
	return c.publishInternal(messages)
}

func (c *PushClient) publishInternal(messages []PushMessage) ([]PushResponse, error) {
	// Validate the messages
	for _, message := range messages {
		if len(message.To) == 0 {
			return nil, errors.New("No recipients")
		}
		for _, recipient := range message.To {
			if recipient == "" {
				return nil, errors.New("Invalid push token")
			}
		}
	}

	url := fmt.Sprintf("%s%s/push/send", c.host, c.apiURL)
	if c.url != "" {
		url = c.url
	}

	jsonBytes, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}

	// Create request w/ body
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	// Add appropriate headers
	req.Header.Add("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.accessToken)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check that we didn't receive an invalid response
	err = checkStatus(resp)
	if err != nil {
		return nil, err
	}

	// Validate the response format first
	var r *Response
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		// The response isn't json
		return nil, err
	}
	// If there are errors with the entire request, raise an error now.
	if r.Errors != nil {
		return nil, NewPushServerError("Invalid server response", resp, r, r.Errors)
	}
	// We expect the response to have a 'data' field with the responses.
	if r.Data == nil {
		return nil, NewPushServerError("Invalid server response", resp, r, nil)
	}
	// Sanity check the response
	if len(messages) != len(r.Data) {
		message := "Mismatched response length. Expected %d receipts but only received %d"
		errorMessage := fmt.Sprintf(message, len(messages), len(r.Data))
		return nil, NewPushServerError(errorMessage, resp, r, nil)
	}
	// Add the original message to each response for reference
	for i := range r.Data {
		r.Data[i].PushMessage = messages[i]
	}
	return r.Data, nil
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return fmt.Errorf("invalid access token %s", resp.StatusCode)
	}

	return fmt.Errorf("invalid response %s", resp.Status)
}
