package jamfproapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/kenchan0130/go-jamf-pro/jamf"
	"github.com/kenchan0130/go-jamf-pro/utils"
)

type SSOFailoverService service

type SSOFailover struct {
	FailoverURL    *string `json:"failoverUrl,omitempty"`
	GenerationTime *int64  `json:"generationTime,omitempty"`
}

const ssoFailoverPath = "/v1/sso/failover"

func (s *SSOFailoverService) Get(ctx context.Context) (*SSOFailover, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: ssoFailoverPath,
		},
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Get(): %v", err)
	}
	defer utils.HandleCloseFunc(resp.Body, s.client.RetryableClient.Logger)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var ssoFailover SSOFailover
	if err := json.Unmarshal(respBody, &ssoFailover); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &ssoFailover, resp, nil
}

func (s *SSOFailoverService) Generate(ctx context.Context) (*SSOFailover, *jamf.Response, error) {
	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: path.Join(ssoFailoverPath, "generate"),
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

	var ssoFailover SSOFailover
	if err := json.Unmarshal(respBody, &ssoFailover); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &ssoFailover, resp, nil
}
