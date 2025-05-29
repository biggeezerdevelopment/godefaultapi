package godefaultapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// ContentType represents the supported content types
type ContentType string

const (
	// ContentTypeJSON represents application/json
	ContentTypeJSON ContentType = "application/json"
	// ContentTypeXML represents application/xml
	ContentTypeXML ContentType = "application/xml"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	// HeaderName is the name of the response header that contains rate limit information
	HeaderName string
	// MaxRetries is the maximum number of retries before giving up
	MaxRetries int
	// DefaultWaitTime is the time to wait if no rate limit header is found
	DefaultWaitTime time.Duration
}

// DefaultRateLimitConfig returns a default rate limit configuration
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		HeaderName:      "X-RateLimit-Reset",
		MaxRetries:      3,
		DefaultWaitTime: 5 * time.Second,
	}
}

// Client represents the API client
type Client struct {
	baseURL         string
	httpClient      *http.Client
	headers         map[string]string
	requestType     ContentType
	responseType    ContentType
	rateLimitConfig *RateLimitConfig
}

// NewClient creates a new API client with default configuration
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:      baseURL,
		requestType:  ContentTypeJSON,
		responseType: ContentTypeXML,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		headers:         make(map[string]string),
		rateLimitConfig: DefaultRateLimitConfig(),
	}
}

// SetRequestType sets the content type for requests
func (c *Client) SetRequestType(contentType ContentType) {
	c.requestType = contentType
}

// SetResponseType sets the content type for responses
func (c *Client) SetResponseType(contentType ContentType) {
	c.responseType = contentType
}

// SetRateLimitConfig sets the rate limiting configuration
func (c *Client) SetRateLimitConfig(config *RateLimitConfig) {
	c.rateLimitConfig = config
}

// SetContentType sets the content type for requests
func (c *Client) SetContentType(contentType ContentType) {
	c.requestType = contentType
}

// SetBearerToken sets the Authorization header with a Bearer token
func (c *Client) SetBearerToken(token string) {
	c.headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
}

// SetBasicAuth sets the Authorization header with Basic authentication
func (c *Client) SetBasicAuth(username, password string) {
	auth := username + ":" + password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	c.headers["Authorization"] = fmt.Sprintf("Basic %s", encodedAuth)
}

// SetHeader sets a custom header
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, body, result interface{}) error {
	var reqBody []byte
	if body != nil {
		if b, ok := body.([]byte); ok {
			reqBody = b
		} else {
			return fmt.Errorf("body must be []byte")
		}
	}
	return c.doRequest(ctx, http.MethodGet, path, reqBody, result)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	var reqBody []byte
	if body != nil {
		if b, ok := body.([]byte); ok {
			reqBody = b
		} else {
			return fmt.Errorf("body must be []byte")
		}
	}
	return c.doRequest(ctx, http.MethodPost, path, reqBody, result)
}

// doRequest performs the actual HTTP request with rate limiting support
func (c *Client) doRequest(ctx context.Context, method, path string, body []byte, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set content type headers
	req.Header.Set("Content-Type", string(c.requestType))
	req.Header.Set("Accept", string(c.responseType))

	// Set custom headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Retry loop for rate limiting
	var resp *http.Response
	for retry := 0; retry <= c.rateLimitConfig.MaxRetries; retry++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("error performing request: %w", err)
		}

		// Check for rate limit header
		resetTime := resp.Header.Get(c.rateLimitConfig.HeaderName)
		if resetTime != "" {
			// Parse the reset time
			resetUnix, err := strconv.ParseInt(resetTime, 10, 64)
			if err != nil {
				// If we can't parse the reset time, use default wait time
				time.Sleep(c.rateLimitConfig.DefaultWaitTime)
				continue
			}

			// Calculate wait time
			waitTime := time.Until(time.Unix(resetUnix, 0))
			if waitTime > 0 {
				// Wait for the specified time
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(waitTime):
					// Continue with retry
					continue
				}
			}
		}

		// If we get here, either there was no rate limit or we've waited
		break
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	//fmt.Println("Response body: ", string(respBody))

	if result != nil {
		switch c.responseType {
		case ContentTypeJSON:
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("error decoding JSON response: %w", err)
			}
		case ContentTypeXML:
			if err := xml.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("error decoding XML response: %w", err)
			}
		default:
			return fmt.Errorf("unsupported response content type: %s", c.responseType)
		}
	}

	return nil
}
