// Referred https://github.com/manicminer/hamilton
// Apache License 2.0 https://github.com/manicminer/hamilton/blob/main/LICENSE

package jamf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

// ConsistencyFailureFunc is a function that determines whether an HTTP request has failed due to eventual consistency and should be retried.
type ConsistencyFailureFunc func(*http.Response) bool

// RequestMiddlewareFunc is a function that modifies an HTTP request before it is sent.
type RequestMiddlewareFunc func(r *http.Request)

// ValidStatusFunc is a function that tests whether an HTTP response is considered valid for the particular request.
type ValidStatusFunc func(*http.Response) bool

// HttpRequestInput is any type that can validate the response to an HTTP request.
type HttpRequestInput interface {
	GetRequestMiddlewareFunc() RequestMiddlewareFunc
	GetConsistencyFailureFunc() ConsistencyFailureFunc
	GetContentType() string
	GetValidStatusCodes() []int
	GetValidStatusFunc() ValidStatusFunc
}

type Uri struct {
	Entity string
	Params url.Values
}

// RetryableErrorHandler ensures that the response is returned after exhausting retries for a request
func RetryableErrorHandler(resp *http.Response, err error, numTries int) (*http.Response, error) {
	return resp, nil
}

type BaseClient struct {
	// BaseURL is the base endpoint for Jamf Pro, e.g. https://yourdomain.jamfcloud.com
	BaseURL            *url.URL
	AuthorizationToken *string
	DisableRetries     bool

	// HttpClient is the underlying http.Client, which by default uses a retryable client
	HttpClient      *http.Client
	RetryableClient *retryablehttp.Client

	// Logger is the logger, which by default uses the standard log package
	Logger *log.Logger
}

type Response struct {
	*http.Response
}

func NewBaseClient(baseURL *url.URL) *BaseClient {
	r := retryablehttp.NewClient()
	r.ErrorHandler = RetryableErrorHandler

	c := &BaseClient{
		BaseURL:         baseURL,
		HttpClient:      r.StandardClient(),
		RetryableClient: r,
		DisableRetries:  false,
		Logger:          log.Default(),
	}

	return c
}

func (c *BaseClient) buildUri(uri Uri) string {
	newUrl := c.BaseURL.JoinPath(uri.Entity)
	if uri.Params != nil {
		newUrl.RawQuery = uri.Params.Encode()
	}
	return newUrl.String()
}

func (c *BaseClient) performRequest(req *http.Request, input HttpRequestInput) (*http.Response, int, error) {
	var status int

	req.Header.Add("Content-Type", input.GetContentType())
	if c.AuthorizationToken != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *c.AuthorizationToken))
	}

	f := input.GetRequestMiddlewareFunc()
	if f != nil {
		f(req)
	}

	c.RetryableClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp != nil && !c.DisableRetries {
			if resp.StatusCode == http.StatusFailedDependency {
				return true, nil
			}

			f := input.GetConsistencyFailureFunc()
			if f != nil && f(resp) {
				return true, nil
			}
		}

		//nolint:golint,wrapcheck // retryablehttp.DefaultRetryPolicy returns a non-wrapped error
		return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, status, fmt.Errorf("http.Client#Do(): %v", err)
	}

	if resp == nil {
		return resp, status, fmt.Errorf("nil response received")
	}

	status = resp.StatusCode
	if !containsStatusCode(input.GetValidStatusCodes(), status) {
		f := input.GetValidStatusFunc()
		if f != nil && f(resp) {
			return resp, status, nil
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				c.Logger.Printf("Error closing response body: %v", err)
			}
		}()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, status, fmt.Errorf("unexpected status %d, could not read response body", status)
		}
		if len(respBody) == 0 {
			return nil, status, fmt.Errorf("unexpected status %d received with no body", status)
		}
		errText := fmt.Sprintf("response: %s", respBody)
		return nil, status, fmt.Errorf("unexpected status %d with %s", status, errText)
	}

	return resp, status, nil
}

func containsStatusCode(expected []int, actual int) bool {
	for _, v := range expected {
		if actual == v {
			return true
		}
	}

	return false
}

type DeleteHttpRequestInput struct {
	ConsistencyFailureFunc ConsistencyFailureFunc
	RequestMiddlewareFunc  RequestMiddlewareFunc
	Uri                    Uri
	ValidStatusCodes       []int
	ValidStatusFunc        ValidStatusFunc
}

// GetConsistencyFailureFunc returns a function used to evaluate whether a failed request is due to eventual consistency and should be retried.
func (i DeleteHttpRequestInput) GetConsistencyFailureFunc() ConsistencyFailureFunc {
	return i.ConsistencyFailureFunc
}

// GetContentType returns the content type for the request, currently only application/json is supported
func (i DeleteHttpRequestInput) GetContentType() string {
	return "application/json; charset=utf-8"
}

func (i DeleteHttpRequestInput) GetRequestMiddlewareFunc() RequestMiddlewareFunc {
	return i.RequestMiddlewareFunc
}

// GetValidStatusCodes returns a []int of status codes considered valid for a DELETE request.
func (i DeleteHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

// GetValidStatusFunc returns a function used to evaluate whether the response to a DELETE request is considered valid.
func (i DeleteHttpRequestInput) GetValidStatusFunc() ValidStatusFunc {
	return i.ValidStatusFunc
}

// Delete performs a DELETE request.
func (c *BaseClient) Delete(ctx context.Context, input DeleteHttpRequestInput) (*Response, int, error) {
	var status int

	u := c.buildUri(input.Uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, http.NoBody)
	if err != nil {
		return nil, status, fmt.Errorf("http.NewRequestWithContext(): %v", err)
	}

	resp, status, err := c.performRequest(req, input)
	if err != nil {
		return nil, status, err
	}

	return &Response{resp}, status, nil
}

// GetHttpRequestInput configures a GET request.
type GetHttpRequestInput struct {
	ConsistencyFailureFunc ConsistencyFailureFunc
	ContentType            string
	RequestMiddlewareFunc  RequestMiddlewareFunc
	Uri                    Uri
	ValidStatusCodes       []int
	ValidStatusFunc        ValidStatusFunc
}

// GetConsistencyFailureFunc returns a function used to evaluate whether a failed request is due to eventual consistency and should be retried.
func (i GetHttpRequestInput) GetConsistencyFailureFunc() ConsistencyFailureFunc {
	return i.ConsistencyFailureFunc
}

// GetContentType returns the content type for the request, defaults to application/json
func (i GetHttpRequestInput) GetContentType() string {
	if i.ContentType != "" {
		return i.ContentType
	}
	return "application/json; charset=utf-8"
}

func (i GetHttpRequestInput) GetRequestMiddlewareFunc() RequestMiddlewareFunc {
	return i.RequestMiddlewareFunc
}

// GetValidStatusCodes returns a []int of status codes considered valid for a GET request.
func (i GetHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

// GetValidStatusFunc returns a function used to evaluate whether the response to a GET request is considered valid.
func (i GetHttpRequestInput) GetValidStatusFunc() ValidStatusFunc {
	return i.ValidStatusFunc
}

// Get performs a GET request.
func (c *BaseClient) Get(ctx context.Context, input GetHttpRequestInput) (*Response, int, error) {
	var status int

	u := c.buildUri(input.Uri)

	// Build a new request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, status, fmt.Errorf("http.NewRequestWithContext(): %v", err)
	}

	// Perform the request
	resp, status, err := c.performRequest(req, input)
	if err != nil {
		return nil, status, err
	}

	return &Response{resp}, status, nil
}

// PatchHttpRequestInput configures a PATCH request.
type PatchHttpRequestInput struct {
	Body                   *bytes.Buffer
	ConsistencyFailureFunc ConsistencyFailureFunc
	ContentType            string
	RequestMiddlewareFunc  RequestMiddlewareFunc
	Uri                    Uri
	ValidStatusCodes       []int
	ValidStatusFunc        ValidStatusFunc
}

// GetConsistencyFailureFunc returns a function used to evaluate whether a failed request is due to eventual consistency and should be retried.
func (i PatchHttpRequestInput) GetConsistencyFailureFunc() ConsistencyFailureFunc {
	return i.ConsistencyFailureFunc
}

// GetContentType returns the content type for the request, defaults to application/json
func (i PatchHttpRequestInput) GetContentType() string {
	if i.ContentType != "" {
		return i.ContentType
	}
	return "application/json; charset=utf-8"
}

func (i PatchHttpRequestInput) GetRequestMiddlewareFunc() RequestMiddlewareFunc {
	return i.RequestMiddlewareFunc
}

// GetValidStatusCodes returns a []int of status codes considered valid for a PATCH request.
func (i PatchHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

// GetValidStatusFunc returns a function used to evaluate whether the response to a PATCH request is considered valid.
func (i PatchHttpRequestInput) GetValidStatusFunc() ValidStatusFunc {
	return i.ValidStatusFunc
}

// Patch performs a PATCH request.
func (c *BaseClient) Patch(ctx context.Context, input PatchHttpRequestInput) (*Response, int, error) {
	var status int

	u := c.buildUri(input.Uri)

	var inputBody io.Reader
	if input.Body == nil {
		inputBody = http.NoBody
	} else {
		inputBody = input.Body
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u, inputBody)
	if err != nil {
		return nil, status, fmt.Errorf("http.NewRequestWithContext(): %v", err)
	}

	resp, status, err := c.performRequest(req, input)
	if err != nil {
		return nil, status, err
	}

	return &Response{resp}, status, nil
}

// PostHttpRequestInput configures a POST request.
type PostHttpRequestInput struct {
	Body                   *bytes.Buffer
	ContentType            string
	ConsistencyFailureFunc ConsistencyFailureFunc
	RequestMiddlewareFunc  RequestMiddlewareFunc
	Uri                    Uri
	ValidStatusCodes       []int
	ValidStatusFunc        ValidStatusFunc
}

// GetConsistencyFailureFunc returns a function used to evaluate whether a failed request is due to eventual consistency and should be retried.
func (i PostHttpRequestInput) GetConsistencyFailureFunc() ConsistencyFailureFunc {
	return i.ConsistencyFailureFunc
}

// GetContentType returns the content type for the request, defaults to application/json
func (i PostHttpRequestInput) GetContentType() string {
	if i.ContentType != "" {
		return i.ContentType
	}
	return "application/json; charset=utf-8"
}

func (i PostHttpRequestInput) GetRequestMiddlewareFunc() RequestMiddlewareFunc {
	return i.RequestMiddlewareFunc
}

// GetValidStatusCodes returns a []int of status codes considered valid for a POST request.
func (i PostHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

// GetValidStatusFunc returns a function used to evaluate whether the response to a POST request is considered valid.
func (i PostHttpRequestInput) GetValidStatusFunc() ValidStatusFunc {
	return i.ValidStatusFunc
}

// Post performs a POST request.
func (c *BaseClient) Post(ctx context.Context, input PostHttpRequestInput) (*Response, int, error) {
	var status int

	u := c.buildUri(input.Uri)

	var inputBody io.Reader
	if input.Body == nil {
		inputBody = http.NoBody
	} else {
		inputBody = input.Body
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, inputBody)
	if err != nil {
		return nil, status, fmt.Errorf("http.NewRequestWithContext(): %v", err)
	}

	resp, status, err := c.performRequest(req, input)
	if err != nil {
		return nil, status, err
	}

	return &Response{resp}, status, nil
}

// PutHttpRequestInput configures a PUT request.
type PutHttpRequestInput struct {
	Body                   *bytes.Buffer
	ConsistencyFailureFunc ConsistencyFailureFunc
	ContentType            string
	RequestMiddlewareFunc  RequestMiddlewareFunc
	Uri                    Uri
	ValidStatusCodes       []int
	ValidStatusFunc        ValidStatusFunc
}

// GetConsistencyFailureFunc returns a function used to evaluate whether a failed request is due to eventual consistency and should be retried.
func (i PutHttpRequestInput) GetConsistencyFailureFunc() ConsistencyFailureFunc {
	return i.ConsistencyFailureFunc
}

// GetContentType returns the content type for the request, defaults to application/json
func (i PutHttpRequestInput) GetContentType() string {
	if i.ContentType != "" {
		return i.ContentType
	}
	return "application/json; charset=utf-8"
}

func (i PutHttpRequestInput) GetRequestMiddlewareFunc() RequestMiddlewareFunc {
	return i.RequestMiddlewareFunc
}

// GetValidStatusCodes returns a []int of status codes considered valid for a PUT request.
func (i PutHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

// GetValidStatusFunc returns a function used to evaluate whether the response to a PUT request is considered valid.
func (i PutHttpRequestInput) GetValidStatusFunc() ValidStatusFunc {
	return i.ValidStatusFunc
}

// Put performs a PUT request.
func (c *BaseClient) Put(ctx context.Context, input PutHttpRequestInput) (*Response, int, error) {
	var status int

	u := c.buildUri(input.Uri)

	var inputBody io.Reader
	if input.Body == nil {
		inputBody = http.NoBody
	} else {
		inputBody = input.Body
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, inputBody)
	if err != nil {
		return nil, status, fmt.Errorf("http.NewRequestWithContext(): %v", err)
	}

	resp, status, err := c.performRequest(req, input)
	if err != nil {
		return nil, status, err
	}

	return &Response{resp}, status, nil
}
