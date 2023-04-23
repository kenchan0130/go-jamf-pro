package classic

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestPoliciesService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(policiesPath, "id", "0"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, []byte("<policy><general><name>Test Policy</name></general></policy>"))

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<policy>
  <id>1</id>
</policy>`))
	})

	ctx := context.Background()
	policyID, _, err := client.Policies.Create(ctx, &Policy{
		General: &PolicyGeneral{
			Name: ptr("Test Policy"),
		},
	})
	if err != nil {
		t.Errorf("Policies.Create(): %v", err)
	}

	want := ptr(1)
	if !cmp.Equal(policyID, want) {
		t.Errorf("Policies.Create() returned %s, want %s", formatWithSpew(policyID), formatWithSpew(want))
	}
}

func TestPoliciesService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(policiesPath, "id", "1"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<policy>
  <id>1</id>
</policy>`))
	})

	ctx := context.Background()
	policyID := 1
	_, err := client.Policies.Delete(ctx, policyID)
	if err != nil {
		t.Errorf("Policies.Delete(): %v", err)
	}
}

func TestPoliciesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	policyID := 1
	mux.HandleFunc(buildHandlePath(policiesPath, "id", fmt.Sprint(policyID)), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<policy>
  <general>
    <id>1</id>
    <name>Test Policy Name</name>
    <enabled>true</enabled>
    <trigger>EVENT</trigger>
    <trigger_checkin>false</trigger_checkin>
    <trigger_enrollment_complete>false</trigger_enrollment_complete>
    <trigger_login>false</trigger_login>
    <trigger_network_state_changed>false</trigger_network_state_changed>
    <trigger_startup>false</trigger_startup>
    <trigger_other>Test Trigger</trigger_other>
    <frequency>Once every day</frequency>
    <retry_event>none</retry_event>
    <retry_attempts>-1</retry_attempts>
    <notify_on_each_failed_retry>false</notify_on_each_failed_retry>
    <location_user_only>false</location_user_only>
    <target_drive>/</target_drive>
    <offline>false</offline>
    <category>
      <id>1</id>
      <name>Test Category</name>
    </category>
    <date_time_limitations>
      <activation_date>2023-01-01 14:00:00</activation_date>
      <activation_date_epoch>1672581600000</activation_date_epoch>
      <activation_date_utc>2023-01-01T14:00:00.000+0000</activation_date_utc>
      <expiration_date>2027-01-07 12:01:00</expiration_date>
      <expiration_date_epoch>1799323260000</expiration_date_epoch>
      <expiration_date_utc>2027-01-07T12:01:00.000+0000</expiration_date_utc>
      <no_execute_on>
        <day>Sun</day>
        <day>Mon</day>
      </no_execute_on>
      <no_execute_start>10:00 AM</no_execute_start>
      <no_execute_end>1:00 PM</no_execute_end>
    </date_time_limitations>
    <network_limitations>
      <minimum_network_connection>Ethernet</minimum_network_connection>
      <any_ip_address>false</any_ip_address>
      <network_segments>
        <network_segment>
          <id>1</id>
          <name>Test Network Segment</name>
        </network_segment>
      </network_segments>
    </network_limitations>
    <override_default_settings>
      <target_drive>/</target_drive>
      <distribution_point>default</distribution_point>
      <force_afp_smb>false</force_afp_smb>
      <sus>default</sus>
    </override_default_settings>
    <network_requirements>Any</network_requirements>
    <site>
      <id>-1</id>
      <name>None</name>
    </site>
  </general>
  <scope>
    <all_computers>false</all_computers>
    <computers/>
    <computer_groups/>
    <buildings/>
    <departments/>
    <limit_to_users>
      <user_groups>
        <user_group>scientists</user_group>
      </user_groups>
    </limit_to_users>
    <limitations>
      <users>
        <user>
          <name>riemann</name>
        </user>
        <user>
          <name>gauss</name>
        </user>
      </users>
      <user_groups>
        <user_group>
          <id>scientists</id>
          <name>scientists</name>
        </user_group>
      </user_groups>
      <network_segments>
        <network_segment>
          <id>1</id>
          <name>Test Network Segment</name>
        </network_segment>
      </network_segments>
      <ibeacons/>
    </limitations>
    <exclusions>
      <computers/>
      <computer_groups/>
      <buildings/>
      <departments/>
      <users/>
      <user_groups/>
      <network_segments/>
      <ibeacons/>
    </exclusions>
  </scope>
  <self_service>
    <use_for_self_service>false</use_for_self_service>
    <self_service_display_name>Test Policy</self_service_display_name>
    <install_button_text>Install</install_button_text>
    <reinstall_button_text>Reinstall</reinstall_button_text>
    <self_service_description>Test Description</self_service_description>
    <force_users_to_view_description>false</force_users_to_view_description>
    <self_service_icon>
      <id>1</id>
      <filename>test.png</filename>
      <uri>https://jamf.jamfcloud.com//api/v1/icon/download/1</uri>
    </self_service_icon>
    <feature_on_main_page>false</feature_on_main_page>
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
    <notification>true</notification>
    <notification>Self Service</notification>
    <notification_subject>Test Notification Subject</notification_subject>
    <notification_message>Test notification message</notification_message>
  </self_service>
  <package_configuration>
    <packages>
      <size>1</size>
      <package>
       <id>1</id>
       <name>test.dmg</name>
       <action>Install</action>
       <fut>false</fut>
       <feu>false</feu>
	  </package>
    </packages>
    <distribution_point>default</distribution_point>
  </package_configuration>
  <scripts>
    <size>1</size>
    <script>
      <id>1</id>
      <name>Test Script</name>
      <priority>After</priority>
      <parameter4>Test Parameter</parameter4>
      <parameter5/>
      <parameter6/>
      <parameter7/>
      <parameter8/>
      <parameter9/>
      <parameter10/>
      <parameter11/>
    </script>
  </scripts>
  <printers>
	<size>1</size>
    <leave_existing_default/>
    <printer>
      <id>1</id>
      <name>Test Printer</name>
      <action>install</action>
      <make_default>false</make_default>
    </printer>
  </printers>
  <dock_items>
    <size>1</size>
    <dock_item>
      <id>1</id>
      <name>Test Dock Item</name>
      <action>Add To Beginning</action>
    </dock_item>
  </dock_items>
  <account_maintenance>
    <accounts>
      <size>1</size>
      <account>
        <action>Create</action>
        <username>Test Username</username>
        <realname>Test Realname</realname>
        <password_sha256>********************</password_sha256>
        <home>/Users/test</home>
        <hint/>
        <picture/>
        <admin>false</admin>
        <filevault_enabled>false</filevault_enabled>
      </account>
    </accounts>
    <directory_bindings>
      <size>0</size>
    </directory_bindings>
    <management_account>
      <action>random</action>
      <managed_password_length>8</managed_password_length>
    </management_account>
    <open_firmware_efi_password>
      <of_mode>command</of_mode>
      <of_password_sha256>********************</of_password_sha256>
    </open_firmware_efi_password>
  </account_maintenance>
  <reboot>
    <message>This is a reboot message</message>
    <startup_disk>Specify Local Startup Disk</startup_disk>
    <specify_startup>/Volumes/test</specify_startup>
    <no_user_logged_in>Restart immediately</no_user_logged_in>
    <user_logged_in>Restart if a package or update requires it</user_logged_in>
    <minutes_until_reboot>5</minutes_until_reboot>
    <start_reboot_timer_immediately>true</start_reboot_timer_immediately>
    <file_vault_2_reboot>true</file_vault_2_reboot>
  </reboot>
  <maintenance>
    <recon>true</recon>
    <reset_name>false</reset_name>
    <install_all_cached_packages>false</install_all_cached_packages>
    <heal>false</heal>
    <prebindings>false</prebindings>
    <permissions>false</permissions>
    <byhost>false</byhost>
    <system_cache>false</system_cache>
    <user_cache>false</user_cache>
    <verify>false</verify>
  </maintenance>
  <files_processes>
    <search_by_path/>
    <delete_file>false</delete_file>
    <locate_file/>
    <update_locate_database>false</update_locate_database>
    <spotlight_search/>
    <search_for_process/>
    <kill_process>false</kill_process>
    <run_command/>
  </files_processes>
  <user_interaction>
    <message_start/>
    <allow_users_to_defer>true</allow_users_to_defer>
    <allow_deferral_until_utc>2025-01-01T10:00:00.000+0000</allow_deferral_until_utc>
    <allow_deferral_minutes>0</allow_deferral_minutes>
    <message_finish/>
  </user_interaction>
  <disk_encryption>
    <action>apply</action>
    <disk_encryption_configuration_id>1</disk_encryption_configuration_id>
    <auth_restart>false</auth_restart>
  </disk_encryption>
</policy>`))
	})

	ctx := context.Background()
	policy, _, err := client.Policies.Get(ctx, policyID)
	if err != nil {
		t.Errorf("Policies.Get(): %v", err)
	}

	want := &Policy{
		General: &PolicyGeneral{
			ID:                         ptr(1),
			Name:                       ptr("Test Policy Name"),
			Enabled:                    ptr(true),
			Trigger:                    ptr(PolicyGeneralTriggerEvent),
			TriggerCheckin:             ptr(false),
			TriggerEnrollmentComplete:  ptr(false),
			TriggerLogin:               ptr(false),
			TriggerNetworkStateChanged: ptr(false),
			TriggerStartup:             ptr(false),
			TriggerOther:               ptr("Test Trigger"),
			Frequency:                  ptr(PolicyGeneralFrequencyOnceEveryDay),
			RetryEvent:                 ptr(PolicyGeneralRetryEventNone),
			RetryAttempts:              ptr(-1),
			NotifyOnEachFailedRetry:    ptr(false),
			LocationUserOnly:           ptr(false),
			TargetDrive:                ptr("/"),
			Offline:                    ptr(false),
			Category: &GeneralCategory{
				ID:   ptr(1),
				Name: ptr("Test Category"),
			},
			DateTimeLimitations: &PolicyGeneralDateTimeLimitations{
				ActivationDate:      ptr(time.Date(2023, time.January, 1, 14, 0, 0, 0, time.UTC)),
				ActivationDateEpoch: ptr(int64(1672581600000)),
				ActivationDateUTC:   ptr(time.Date(2023, time.January, 1, 14, 0, 0, 0, time.UTC)),
				ExpirationDate:      ptr(time.Date(2027, time.January, 7, 12, 1, 0, 0, time.UTC)),
				ExpirationDateEpoch: ptr(int64(1799323260000)),
				ExpirationDateUTC:   ptr(time.Date(2027, time.January, 7, 12, 1, 0, 0, time.UTC)),
				NoExecuteOn:         &[]PolicyGeneralDateTimeLimitationsNoExecuteOnDay{PolicyGeneralDateTimeLimitationsNoExecuteOnDaySun, PolicyGeneralDateTimeLimitationsNoExecuteOnDayMon},
				NoExecuteStart:      ptr("10:00 AM"),
				NoExecuteEnd:        ptr("1:00 PM"),
			},
			NetworkLimitations: &PolicyGeneralNetworkLimitations{
				MinimumNetworkConnection: ptr(PolicyGeneralNetworkLimitationsMinimumNetworkConnectionEthernet),
				AnyIpAddress:             ptr(false),
			},
			OverrideDefaultSettings: &PolicyGeneralOverrideDefaultSettings{
				TargetDrive:       ptr("/"),
				DistributionPoint: ptr("default"),
				ForceAfpSmb:       ptr(false),
				Sus:               ptr("default"),
			},
			NetworkRequirements: ptr(PolicyGeneralNetworkRequirementsAny),
			Site: &Site{
				ID:   ptr(-1),
				Name: ptr("None"),
			},
		},
		Scope: &PolicyScope{
			AllComputers: ptr(false),
			Limitations: &PolicyScopeLimitations{
				Users:           &[]PolicyScopeUser{{Name: ptr("riemann")}, {Name: ptr("gauss")}},
				UserGroups:      &[]PolicyScopeUserGroup{{ID: ptr("scientists"), Name: ptr("scientists")}},
				NetworkSegments: &[]PolicyScopeNetworkSegment{{ID: ptr(1), Name: ptr("Test Network Segment")}},
			},
			Exclusions: &PolicyScopeExclusions{},
		},
		SelfService: &PolicySelfService{
			UseForSelfService:           ptr(false),
			SelfServiceDisplayName:      ptr("Test Policy"),
			InstallButtonText:           ptr("Install"),
			ReinstallButtonText:         ptr("Reinstall"),
			SelfServiceDescription:      ptr("Test Description"),
			ForceUsersToViewDescription: ptr(false),
			SelfServiceIcon: &Icon{
				ID:       ptr(1),
				Filename: ptr("test.png"),
				URI:      ptr("https://jamf.jamfcloud.com//api/v1/icon/download/1"),
			},
			FeatureOnMainPage: ptr(false),
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
			NotificationEnabled: ptr(true),
			NotificationType:    ptr(PolicySelfServiceNotificationTypeSelfService),
			NotificationSubject: ptr("Test Notification Subject"),
			NotificationMessage: ptr("Test notification message"),
		},
		PackageConfiguration: &PolicyPackageConfiguration{
			Packages: &PolicyPackageConfigurationPackages{
				Size: ptr(1),
				Packages: &[]PolicyPackageConfigurationPackage{{
					ID:                              ptr(1),
					Name:                            ptr("test.dmg"),
					Action:                          ptr(PolicyPackageConfigurationPackageActionInstall),
					FillUserTemplates:               ptr(false),
					FillExistingUserHomeDirectories: ptr(false),
				}},
			},
			DistributionPoint: ptr("default"),
		},
		Scripts: &PolicyScripts{
			Size: ptr(1),
			Scripts: &[]PolicyScript{{
				ID:          ptr(1),
				Name:        ptr("Test Script"),
				Priority:    ptr(ScriptPriorityAfter),
				Parameter4:  ptr("Test Parameter"),
				Parameter5:  ptr(""),
				Parameter6:  ptr(""),
				Parameter7:  ptr(""),
				Parameter8:  ptr(""),
				Parameter9:  ptr(""),
				Parameter10: ptr(""),
				Parameter11: ptr(""),
			}},
		},
		Printers: &PolicyPrinters{
			Size:                 ptr(1),
			LeaveExistingDefault: ptr(""),
			Printers: &[]PolicyPrinter{{
				ID:          ptr(1),
				Name:        ptr("Test Printer"),
				Action:      ptr(PolicyPrinterActionInstall),
				MakeDefault: ptr(false),
			}},
		},
		DockItems: &PolicyDockItems{
			Size: ptr(1),
			DockItems: &[]PolicyDockItem{{
				ID:     ptr(1),
				Name:   ptr("Test Dock Item"),
				Action: ptr(PolicyDockItemActionAddToBeginning),
			}},
		},
		AccountMaintenance: &PolicyAccountMaintenance{
			Accounts: &PolicyAccountMaintenanceAccounts{
				Size: ptr(1),
				Accounts: &[]PolicyAccountMaintenanceAccount{{
					Action:           ptr(PolicyAccountMaintenanceAccountActionCreate),
					Username:         ptr("Test Username"),
					Realname:         ptr("Test Realname"),
					Home:             ptr("/Users/test"),
					Picture:          ptr(""),
					Admin:            ptr(false),
					FileVaultEnabled: ptr(false),
				}},
			},
			DirectoryBindings: &PolicyAccountMaintenanceDirectoryBindings{
				Size: ptr(0),
			},
			ManagementAccount: &PolicyAccountMaintenanceManagementAccount{
				Action:                ptr(PolicyAccountMaintenanceManagementAccountActionRandom),
				ManagedPasswordLength: ptr(8),
			},
			OpenFirmwareEFIPassword: &PolicyAccountMaintenanceOpenFirmwareEFIPassword{
				OfMode: ptr(PolicyAccountMaintenanceOpenFirmwareEFIPasswordOfModeCommand),
			},
		},
		Reboot: &PolicyReboot{
			Message:                     ptr("This is a reboot message"),
			StartupDisk:                 ptr("Specify Local Startup Disk"),
			SpecifyStartup:              ptr("/Volumes/test"),
			NoUserLoggedIn:              ptr(PolicyRebootNoUserLoggedInRestartImmediately),
			UserLoggedIn:                ptr(PolicyRebootUserLoggedInRestartIfAPackageOrUpdateRequireIt),
			MinutesUntilReboot:          ptr(5),
			StartRebootTimerImmediately: ptr(true),
			FileVault2Reboot:            ptr(true),
		},
		Maintenance: &PolicyMaintenance{
			Recon:                    ptr(true),
			ResetName:                ptr(false),
			InstallAllCachedPackages: ptr(false),
			Heal:                     ptr(false),
			Prebindings:              ptr(false),
			Permissions:              ptr(false),
			Byhost:                   ptr(false),
			SystemCache:              ptr(false),
			UserCache:                ptr(false),
			Verify:                   ptr(false),
		},
		FilesProcesses: &PolicyFilesProcesses{
			SearchByPath:         ptr(""),
			DeleteFile:           ptr(false),
			LocateFile:           ptr(""),
			UpdateLocateDatabase: ptr(false),
			SpotlightSearch:      ptr(""),
			SearchForProcess:     ptr(""),
			KillProcess:          ptr(false),
			RunCommand:           ptr(""),
		},
		UserInteraction: &PolicyUserInteraction{
			MessageStart:          ptr(""),
			AllowUsersToDefer:     ptr(true),
			AllowDeferralUntilUTC: ptr(time.Date(2025, time.January, 1, 10, 0, 0, 0, time.UTC)),
			AllowDeferralMinutes:  ptr(0),
			MessageFinish:         ptr(""),
		},
		DiskEncryption: &PolicyDiskEncryption{
			Action:                        ptr(PolicyDiskEncryptionActionApply),
			DiskEncryptionConfigurationID: ptr(1),
			AuthRestart:                   ptr(false),
		},
	}
	if !cmp.Equal(policy, want) {
		t.Errorf("Policies.Get()returned %s, want %s", formatWithSpew(policy), formatWithSpew(want))
	}
}

func TestPoliciesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(policiesPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<policies>
  <size>1</size>
  <policy>
    <id>1</id>
    <name>Test Policy</name>
  </policy>
</policies>`))
	})

	ctx := context.Background()
	policies, _, err := client.Policies.List(ctx)
	if err != nil {
		t.Errorf("Policies.List(): %v", err)
	}

	want := &ListPolicies{
		Size: ptr(1),
		Policies: &[]ListPolicy{{
			ID:   ptr(1),
			Name: ptr("Test Policy"),
		}},
	}
	if !cmp.Equal(policies, want) {
		t.Errorf("Policies.List() returned %s, want %s", formatWithSpew(policies), formatWithSpew(want))
	}
}

func TestPoliciesService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	policyID := 1

	mux.HandleFunc(buildHandlePath(policiesPath, "id", fmt.Sprint(policyID)), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, []byte("<policy><general><id>1</id><name>Test Policy Updated</name></general></policy>"))

		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	policy := &Policy{
		General: &PolicyGeneral{
			ID:   ptr(policyID),
			Name: ptr("Test Policy Updated"),
		},
	}
	_, err := client.Policies.Update(ctx, policy)
	if err != nil {
		t.Errorf("Policies.Update(): %v", err)
	}
}
