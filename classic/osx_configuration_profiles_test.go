package classic

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOSXConfigurationProfilesService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(osxConfigurationProfilesPath, "id", "0"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, []byte("<os_x_configuration_profile><general><name>Test OSX Configuration Profile</name></general></os_x_configuration_profile>"))

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<os_x_configuration_profile>
  <id>1</id>
</os_x_configuration_profile>`))
	})

	ctx := context.Background()
	osxConfigurationProfileID, _, err := client.OSXConfigurationProfiles.Create(ctx, &OSXConfigurationProfile{
		General: &OSXConfigurationProfileGeneral{
			Name: ptr("Test OSX Configuration Profile"),
		},
	})
	if err != nil {
		t.Errorf("OSXConfigurationProfiles.Create(): %v", err)
	}

	want := ptr(1)
	if !cmp.Equal(osxConfigurationProfileID, want) {
		t.Errorf("OSXConfigurationProfiles.Create() returned %s, want %s", formatWithSpew(osxConfigurationProfileID), formatWithSpew(want))
	}
}

func TestOSXConfigurationProfilesService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(osxConfigurationProfilesPath, "id", "1"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<os_x_configuration_profile>
  <id>1</id>
</os_x_configuration_profile>`))
	})

	ctx := context.Background()
	osxConfigurationProfileID := 1
	_, err := client.OSXConfigurationProfiles.Delete(ctx, osxConfigurationProfileID)
	if err != nil {
		t.Errorf("OSXConfigurationProfiles.Delete(): %v", err)
	}
}

func TestOSXConfigurationProfilesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	osxConfigurationProfileID := 1
	mux.HandleFunc(buildHandlePath(osxConfigurationProfilesPath, "id", fmt.Sprint(osxConfigurationProfileID)), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<os_x_configuration_profile>
  <general>
    <id>1</id>
    <name>Test OSX Configuration Profile Name</name>
    <description>Test OSX Configuration Profile Description</description>
    <site>
      <id>-1</id>
      <name>None</name>
    </site>
    <category>
      <id>-1</id>
      <name>No category assigned</name>
    </category>
    <distribution_method>Make Available in Self Service</distribution_method>
    <user_removable>false</user_removable>
    <level>System</level>
    <uuid>DAF31D34-0650-4B57-88C3-0F75F8F56464</uuid>
    <redeploy_on_update>Newly Assigned</redeploy_on_update>
    <payloads>&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;&lt;!DOCTYPE plist PUBLIC &quot;-//Apple//DTD PLIST 1.0//EN&quot; &quot;http://www.apple.com/DTDs/PropertyList-1.0.dtd&quot;&gt;
&lt;plist version=&quot;1&quot;&gt;&lt;dict&gt;&lt;key&gt;PayloadUUID&lt;/key&gt;&lt;string&gt;DAF31D34-0650-4B57-88C3-0F75F8F56464&lt;/string&gt;&lt;key&gt;PayloadType&lt;/key&gt;&lt;string&gt;Configuration&lt;/string&gt;&lt;key&gt;PayloadOrganization&lt;/key&gt;&lt;string&gt;Test Organization&lt;/string&gt;&lt;key&gt;PayloadIdentifier&lt;/key&gt;&lt;string&gt;DAF31D34-0650-4B57-88C3-0F75F8F56464&lt;/string&gt;&lt;key&gt;PayloadDisplayName&lt;/key&gt;&lt;string&gt;Test OSX Configuration Profile&lt;/string&gt;&lt;key&gt;PayloadDescription&lt;/key&gt;&lt;string/&gt;&lt;key&gt;PayloadVersion&lt;/key&gt;&lt;integer&gt;1&lt;/integer&gt;&lt;key&gt;PayloadEnabled&lt;/key&gt;&lt;true/&gt;&lt;key&gt;PayloadRemovalDisallowed&lt;/key&gt;&lt;true/&gt;&lt;key&gt;PayloadScope&lt;/key&gt;&lt;string&gt;System&lt;/string&gt;&lt;key&gt;PayloadContent&lt;/key&gt;&lt;array&gt;&lt;dict&gt;&lt;key&gt;PayloadDisplayName&lt;/key&gt;&lt;string&gt;Passcode Payload&lt;/string&gt;&lt;key&gt;PayloadIdentifier&lt;/key&gt;&lt;string&gt;CC4BDA54-066F-42E9-B1F1-C8906B3FBF67&lt;/string&gt;&lt;key&gt;PayloadOrganization&lt;/key&gt;&lt;string&gt;JAMF Software&lt;/string&gt;&lt;key&gt;PayloadType&lt;/key&gt;&lt;string&gt;com.apple.mobiledevice.passwordosxConfigurationProfile&lt;/string&gt;&lt;key&gt;PayloadUUID&lt;/key&gt;&lt;string&gt;CC4BDA54-066F-42E9-B1F1-C8906B3FBF67&lt;/string&gt;&lt;key&gt;PayloadVersion&lt;/key&gt;&lt;integer&gt;1&lt;/integer&gt;&lt;key&gt;forcePIN&lt;/key&gt;&lt;true/&gt;&lt;/dict&gt;&lt;/array&gt;&lt;/dict&gt;&lt;/plist&gt;</payloads>
  </general>
  <scope>
    <all_computers>false</all_computers>
    <all_jss_users>false</all_jss_users>
    <computers>
      <computer>
        <id>1</id>
        <name>Test Computer</name>
        <udid>ACFE2853-D815-AEB2-3318-F4B7931934D1</udid>
      </computer>
    </computers>
    <buildings/>
    <departments/>
    <computer_groups>
      <computer_group>
        <id>1</id>
        <name>Test Computer Group</name>
      </computer_group>
    </computer_groups>
    <jss_users>
      <user>
        <id>1</id>
        <name>Test User</name>
      </user>
    </jss_users>
    <jss_user_groups>
      <user_group>
        <id>1</id>
        <name>Test User Group</name>
      </user_group>
    </jss_user_groups>
    <limitations>
      <users>
        <user>
          <name>riemann</name>
        </user>
      </users>
      <user_groups>
        <user_group>
          <name>scientists</name>
        </user_group>
      </user_groups>
      <network_segments>
        <network_segment>
          <id>1</id>
          <uid>43_1</uid>
          <name>Test Network Segment</name>
        </network_segment>
      </network_segments>
      <ibeacons/>
    </limitations>
    <exclusions>
      <computers/>
      <buildings/>
      <departments/>
      <computer_groups/>
      <users/>
      <user_groups/>
      <network_segments/>
      <ibeacons/>
      <jss_users/>
      <jss_user_groups/>
    </exclusions>
  </scope>
  <self_service>
    <self_service_display_name>Test OSX Configuration Profile Name</self_service_display_name>
    <install_button_text>Install</install_button_text>
    <self_service_description>Test OSX Configuration Profile Description</self_service_description>
    <force_users_to_view_description>true</force_users_to_view_description>
    <security>
      <removal_disallowed>Never</removal_disallowed>
    </security>
    <self_service_icon/>
    <feature_on_main_page>true</feature_on_main_page>
    <self_service_categories>
      <category>
        <id>1</id>
        <name>Test Category</name>
        <display_in>true</display_in>
        <feature_in>true</feature_in>
      </category>
      <category>
        <id>2</id>
        <name>Test Category Second</name>
        <display_in>true</display_in>
        <feature_in>false</feature_in>
      </category>
    </self_service_categories>
    <notification>false</notification>
    <notification>Self Service</notification>
    <notification_subject/>
    <notification_message/>
  </self_service>
</os_x_configuration_profile>`))
	})

	ctx := context.Background()
	osxConfigurationProfile, _, err := client.OSXConfigurationProfiles.Get(ctx, osxConfigurationProfileID)
	if err != nil {
		t.Errorf("OSXConfigurationProfiles.Get(): %v", err)
	}

	want := &OSXConfigurationProfile{
		General: &OSXConfigurationProfileGeneral{
			ID:          ptr(1),
			Name:        ptr("Test OSX Configuration Profile Name"),
			Description: ptr("Test OSX Configuration Profile Description"),
			Site: &Site{
				ID:   ptr(-1),
				Name: ptr("None"),
			},
			Category: &GeneralCategory{
				ID:   ptr(-1),
				Name: ptr("No category assigned"),
			},
			DistributionMethod: ptr(OSXConfigurationProfileGeneralDistributionMethodMakeAvailableInSelfService),
			UserRemovable:      ptr(false),
			Level:              ptr(OSXConfigurationProfileGeneralLevelSystem),
			UUID:               ptr("DAF31D34-0650-4B57-88C3-0F75F8F56464"),
			RedeployOnUpdate:   ptr(OSXConfigurationProfileGeneralRedeployOnUpdateNewlyAssigned),
			Payloads: ptr(`<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1"><dict><key>PayloadUUID</key><string>DAF31D34-0650-4B57-88C3-0F75F8F56464</string><key>PayloadType</key><string>Configuration</string><key>PayloadOrganization</key><string>Test Organization</string><key>PayloadIdentifier</key><string>DAF31D34-0650-4B57-88C3-0F75F8F56464</string><key>PayloadDisplayName</key><string>Test OSX Configuration Profile</string><key>PayloadDescription</key><string/><key>PayloadVersion</key><integer>1</integer><key>PayloadEnabled</key><true/><key>PayloadRemovalDisallowed</key><true/><key>PayloadScope</key><string>System</string><key>PayloadContent</key><array><dict><key>PayloadDisplayName</key><string>Passcode Payload</string><key>PayloadIdentifier</key><string>CC4BDA54-066F-42E9-B1F1-C8906B3FBF67</string><key>PayloadOrganization</key><string>JAMF Software</string><key>PayloadType</key><string>com.apple.mobiledevice.passwordosxConfigurationProfile</string><key>PayloadUUID</key><string>CC4BDA54-066F-42E9-B1F1-C8906B3FBF67</string><key>PayloadVersion</key><integer>1</integer><key>forcePIN</key><true/></dict></array></dict></plist>`),
		},
		Scope: &OSXConfigurationProfileScope{
			AllComputers: ptr(false),
			AllJSSUsers:  ptr(false),
			Computers: &[]OSXConfigurationProfileScopeComputer{{
				ID:   ptr(1),
				Name: ptr("Test Computer"),
				UDID: ptr("ACFE2853-D815-AEB2-3318-F4B7931934D1"),
			}},
			ComputerGroups: &[]OSXConfigurationProfileScopeComputerGroup{{
				ID:   ptr(1),
				Name: ptr("Test Computer Group"),
			}},
			JSSUsers: &[]OSXConfigurationProfileScopeUser{{
				ID:   ptr(1),
				Name: ptr("Test User"),
			}},
			JSSUserGroups: &[]OSXConfigurationProfileScopeUserGroup{{
				ID:   ptr(1),
				Name: ptr("Test User Group"),
			}},
			Limitations: &OSXConfigurationProfileScopeLimitations{
				Users: &[]OSXConfigurationProfileScopeLimitationsUser{{
					Name: ptr("riemann"),
				}},
				UserGroups: &[]OSXConfigurationProfileScopeLimitationsUserGroup{{
					Name: ptr("scientists"),
				}},
				NetworkSegments: &[]OSXConfigurationProfileScopeNetworkSegment{{
					ID:   ptr(1),
					Name: ptr("Test Network Segment"),
				}},
			},
			Exclusions: &OSXConfigurationProfileScopeExclusions{},
		},
		SelfService: &OSXConfigurationProfileSelfService{
			SelfServiceDisplayName:      ptr("Test OSX Configuration Profile Name"),
			InstallButtonText:           ptr("Install"),
			SelfServiceDescription:      ptr("Test OSX Configuration Profile Description"),
			ForceUsersToViewDescription: ptr(true),
			Security: &OSXConfigurationProfileSelfServiceSecurity{
				RemovalDisallowed: ptr(OSXConfigurationProfileSelfServiceSecurityRemovalDisallowedNever),
			},
			FeatureOnMainPage: ptr(true),
			SelfServiceCategories: &[]SelfServiceCategory{
				{
					ID:        ptr(1),
					Name:      ptr("Test Category"),
					DisplayIn: ptr(true),
					FeatureIn: ptr(true),
				},
				{
					ID:        ptr(2),
					Name:      ptr("Test Category Second"),
					DisplayIn: ptr(true),
					FeatureIn: ptr(false),
				},
			},
		},
	}
	if !cmp.Equal(osxConfigurationProfile, want) {
		t.Errorf("OSXConfigurationProfiles.Get() returned %s, want %s", formatWithSpew(osxConfigurationProfile), formatWithSpew(want))
	}
}

func TestOSXConfigurationProfilesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(osxConfigurationProfilesPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<policies>
  <size>1</size>
  <os_x_configuration_profile>
    <id>1</id>
    <name>Test OSX Configuration Profile</name>
  </os_x_configuration_profile>
</policies>`))
	})

	ctx := context.Background()
	policies, _, err := client.OSXConfigurationProfiles.List(ctx)
	if err != nil {
		t.Errorf("OSXConfigurationProfiles.List(): %v", err)
	}

	want := &ListOSXConfigurationProfiles{
		Size: ptr(1),
		OSXConfigurationProfiles: &[]ListOSXConfigurationProfile{{
			ID:   ptr(1),
			Name: ptr("Test OSX Configuration Profile"),
		}},
	}
	if !cmp.Equal(policies, want) {
		t.Errorf("OSXConfigurationProfiles.List() returned %s, want %s", formatWithSpew(policies), formatWithSpew(want))
	}
}

func TestOSXConfigurationProfilesService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	osxConfigurationProfileID := 1

	mux.HandleFunc(buildHandlePath(osxConfigurationProfilesPath, "id", fmt.Sprint(osxConfigurationProfileID)), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, []byte("<os_x_configuration_profile><general><id>1</id><name>Test OSX Configuration Profile Updated</name></general></os_x_configuration_profile>"))

		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	osxConfigurationProfile := &OSXConfigurationProfile{
		General: &OSXConfigurationProfileGeneral{
			ID:   ptr(osxConfigurationProfileID),
			Name: ptr("Test OSX Configuration Profile Updated"),
		},
	}
	_, err := client.OSXConfigurationProfiles.Update(ctx, osxConfigurationProfile)
	if err != nil {
		t.Errorf("OSXConfigurationProfiles.Update(): %v", err)
	}
}
