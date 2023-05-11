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

type CategoriesService service

type Category struct {
	ID       *string `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Priority *int32  `json:"priority,omitempty"`
}

type ListCategory struct {
	TotalCount *int        `json:"totalCount,omitempty"`
	Categories *[]Category `json:"results,omitempty"`
}

const categoriesPath = "/v1/categories"

func (s *CategoriesService) Create(ctx context.Context, category *Category) (*string, *jamf.Response, error) {
	if category.Name == nil {
		return nil, nil, errors.New("CategoriesService.Create(): cannot create category with nil Name")
	}
	if category.Priority == nil {
		return nil, nil, errors.New("CategoriesService.Create(): cannot create category with nil Priority")
	}

	body, err := json.Marshal(category)
	if err != nil {
		return nil, nil, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: jamf.Uri{
			Entity: categoriesPath,
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Post(): %v", err)
	}
	defer utils.HandleCloseFunc(resp.Body, s.client.RetryableClient.Logger)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("io.ReadAll(): %v", err)
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

func (s *CategoriesService) Delete(ctx context.Context, categoryID string) (*jamf.Response, error) {
	resp, _, err := s.client.Delete(ctx, jamf.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: jamf.Uri{
			Entity: path.Join(categoriesPath, categoryID),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("client.Delete(): %v", err)
	}

	return resp, nil
}

func (s *CategoriesService) DeleteMultiple(ctx context.Context, categoryIDs []string) (*jamf.Response, error) {
	var data struct {
		CategoryIDs []string `json:"ids"`
	}
	data.CategoryIDs = categoryIDs

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: jamf.Uri{
			Entity: path.Join(categoriesPath, "delete-multiple"),
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, fmt.Errorf("client.Post(): %v", err)
	}

	return resp, nil
}

func (s *CategoriesService) Get(ctx context.Context, categoryID string) (*Category, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: path.Join(categoriesPath, categoryID),
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

	var category Category
	if err := json.Unmarshal(respBody, &category); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &category, resp, nil
}

func (s *CategoriesService) List(ctx context.Context, options ListOptions) (*ListCategory, *jamf.Response, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, nil, fmt.Errorf("query.Values(): %v", err)
	}

	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: categoriesPath,
			Params: params,
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

	var listCategory ListCategory
	if err := json.Unmarshal(respBody, &listCategory); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &listCategory, resp, nil
}

func (s *CategoriesService) Update(ctx context.Context, category *Category) (*Category, *jamf.Response, error) {
	if category.ID == nil {
		return nil, nil, errors.New("CategoriesService.Update(): cannot update category with nil ID")
	}
	if category.Name == nil {
		return nil, nil, errors.New("CategoriesService.Update(): cannot update category with nil Name")
	}
	if category.Priority == nil {
		return nil, nil, errors.New("CategoriesService.Update(): cannot update category with nil Priority")
	}

	body, err := json.Marshal(category)
	if err != nil {
		return nil, nil, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, _, err := s.client.Put(ctx, jamf.PutHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: path.Join(categoriesPath, *category.ID),
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Put(): %v", err)
	}
	defer utils.HandleCloseFunc(resp.Body, s.client.RetryableClient.Logger)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newCategory Category
	if err := json.Unmarshal(respBody, &newCategory); err != nil {
		return nil, resp, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newCategory, resp, nil
}
