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

type PackagesService service

type Package struct {
	ID                         *int    `xml:"id,omitempty"`
	Name                       *string `xml:"name,omitempty"`
	Category                   *string `xml:"category,omitempty"`
	Filename                   *string `xml:"filename,omitempty"`
	Info                       *string `xml:"info,omitempty"`
	Notes                      *string `xml:"notes,omitempty"`
	Priority                   *int    `xml:"priority,omitempty"`
	RebootRequired             *bool   `xml:"reboot_required,omitempty"`
	FillUserTemplate           *bool   `xml:"fill_user_template,omitempty"`
	FillExistingUsers          *bool   `xml:"fill_existing_users,omitempty"`
	AllowUninstalled           *bool   `xml:"allow_uninstalled,omitempty"`
	OSRequirements             *string `xml:"os_requirements,omitempty"`
	RequiredProcessor          *string `xml:"required_processor,omitempty"`
	HashType                   *string `xml:"hash_type,omitempty"`
	HashValue                  *string `xml:"hash_value,omitempty"`
	SwitchWithPackage          *string `xml:"switch_with_package,omitempty"`
	InstallIfReportedAvailable *bool   `xml:"install_if_reported_available,omitempty"`
	ReinstallOption            *string `xml:"reinstall_option,omitempty"`
	TriggeringFiles            *string `xml:"triggering_files,omitempty"`
	SendNotification           *bool   `xml:"send_notification,omitempty"`
}

type PackageRequiredProcessor string

const (
	PackageRequiredProcessorNone PackageRequiredProcessor = "None"
	PackageRequiredProcessorPpc  PackageRequiredProcessor = "ppc"
	PackageRequiredProcessorX86  PackageRequiredProcessor = "x86"
)

type ListPackages struct {
	Size     *int           `xml:"size,omitempty"`
	Packages *[]ListPackage `xml:"package,omitempty"`
}

type ListPackage struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

const packagesPath = "/packages"

func (s *PackagesService) Create(ctx context.Context, pkg *Package) (*int, *jamf.Response, error) {
	if pkg == nil {
		return nil, nil, errors.New("PackagesService.Create(): cannot create nil package")
	}
	if pkg.Name == nil {
		return nil, nil, errors.New("PackagesService.Create(): cannot create package with nil Name")
	}
	if pkg.Filename == nil {
		return nil, nil, errors.New("PackagesService.Create(): cannot create package with nil Filename")
	}

	reqBody := &struct {
		*Package
		XMLName xml.Name `xml:"package"`
	}{
		Package: pkg,
	}
	body, err := xml.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("xml.Marshal(): %v", err)
	}

	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: jamf.Uri{
			// When the ID is 0, the ID will be generated automatically at the server side.
			Entity: path.Join(packagesPath, "id", "0"),
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
		XMLName xml.Name `xml:"package"`
		ID      *int     `xml:"id"`
	}
	if err := xml.Unmarshal(respBody, &data); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	// The create API only returns the ID of the new Package.
	return data.ID, resp, nil
}

func (s *PackagesService) Delete(ctx context.Context, packageID int) (*jamf.Response, error) {
	resp, _, err := s.client.Delete(ctx, jamf.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(packagesPath, "id", fmt.Sprint(packageID)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("client.Delete(): %v", err)
	}

	return resp, nil
}

func (s *PackagesService) Get(ctx context.Context, packageID int) (*Package, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(packagesPath, "id", fmt.Sprint(packageID)),
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

	var pkg Package
	if err := xml.Unmarshal(respBody, &pkg); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	return &pkg, resp, nil
}

func (s *PackagesService) List(ctx context.Context) (*ListPackages, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: packagesPath,
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

	var listPackages ListPackages
	if err := xml.Unmarshal(respBody, &listPackages); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	return &listPackages, resp, nil
}

func (s *PackagesService) Update(ctx context.Context, pkg *Package) (*jamf.Response, error) {
	if pkg == nil {
		return nil, errors.New("PackagesService.Update(): cannot update nil package")
	}

	reqBody := &struct {
		*Package
		XMLName xml.Name `xml:"package"`
	}{
		Package: pkg,
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
			Entity: path.Join(packagesPath, "id", fmt.Sprint(*pkg.ID)),
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, fmt.Errorf("client.Put(): %v", err)
	}

	return resp, nil
}
