// Package rest provides a simple HTTP and REST request builder and client.
package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"go.yaml.in/yaml/v3"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/rest/helpers"
)

// ErrEmptyResponseBody means the client receives an unexpected empty response from server
var ErrEmptyResponseBody = errors.New("empty response body")
var bufferSize = 1024

// ErrorResponse is the status code and response received from the server when an error occurs.
type ErrorResponse struct {
	StatusCode int    //  Response status code
	Message    string // Response text
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Error response from server. Status code: %v; message: %v", e.StatusCode, e.Message)
}

// Client is a simple HTTP and REST client. Create it with NewClient method.
type Client struct {
	HTTPClient    *http.Client // HTTP client, default is HTTP DefaultClient
	DefaultHeader http.Header  // Default header applied to all outgoing HTTP request.
}

// NewClient creates a client.
func NewClient() *Client {
	return &Client{
		HTTPClient:    http.DefaultClient,
		DefaultHeader: make(http.Header),
	}
}

// DoWithContext sends a request and returns a HTTP response whose body is consumed and
// closed. The context controls the lifetime of the outgoing request and its response.
//
// If respV is not nil, the value it points to is JSON decoded when server
// returns a successful response.
//
// If errV is not nil, the value it points to is JSON decoded when server
// returns an unsuccessfully response. If the response text is not a JSON
// string, a more generic ErrorResponse error is returned.
func (c *Client) DoWithContext(ctx context.Context, r *Request, respV interface{}, errV interface{}) (*http.Response, error) {
	req, err := c.makeRequest(ctx, r)
	if err != nil {
		return nil, err
	}

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	defer func() error {
		if err := resp.Body.Close(); err != nil {
			return err
		}
		return nil
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, fmt.Errorf("Error reading response: %v", err)
		}

		if len(raw) > 0 && errV != nil {
			if json.Unmarshal(raw, errV) == nil {
				return resp, nil
			}
		}

		return resp, &ErrorResponse{resp.StatusCode, string(raw)}
	}

	if respV != nil {
		switch respV.(type) {
		case io.Writer:
			_, err = io.Copy(respV.(io.Writer), resp.Body)
		default:
			// Determine the response type and the decoder that should be used.
			// If buffer is identified as json, use JSON decoder, otherwise
			// assume the buffer contains yaml bytes
			body, isJSON := IsJSONStream(resp.Body, bufferSize)
			if isJSON {
				err = json.NewDecoder(body).Decode(respV)
			} else {
				err = yaml.NewDecoder(body).Decode(respV)
			}
			// For 204 No Content we should not throw an error
			// if there is an empty response body
			if err == io.EOF && resp.StatusCode == http.StatusNoContent {
				err = nil
			} else if err == io.EOF {
				err = ErrEmptyResponseBody
			}
		}
	}

	return resp, err
}

// Do wraps DoWithContext using the background context.
func (c *Client) Do(r *Request, respV interface{}, errV interface{}) (*http.Response, error) {
	return c.DoWithContext(context.Background(), r, respV, errV)
}

func (c *Client) makeRequest(ctx context.Context, r *Request) (*http.Request, error) {
	req, err := r.Build()
	if err != nil {
		return nil, err
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	c.applyDefaultHeader(req)

	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) applyDefaultHeader(req *http.Request) {
	for k, vs := range c.DefaultHeader {
		if req.Header.Get(k) != "" {
			continue
		}
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
}
