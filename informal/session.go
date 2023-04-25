package informal

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kenchan0130/go-jamf-pro/jamf"
)

type SessionService service

func (s *SessionService) Create(ctx context.Context, loginURL string) (*jamf.Response, error) {
	u, err := url.Parse(loginURL)
	if err != nil {
		return nil, fmt.Errorf("url.Parse(): %v", err)
	}

	values := url.Values{
		"username": {s.username},
		"password": {s.password},
	}

	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusFound},
		ContentType:      "application/x-www-form-urlencoded",
		Uri: jamf.Uri{
			Entity: u.Path,
			Params: u.Query(),
		},
		Body: bytes.NewBufferString(values.Encode()),
	})

	if err != nil {
		return nil, fmt.Errorf("client.Post(): %v", err)
	}

	return resp, nil
}
