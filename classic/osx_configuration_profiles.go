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

type OSXConfigurationProfilesService service

type ListOSXConfigurationProfiles struct {
	Size                     *int                           `xml:"size,omitempty"`
	OSXConfigurationProfiles *[]ListOSXConfigurationProfile `xml:"os_x_configuration_profile,omitempty"`
}

type ListOSXConfigurationProfile struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type OSXConfigurationProfile struct {
	General     *OSXConfigurationProfileGeneral     `xml:"general,omitempty"`
	Scope       *OSXConfigurationProfileScope       `xml:"scope,omitempty"`
	SelfService *OSXConfigurationProfileSelfService `xml:"self_service,omitempty"`
}

type OSXConfigurationProfileGeneral struct {
	ID                 *int                                              `xml:"id,omitempty"`
	Name               *string                                           `xml:"name,omitempty"`
	Description        *string                                           `xml:"description,omitempty"`
	Site               *Site                                             `xml:"site,omitempty"`
	Category           *GeneralCategory                                  `xml:"category,omitempty"`
	DistributionMethod *OSXConfigurationProfileGeneralDistributionMethod `xml:"distribution_method,omitempty"`
	UserRemovable      *bool                                             `xml:"user_removable,omitempty"`
	Level              *OSXConfigurationProfileGeneralLevel              `xml:"level,omitempty"`
	UUID               *string                                           `xml:"uuid,omitempty"`
	RedeployOnUpdate   *OSXConfigurationProfileGeneralRedeployOnUpdate   `xml:"redeploy_on_update,omitempty"`
	Payloads           *string                                           `xml:"payloads,omitempty"`
}

type OSXConfigurationProfileGeneralDistributionMethod string

const (
	OSXConfigurationProfileGeneralDistributionMethodInstallAutomatically       OSXConfigurationProfileGeneralDistributionMethod = "Install Automatically"
	OSXConfigurationProfileGeneralDistributionMethodMakeAvailableInSelfService OSXConfigurationProfileGeneralDistributionMethod = "Make Available in Self Service"
)

type OSXConfigurationProfileGeneralLevel string

const (
	OSXConfigurationProfileGeneralLevelSystem OSXConfigurationProfileGeneralLevel = "System"
	OsxConfigurationProfileGeneralLevelUser   OSXConfigurationProfileGeneralLevel = "User"
)

type OSXConfigurationProfileGeneralRedeployOnUpdate string

const (
	OSXConfigurationProfileGeneralRedeployOnUpdateNewlyAssigned OSXConfigurationProfileGeneralRedeployOnUpdate = "Newly Assigned"
	OSXConfigurationProfileGeneralRedeployOnUpdateAll           OSXConfigurationProfileGeneralRedeployOnUpdate = "All"
)

type OSXConfigurationProfileScope struct {
	AllComputers   *bool                                        `xml:"all_computers,omitempty"`
	AllJSSUsers    *bool                                        `xml:"all_jss_users,omitempty"`
	Computers      *[]OSXConfigurationProfileScopeComputer      `xml:"computers>computer,omitempty"`
	ComputerGroups *[]OSXConfigurationProfileScopeComputerGroup `xml:"computer_groups>computer_group,omitempty"`
	Buildings      *[]Building                                  `xml:"buildings>building,omitempty"`
	Departments    *[]Department                                `xml:"departments>department,omitempty"`
	JSSUsers       *[]OSXConfigurationProfileScopeUser          `xml:"jss_users>user,omitempty"`
	JSSUserGroups  *[]OSXConfigurationProfileScopeUserGroup     `xml:"jss_user_groups>user_group,omitempty"`
	Limitations    *OSXConfigurationProfileScopeLimitations     `xml:"limitations,omitempty"`
	Exclusions     *OSXConfigurationProfileScopeExclusions      `xml:"exclusions,omitempty"`
}

type OSXConfigurationProfileScopeComputer struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
	UDID *string `xml:"udid,omitempty"`
}

type OSXConfigurationProfileScopeComputerGroup struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type OSXConfigurationProfileScopeExclusions struct {
	Computers       *[]OSXConfigurationProfileScopeComputer       `xml:"computers>computer,omitempty"`
	ComputerGroups  *[]OSXConfigurationProfileScopeComputerGroup  `xml:"computer_groups>computer_group,omitempty"`
	Buildings       *[]Building                                   `xml:"buildings>building,omitempty"`
	Departments     *[]Department                                 `xml:"departments>department,omitempty"`
	Users           *[]OSXConfigurationProfileScopeUser           `xml:"jss_users>user,omitempty"`
	UserGroups      *[]OSXConfigurationProfileScopeUserGroup      `xml:"jss_user_groups>user_group,omitempty"`
	NetworkSegments *[]OSXConfigurationProfileScopeNetworkSegment `xml:"network_segments>network_segment,omitempty"`
	Ibeacons        *[]OSXConfigurationProfileScopeIbeacon        `xml:"ibeacons>ibeacon,omitempty"`
}

type OSXConfigurationProfileScopeLimitations struct {
	Users           *[]OSXConfigurationProfileScopeLimitationsUser      `xml:"users>user,omitempty"`
	UserGroups      *[]OSXConfigurationProfileScopeLimitationsUserGroup `xml:"user_groups>user_group,omitempty"`
	NetworkSegments *[]OSXConfigurationProfileScopeNetworkSegment       `xml:"network_segments>network_segment,omitempty"`
	Ibeacons        *[]OSXConfigurationProfileScopeIbeacon              `xml:"ibeacons>ibeacon,omitempty"`
}

type OSXConfigurationProfileScopeLimitationsUser struct {
	Name *string `xml:"name,omitempty"`
}

type OSXConfigurationProfileScopeLimitationsUserGroup struct {
	Name *string `xml:"name,omitempty"`
}

type OSXConfigurationProfileScopeIbeacon struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type OSXConfigurationProfileScopeNetworkSegment struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type OSXConfigurationProfileScopeUser struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type OSXConfigurationProfileScopeUserGroup struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type OSXConfigurationProfileSelfService struct {
	// The 'self_service_display_name' filed is not documented, but it actually exists in the API response.
	SelfServiceDisplayName      *string `xml:"self_service_display_name,omitempty"`
	InstallButtonText           *string `xml:"install_button_text,omitempty"`
	SelfServiceDescription      *string `xml:"self_service_description,omitempty"`
	ForceUsersToViewDescription *bool   `xml:"force_users_to_view_description,omitempty"`
	// The 'security' filed is not documented, but it actually exists in the API response.
	Security *OSXConfigurationProfileSelfServiceSecurity `xml:"security,omitempty"`
	// The attributes related to SelfServiceIcon are not working.
	//SelfServiceIcon *interface{} `xml:"self_service_icon,omitempty"`
	FeatureOnMainPage     *bool                  `xml:"feature_on_main_page,omitempty"`
	SelfServiceCategories *[]SelfServiceCategory `xml:"self_service_categories>category,omitempty"`
	// The attributes related to Notification are not working.
	//NotificationEnabled   *bool                                               `xml:"-"`
	//NotificationType      *OSXConfigurationProfileSelfServiceNotificationType `xml:"-"`
	//NotificationSubject   *string                                             `xml:"notification_subject,omitempty"`
	//NotificationMessage   *string                                             `xml:"notification_message,omitempty"`
}

type OSXConfigurationProfileSelfServiceNotificationType string

const (
	OSXConfigurationProfileSelfServiceNotificationTypeSelfService                      OSXConfigurationProfileSelfServiceNotificationType = "Self Service"
	OSXConfigurationProfileSelfServiceNotificationTypeSelfServiceAndNotificationCenter OSXConfigurationProfileSelfServiceNotificationType = "Self Service and Notification Center"
)

type OSXConfigurationProfileSelfServiceSecurity struct {
	RemovalDisallowed *OSXConfigurationProfileSelfServiceSecurityRemovalDisallowed `xml:"removal_disallowed,omitempty"`
}

type OSXConfigurationProfileSelfServiceSecurityRemovalDisallowed string

const (
	OSXConfigurationProfileSelfServiceSecurityRemovalDisallowedAlways            OSXConfigurationProfileSelfServiceSecurityRemovalDisallowed = "Always"
	OSXConfigurationProfileSelfServiceSecurityRemovalDisallowedNever             OSXConfigurationProfileSelfServiceSecurityRemovalDisallowed = "Never"
	OSXConfigurationProfileSelfServiceSecurityRemovalDisallowedWithAuthorization OSXConfigurationProfileSelfServiceSecurityRemovalDisallowed = "With Authorization"
)

const osxConfigurationProfilesPath = "/osxconfigurationprofiles"

func (s *OSXConfigurationProfilesService) Create(ctx context.Context, osxConfigurationProfile *OSXConfigurationProfile) (*int, *jamf.Response, error) {
	if osxConfigurationProfile == nil {
		return nil, nil, errors.New("OSXConfigurationProfilesService.Create(): cannot create nil OSX configuration profile")
	}
	if osxConfigurationProfile.General == nil {
		return nil, nil, errors.New("OSXConfigurationProfilesService.Create(): cannot create OSX configuration profile with nil General")
	}
	if osxConfigurationProfile.General.Name == nil {
		return nil, nil, errors.New("OSXConfigurationProfilesService.Create(): cannot create OSX configuration profile with nil Name of General")
	}

	reqBody := &struct {
		*OSXConfigurationProfile
		XMLName xml.Name `xml:"os_x_configuration_profile"`
	}{
		OSXConfigurationProfile: osxConfigurationProfile,
	}
	body, err := xml.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("xml.Marshal(): %v", err)
	}

	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: jamf.Uri{
			// When the ID is 0, the ID will be generated automatically at the server side.
			Entity: path.Join(osxConfigurationProfilesPath, "id", "0"),
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
		XMLName xml.Name `xml:"os_x_configuration_profile"`
		ID      *int     `xml:"id"`
	}
	if err := xml.Unmarshal(respBody, &data); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	// The create API only returns the ID of the new OSXConfigurationProfile.
	return data.ID, resp, nil
}

func (s *OSXConfigurationProfilesService) Delete(ctx context.Context, osxConfigurationProfileID int) (*jamf.Response, error) {
	resp, _, err := s.client.Delete(ctx, jamf.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(osxConfigurationProfilesPath, "id", fmt.Sprint(osxConfigurationProfileID)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("client.Delete(): %v", err)
	}

	return resp, nil
}

func (s *OSXConfigurationProfilesService) Get(ctx context.Context, osxConfigurationProfileID int) (*OSXConfigurationProfile, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(osxConfigurationProfilesPath, "id", fmt.Sprint(osxConfigurationProfileID)),
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

	var osxConfigurationProfile OSXConfigurationProfile
	if err := xml.Unmarshal(respBody, &osxConfigurationProfile); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	return &osxConfigurationProfile, resp, nil
}

func (s *OSXConfigurationProfilesService) List(ctx context.Context) (*ListOSXConfigurationProfiles, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: osxConfigurationProfilesPath,
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

	var listOSXConfigurationProfiles ListOSXConfigurationProfiles
	if err := xml.Unmarshal(respBody, &listOSXConfigurationProfiles); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	return &listOSXConfigurationProfiles, resp, nil
}

func (s *OSXConfigurationProfilesService) Update(ctx context.Context, osxConfigurationProfile *OSXConfigurationProfile) (*jamf.Response, error) {
	if osxConfigurationProfile == nil {
		return nil, errors.New("OSXConfigurationProfilesService.Update(): cannot update nil OSX configuration profile")
	}
	if osxConfigurationProfile.General == nil {
		return nil, errors.New("OSXConfigurationProfilesService.Update(): cannot update OSX configuration profile with nil General")
	}
	if osxConfigurationProfile.General.Name == nil {
		return nil, errors.New("OSXConfigurationProfilesService.Create(): cannot update OSX configuration profile with nil Name of General")
	}

	reqBody := &struct {
		*OSXConfigurationProfile
		XMLName xml.Name `xml:"os_x_configuration_profile"`
	}{
		OSXConfigurationProfile: osxConfigurationProfile,
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
			Entity: path.Join(osxConfigurationProfilesPath, "id", fmt.Sprint(*osxConfigurationProfile.General.ID)),
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, fmt.Errorf("client.Put(): %v", err)
	}

	return resp, nil
}
