package jamfproapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/google/go-querystring/query"
	"github.com/kenchan0130/go-jamf-pro/jamf"
	"github.com/kenchan0130/go-jamf-pro/utils"
)

type IconService service

type Icon struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

type IconContentOptions struct {
	Res   *IconContentOptionsRes `url:"res,omitempty"`
	Scale *string                `url:"scale,omitempty"`
}

type IconContentOptionsRes string

const (
	IconContentOptionsResOriginal IconContentOptionsRes = "original"
	IconContentOptionsRes300      IconContentOptionsRes = "300"
	IconContentOptionsRes512      IconContentOptionsRes = "512"
)

const iconPath = "/v1/icon"

func (s *IconService) Download(ctx context.Context, iconID int, options IconContentOptions) (*jamf.Response, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, fmt.Errorf("query.Values(): %v", err)
	}

	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: path.Join(iconPath, "download", fmt.Sprint(iconID)),
			Params: params,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("client.Get(): %v", err)
	}

	return resp, nil
}

func (s *IconService) Get(ctx context.Context, iconID int) (*Icon, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: path.Join(iconPath, fmt.Sprint(iconID)),
		},
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Get(): %v", err)
	}
	defer utils.HandleCloseFunc(resp.Body, s.client.RetryableClient.Logger)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var icon Icon
	if err := json.Unmarshal(respBody, &icon); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &icon, resp, nil
}

func (s *IconService) Upload(ctx context.Context, iconName string, src io.Reader) (*Icon, *jamf.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", iconName)
	if err != nil {
		return nil, nil, fmt.Errorf("writer.CreateFormFile(): %v", err)
	}

	_, err = io.Copy(part, src)
	if err != nil {
		return nil, nil, fmt.Errorf("io.Copy(): %v", err)
	}

	contentType := writer.FormDataContentType()

	err = writer.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("writer.Close(): %v", err)
	}

	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: jamf.Uri{
			Entity: iconPath,
		},
		ContentType: contentType,
		Body:        body,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Post(): %v", err)
	}
	defer utils.HandleCloseFunc(resp.Body, s.client.RetryableClient.Logger)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var icon Icon
	if err := json.Unmarshal(respBody, &icon); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &icon, resp, nil
}
