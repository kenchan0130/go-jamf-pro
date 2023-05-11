package jamfproapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/kenchan0130/go-jamf-pro/jamf"
	"github.com/kenchan0130/go-jamf-pro/utils"
)

type APIAuthenticationService service

type Authorizer struct {
	BaseURL  *url.URL
	Username string
	Password string
}

type AuthToken struct {
	Token   *string `json:"token,omitempty"`
	Expires *string `json:"expires,omitempty"`
}

const apiAuthenticationPath = "/v1/auth"

func (s *APIAuthenticationService) Token(ctx context.Context, username string, password string) (*AuthToken, *jamf.Response, error) {
	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		RequestMiddlewareFunc: func(r *http.Request) {
			r.SetBasicAuth(username, password)
		},
		Uri: jamf.Uri{
			Entity: path.Join(apiAuthenticationPath, "token"),
		},
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Post(): %v", err)
	}
	defer utils.HandleCloseFunc(resp.Body, s.client.RetryableClient.Logger)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var authToken AuthToken
	if err := json.Unmarshal(respBody, &authToken); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &authToken, resp, nil
}
