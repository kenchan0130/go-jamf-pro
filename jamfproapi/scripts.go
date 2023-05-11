package jamfproapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/google/go-querystring/query"
	"github.com/kenchan0130/go-jamf-pro/jamf"
	"github.com/kenchan0130/go-jamf-pro/utils"
)

type ScriptsService service

type ScriptPriority string

const (
	ScriptPriorityBefore   ScriptPriority = "BEFORE"
	ScriptPriorityAfter    ScriptPriority = "AFTER"
	ScriptPriorityAtReboot ScriptPriority = "AT_REBOOT"
)

type Script struct {
	ID             *string         `json:"id,omitempty"`
	Name           *string         `json:"name,omitempty"`
	Info           *string         `json:"info,omitempty"`
	Notes          *string         `json:"notes,omitempty"`
	Priority       *ScriptPriority `json:"priority,omitempty"`
	CategoryID     *string         `json:"categoryId,omitempty"`
	CategoryName   *string         `json:"categoryName,omitempty"`
	Parameter4     *string         `json:"parameter4,omitempty"`
	Parameter5     *string         `json:"parameter5,omitempty"`
	Parameter6     *string         `json:"parameter6,omitempty"`
	Parameter7     *string         `json:"parameter7,omitempty"`
	Parameter8     *string         `json:"parameter8,omitempty"`
	Parameter9     *string         `json:"parameter9,omitempty"`
	Parameter10    *string         `json:"parameter10,omitempty"`
	Parameter11    *string         `json:"parameter11,omitempty"`
	OSRequirements *string         `json:"osRequirements,omitempty"`
	ScriptContents *string         `json:"scriptContents,omitempty"`
}

type ListScript struct {
	TotalCount *int      `json:"totalCount,omitempty"`
	Scripts    *[]Script `json:"results,omitempty"`
}

const scriptsPath = "/v1/scripts"

func (s *ScriptsService) Create(ctx context.Context, script *Script) (*string, *jamf.Response, error) {
	if script.Name == nil {
		return nil, nil, errors.New("ScriptsService.Create(): cannot create script with nil Name")
	}

	body, err := json.Marshal(script)
	if err != nil {
		return nil, nil, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: jamf.Uri{
			Entity: scriptsPath,
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Post(): %v", err)
	}
	defer utils.HandleCloseFunc(resp.Body, s.client.RetryableClient.Logger)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		ID   string `json:"id"`
		Href string `json:"href"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.ID, resp, nil
}

func (s *ScriptsService) Delete(ctx context.Context, scriptID string) (*jamf.Response, error) {
	resp, _, err := s.client.Delete(ctx, jamf.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: jamf.Uri{
			Entity: path.Join(scriptsPath, scriptID),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("client.Delete(): %v", err)
	}

	return resp, nil
}

func (s *ScriptsService) Get(ctx context.Context, scriptID string) (*Script, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: path.Join(scriptsPath, scriptID),
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

	var script Script
	if err := json.Unmarshal(respBody, &script); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &script, resp, nil
}

func (s *ScriptsService) List(ctx context.Context, options ListOptions) (*ListScript, *jamf.Response, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, nil, fmt.Errorf("query.Values(): %v", err)
	}

	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: scriptsPath,
			Params: params,
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

	var listScript ListScript
	if err := json.Unmarshal(respBody, &listScript); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &listScript, resp, nil
}

func (s *ScriptsService) Update(ctx context.Context, script *Script) (*Script, *jamf.Response, error) {
	if script.ID == nil {
		return nil, nil, errors.New("ScriptsService.Update(): cannot update script with nil ID")
	}
	if script.Name == nil {
		return nil, nil, errors.New("ScriptsService.Update(): cannot update script with nil Name")
	}

	body, err := json.Marshal(script)
	if err != nil {
		return nil, nil, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, _, err := s.client.Put(ctx, jamf.PutHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: path.Join(scriptsPath, *script.ID),
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Put(): %v", err)
	}
	defer utils.HandleCloseFunc(resp.Body, s.client.RetryableClient.Logger)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newScript Script
	if err := json.Unmarshal(respBody, &newScript); err != nil {
		return nil, nil, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newScript, nil, nil
}
