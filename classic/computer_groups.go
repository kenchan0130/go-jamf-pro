package classic

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/kenchan0130/go-jamf-pro/jamf"
	"github.com/kenchan0130/go-jamf-pro/utils"
)

type ComputerGroupsService service

type ComputerGroup struct {
	ID        *int                    `xml:"id,omitempty"`
	Name      *string                 `xml:"name,omitempty"`
	IsSmart   *bool                   `xml:"is_smart,omitempty"`
	Site      *Site                   `xml:"site,omitempty"`
	Criteria  *ComputerGroupCriteria  `xml:"criteria,omitempty"`
	Computers *ComputerGroupComputers `xml:"computers,omitempty"`
}

type ComputerGroupComputer struct {
	ID            *int    `xml:"id,omitempty"`
	Name          *string `xml:"name,omitempty"`
	MacAddress    *string `xml:"mac_address,omitempty"`
	AltMacAddress *string `xml:"alt_mac_address,omitempty"`
	SerialNumber  *string `xml:"serial_number,omitempty"`
}

type ComputerGroupComputers struct {
	Size      *int                     `xml:"size,omitempty"`
	Computers *[]ComputerGroupComputer `xml:"computer,omitempty"`
}

type ComputerGroupCriteria struct {
	Size      *int                              `xml:"size,omitempty"`
	Criterion *[]ComputerGroupCriteriaCriterion `xml:"criterion,omitempty"`
}

type ComputerGroupCriteriaCriterion struct {
	Name         *string                                   `xml:"name,omitempty"`
	Priority     *int                                      `xml:"priority,omitempty"`
	AndOr        *ComputerGroupCriteriaCriterionAndOr      `xml:"and_or,omitempty"`
	SearchType   *ComputerGroupCriteriaCriterionSearchType `xml:"search_type,omitempty"`
	Value        *string                                   `xml:"value,omitempty"`
	OpeningParen *bool                                     `xml:"opening_paren,omitempty"`
	ClosingParen *bool                                     `xml:"closing_paren,omitempty"`
}

type ComputerGroupCriteriaCriterionAndOr string

const (
	ComputerGroupCriteriaCriterionAndOrAnd ComputerGroupCriteriaCriterionAndOr = "and"
	ComputerGroupCriteriaCriterionAndOrOr  ComputerGroupCriteriaCriterionAndOr = "or"
)

type ComputerGroupCriteriaCriterionSearchType string

const (
	ComputerGroupCriteriaCriterionSearchTypeIs                 ComputerGroupCriteriaCriterionSearchType = "is"
	ComputerGroupCriteriaCriterionSearchTypeIsNot              ComputerGroupCriteriaCriterionSearchType = "is not"
	ComputerGroupCriteriaCriterionSearchTypeHas                ComputerGroupCriteriaCriterionSearchType = "has"
	ComputerGroupCriteriaCriterionSearchTypeDoesNotHave        ComputerGroupCriteriaCriterionSearchType = "does not have"
	ComputerGroupCriteriaCriterionSearchTypeBefore             ComputerGroupCriteriaCriterionSearchType = "before (yyyy-mm-dd)"
	ComputerGroupCriteriaCriterionSearchTypeAfter              ComputerGroupCriteriaCriterionSearchType = "after (yyyy-mm-dd)"
	ComputerGroupCriteriaCriterionSearchTypeMoreThanXDaysAgo   ComputerGroupCriteriaCriterionSearchType = "more than x days ago"
	ComputerGroupCriteriaCriterionSearchTypeLessThanXDaysAgo   ComputerGroupCriteriaCriterionSearchType = "less than x days ago"
	ComputerGroupCriteriaCriterionSearchTypeLike               ComputerGroupCriteriaCriterionSearchType = "like"
	ComputerGroupCriteriaCriterionSearchTypeNotLike            ComputerGroupCriteriaCriterionSearchType = "note like"
	ComputerGroupCriteriaCriterionSearchTypeGreaterThan        ComputerGroupCriteriaCriterionSearchType = "greater than"
	ComputerGroupCriteriaCriterionSearchTypeLessThan           ComputerGroupCriteriaCriterionSearchType = "less than"
	ComputerGroupCriteriaCriterionSearchTypeGreaterThanOrEqual ComputerGroupCriteriaCriterionSearchType = "greater than or equal"
	ComputerGroupCriteriaCriterionSearchTypeLessThanOrEqual    ComputerGroupCriteriaCriterionSearchType = "less than or equal"
	ComputerGroupCriteriaCriterionSearchTypeMatchesRegex       ComputerGroupCriteriaCriterionSearchType = "matches regex"
	ComputerGroupCriteriaCriterionSearchTypeDoesNotMatchRegex  ComputerGroupCriteriaCriterionSearchType = "does not match regex"
)

type ListComputerGroups struct {
	Size           *int                 `xml:"size,omitempty"`
	ComputerGroups *[]ListComputerGroup `xml:"computer_group,omitempty"`
}

type ListComputerGroup struct {
	ID      *int    `xml:"id,omitempty"`
	Name    *string `xml:"name,omitempty"`
	IsSmart *bool   `xml:"is_smart,omitempty"`
}

const computerGroupsPath = "/computergroups"

func (s *ComputerGroupsService) Create(ctx context.Context, computerGroup *ComputerGroup) (*int, *jamf.Response, error) {
	if computerGroup == nil {
		return nil, nil, errors.New("ComputerGroupsService.Create(): cannot create nil computer group")
	}
	if computerGroup.Name == nil {
		return nil, nil, errors.New("ComputerGroupsService.Create(): cannot create computer group with nil Name")
	}
	if computerGroup.IsSmart == nil {
		return nil, nil, errors.New("ComputerGroupsService.Create(): cannot create computer group with nil IsSmart")
	}

	reqBody := &struct {
		*ComputerGroup
		XMLName xml.Name `xml:"computer_group"`
	}{
		ComputerGroup: computerGroup,
	}
	body, err := xml.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("xml.Marshal(): %v", err)
	}

	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: jamf.Uri{
			// When the ID is 0, the ID will be generated automatically at the server side.
			Entity: path.Join(computerGroupsPath, "id", "0"),
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
		XMLName xml.Name `xml:"computer_group"`
		ID      *int     `xml:"id"`
	}
	if err := xml.Unmarshal(respBody, &data); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	// The create API only returns the ID of the new ComputerGroup.
	return data.ID, resp, nil
}

func (s *ComputerGroupsService) Delete(ctx context.Context, computerGroupID int) (*jamf.Response, error) {
	resp, _, err := s.client.Delete(ctx, jamf.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(computerGroupsPath, "id", fmt.Sprint(computerGroupID)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("client.Delete(): %v", err)
	}

	return resp, nil
}

func (s *ComputerGroupsService) Get(ctx context.Context, computerGroupID int) (*ComputerGroup, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(computerGroupsPath, "id", fmt.Sprint(computerGroupID)),
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

	var computerGroup ComputerGroup
	if err := xml.Unmarshal(respBody, &computerGroup); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	return &computerGroup, resp, nil
}

func (s *ComputerGroupsService) List(ctx context.Context) (*ListComputerGroups, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: computerGroupsPath,
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

	var listComputerGroups ListComputerGroups
	if err := xml.Unmarshal(respBody, &listComputerGroups); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	return &listComputerGroups, resp, nil
}

func (s *ComputerGroupsService) Update(ctx context.Context, computerGroup *ComputerGroup) (*jamf.Response, error) {
	if computerGroup == nil {
		return nil, errors.New("ComputerGroupsService.Update(): cannot update nil computer group")
	}
	if computerGroup.Name == nil {
		return nil, errors.New("ComputerGroupsService.Update(): cannot update computer group with nil Name")
	}
	if computerGroup.IsSmart == nil {
		return nil, errors.New("ComputerGroupsService.Create(): cannot update computer group with nil IsSmart")
	}

	reqBody := &struct {
		*ComputerGroup
		XMLName xml.Name `xml:"computer_group"`
	}{
		ComputerGroup: computerGroup,
	}
	body, err := xml.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("xml.Marshal(): %v", err)
	}

	resp, _, err := s.client.Put(ctx, jamf.PutHttpRequestInput{
		ValidStatusCodes: []int{http.StatusCreated},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(computerGroupsPath, "id", fmt.Sprint(*computerGroup.ID)),
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, fmt.Errorf("client.Put(): %v", err)
	}

	return resp, nil
}
