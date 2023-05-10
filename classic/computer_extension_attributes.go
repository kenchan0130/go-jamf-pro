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
)

type ComputerExtensionAttributesService service

type ComputerExtensionAttribute struct {
	ID               *int                                        `xml:"id,omitempty"`
	Name             *string                                     `xml:"name,omitempty"`
	Enabled          *bool                                       `xml:"enabled,omitempty"`
	Description      *string                                     `xml:"description,omitempty"`
	DataType         *ComputerExtensionAttributeDataType         `xml:"data_type,omitempty"`
	InputType        *ComputerExtensionAttributeInputType        `xml:"input_type,omitempty"`
	InventoryDisplay *ComputerExtensionAttributeInventoryDisplay `xml:"inventory_display,omitempty"`
}

type ComputerExtensionAttributeDataType string

const (
	ComputerExtensionAttributeDataTypeString  ComputerExtensionAttributeDataType = "String"
	ComputerExtensionAttributeDataTypeInteger ComputerExtensionAttributeDataType = "Integer"
	ComputerExtensionAttributeDataTypeDate    ComputerExtensionAttributeDataType = "Date"
)

type ComputerExtensionAttributeInputType struct {
	Type *ComputerExtensionAttributeInputTypeType `xml:"type,omitempty"`
	// The 'platform' field is not documented, but it actually exists in the API response.
	Platform *ComputerExtensionAttributeInputTypePlatform `xml:"platform,omitempty"`
	// The 'popup_choices' field is not documented, but it actually exists in the API response.
	PopupChoices *[]string `xml:"popup_choices>choice,omitempty"`
	// The 'script' field is not documented, but it actually exists in the API response.
	Script *string `xml:"script,omitempty"`
}

type ComputerExtensionAttributeInputTypePlatform string

const (
	ComputerExtensionAttributeInputTypePlatformMac ComputerExtensionAttributeInputTypePlatform = "Mac"
)

type ComputerExtensionAttributeInputTypeType string

const (
	ComputerExtensionAttributeInputTypeTypeScript    ComputerExtensionAttributeInputTypeType = "script"
	ComputerExtensionAttributeInputTypeTypeTextField ComputerExtensionAttributeInputTypeType = "Text Field"
	ComputerExtensionAttributeInputTypeTypePopupMenu ComputerExtensionAttributeInputTypeType = "Pop-up Menu"
	// The 'LDAP Mapping' value is documented, but it does not actually exist in the API response.
	//ComputerExtensionAttributeInputTypeTypeLDAPMapping ComputerExtensionAttributeInputTypeType = "LDAP Mapping"
)

type ComputerExtensionAttributeInventoryDisplay string

const (
	ComputerExtensionAttributeInventoryDisplayGeneral             ComputerExtensionAttributeInventoryDisplay = "General"
	ComputerExtensionAttributeInventoryDisplayHardware            ComputerExtensionAttributeInventoryDisplay = "Hardware"
	ComputerExtensionAttributeInventoryDisplayOperatingSystem     ComputerExtensionAttributeInventoryDisplay = "Operating System"
	ComputerExtensionAttributeInventoryDisplayUserAndLocation     ComputerExtensionAttributeInventoryDisplay = "User and Location"
	ComputerExtensionAttributeInventoryDisplayPurchasing          ComputerExtensionAttributeInventoryDisplay = "Purchasing"
	ComputerExtensionAttributeInventoryDisplayExtensionAttributes ComputerExtensionAttributeInventoryDisplay = "Extension Attributes"
)

type ListComputerExtensionAttributes struct {
	Size                        *int                              `xml:"size,omitempty"`
	ComputerExtensionAttributes *[]ListComputerExtensionAttribute `xml:"computer_extension_attribute,omitempty"`
}

type ListComputerExtensionAttribute struct {
	ID      *int    `xml:"id,omitempty"`
	Name    *string `xml:"name,omitempty"`
	Enabled *bool   `xml:"enabled,omitempty"`
}

const computerExtensionAttributesPath = "/computerextensionattributes"

func (s *ComputerExtensionAttributesService) Create(ctx context.Context, computerExtensionAttribute *ComputerExtensionAttribute) (*int, *jamf.Response, error) {
	if computerExtensionAttribute == nil {
		return nil, nil, errors.New("ComputerExtensionAttributesService.Create(): cannot create nil computer extension attribute")
	}
	if computerExtensionAttribute.Name == nil {
		return nil, nil, errors.New("ComputerExtensionAttributesService.Create(): cannot create computer extension attribute with nil Name")
	}

	reqBody := &struct {
		*ComputerExtensionAttribute
		XMLName xml.Name `xml:"computer_extension_attribute"`
	}{
		ComputerExtensionAttribute: computerExtensionAttribute,
	}
	body, err := xml.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("xml.Marshal(): %v", err)
	}

	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: jamf.Uri{
			// When the ID is 0, the ID will be generated automatically at the server side.
			Entity: path.Join(computerExtensionAttributesPath, "id", "0"),
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Post(): %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.client.Logger.Printf("Error closing response body: %v", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		XMLName xml.Name `xml:"computer_extension_attribute"`
		ID      *int     `xml:"id"`
	}
	if err := xml.Unmarshal(respBody, &data); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	// The create API only returns the ID of the new ComputerExtensionAttribute.
	return data.ID, resp, nil
}

func (s *ComputerExtensionAttributesService) Delete(ctx context.Context, computerExtensionAttributeID int) (*jamf.Response, error) {
	resp, _, err := s.client.Delete(ctx, jamf.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(computerExtensionAttributesPath, "id", fmt.Sprint(computerExtensionAttributeID)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("client.Delete(): %v", err)
	}

	return resp, nil
}

func (s *ComputerExtensionAttributesService) Get(ctx context.Context, computerExtensionAttributeID int) (*ComputerExtensionAttribute, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(computerExtensionAttributesPath, "id", fmt.Sprint(computerExtensionAttributeID)),
		},
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Get(): %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.client.Logger.Printf("Error closing response body: %v", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var computerExtensionAttribute ComputerExtensionAttribute
	if err := xml.Unmarshal(respBody, &computerExtensionAttribute); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	return &computerExtensionAttribute, resp, nil
}

func (s *ComputerExtensionAttributesService) List(ctx context.Context) (*ListComputerExtensionAttributes, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: computerExtensionAttributesPath,
		},
	})

	if err != nil {
		return nil, nil, fmt.Errorf("client.Get(): %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.client.Logger.Printf("Error closing response body: %v", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var listComputerExtensionAttributes ListComputerExtensionAttributes
	if err := xml.Unmarshal(respBody, &listComputerExtensionAttributes); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	return &listComputerExtensionAttributes, resp, nil
}

func (s *ComputerExtensionAttributesService) Update(ctx context.Context, computerExtensionAttribute *ComputerExtensionAttribute) (*jamf.Response, error) {
	if computerExtensionAttribute == nil {
		return nil, errors.New("ComputerExtensionAttributesService.Update(): cannot create nil computer extension attribute")
	}
	if computerExtensionAttribute.Name == nil {
		return nil, errors.New("ComputerExtensionAttributesService.Update(): cannot update computer extension attribute with nil Name")
	}

	reqBody := &struct {
		*ComputerExtensionAttribute
		XMLName xml.Name `xml:"computer_extension_attribute"`
	}{
		ComputerExtensionAttribute: computerExtensionAttribute,
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
			Entity: path.Join(computerExtensionAttributesPath, "id", fmt.Sprint(*computerExtensionAttribute.ID)),
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, fmt.Errorf("client.Put(): %v", err)
	}

	return resp, nil
}
