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
	"strconv"
	"time"

	"github.com/kenchan0130/go-jamf-pro/jamf"
	"github.com/kenchan0130/go-jamf-pro/utils"
)

type PoliciesService service

type ListPolicies struct {
	Size     *int          `xml:"size,omitempty"`
	Policies *[]ListPolicy `xml:"policy,omitempty"`
}

type ListPolicy struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type Policy struct {
	General              *PolicyGeneral              `xml:"general,omitempty"`
	Scope                *PolicyScope                `xml:"scope,omitempty"`
	SelfService          *PolicySelfService          `xml:"self_service,omitempty"`
	PackageConfiguration *PolicyPackageConfiguration `xml:"package_configuration,omitempty"`
	Scripts              *PolicyScripts              `xml:"scripts,omitempty"`
	Printers             *PolicyPrinters             `xml:"printers,omitempty"`
	DockItems            *PolicyDockItems            `xml:"dock_items,omitempty"`
	AccountMaintenance   *PolicyAccountMaintenance   `xml:"account_maintenance,omitempty"`
	// The 'reboot' filed is not documented, but it actually exists in the API response.
	Reboot          *PolicyReboot          `xml:"reboot,omitempty"`
	Maintenance     *PolicyMaintenance     `xml:"maintenance,omitempty"`
	FilesProcesses  *PolicyFilesProcesses  `xml:"files_processes,omitempty"`
	UserInteraction *PolicyUserInteraction `xml:"user_interaction,omitempty"`
	DiskEncryption  *PolicyDiskEncryption  `xml:"disk_encryption,omitempty"`
}

type PolicyAccountMaintenance struct {
	Accounts                *PolicyAccountMaintenanceAccounts                `xml:"accounts,omitempty"`
	DirectoryBindings       *PolicyAccountMaintenanceDirectoryBindings       `xml:"directory_bindings,omitempty"`
	ManagementAccount       *PolicyAccountMaintenanceManagementAccount       `xml:"management_account,omitempty"`
	OpenFirmwareEFIPassword *PolicyAccountMaintenanceOpenFirmwareEFIPassword `xml:"open_firmware_efi_password,omitempty"`
}

type PolicyAccountMaintenanceAccount struct {
	Action                 *PolicyAccountMaintenanceAccountAction `xml:"action,omitempty"`
	Username               *string                                `xml:"username,omitempty"`
	Realname               *string                                `xml:"realname,omitempty"`
	Password               *string                                `xml:"password,omitempty"`
	ArchiveHomeDirectory   *bool                                  `xml:"archive_home_directory,omitempty"`
	ArchiveHomeDirectoryTo *string                                `xml:"archive_home_directory_to,omitempty"`
	Home                   *string                                `xml:"home,omitempty"`
	Picture                *string                                `xml:"picture,omitempty"`
	Admin                  *bool                                  `xml:"admin,omitempty"`
	FileVaultEnabled       *bool                                  `xml:"filevault_enabled,omitempty"`
}

type PolicyAccountMaintenanceAccountAction string

const (
	PolicyAccountMaintenanceAccountActionCreate           PolicyAccountMaintenanceAccountAction = "Create"
	PolicyAccountMaintenanceAccountActionReset            PolicyAccountMaintenanceAccountAction = "Reset"
	PolicyAccountMaintenanceAccountActionDelete           PolicyAccountMaintenanceAccountAction = "Delete"
	PolicyAccountMaintenanceAccountActionDisableFileVault PolicyAccountMaintenanceAccountAction = "DisableFileVault"
)

type PolicyAccountMaintenanceAccounts struct {
	Size     *int                               `xml:"size,omitempty"`
	Accounts *[]PolicyAccountMaintenanceAccount `xml:"account,omitempty"`
}

type PolicyAccountMaintenanceDirectoryBindings struct {
	Size     *int                               `xml:"size,omitempty"`
	Bindings *[]PolicyAccountMaintenanceBinding `xml:"binding,omitempty"`
}

type PolicyAccountMaintenanceBinding struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type PolicyAccountMaintenanceManagementAccount struct {
	Action                *PolicyAccountMaintenanceManagementAccountAction `xml:"action,omitempty"`
	ManagedPassword       *string                                          `xml:"managed_password,omitempty"`
	ManagedPasswordLength *int                                             `xml:"managed_password_length,omitempty"`
}

type PolicyAccountMaintenanceManagementAccountAction string

const (
	PolicyAccountMaintenanceManagementAccountActionDoNotChange      PolicyAccountMaintenanceManagementAccountAction = "doNotChange"
	PolicyAccountMaintenanceManagementAccountActionSpecified        PolicyAccountMaintenanceManagementAccountAction = "specified"
	PolicyAccountMaintenanceManagementAccountActionRandom           PolicyAccountMaintenanceManagementAccountAction = "random"
	PolicyAccountMaintenanceManagementAccountActionReset            PolicyAccountMaintenanceManagementAccountAction = "reset"
	PolicyAccountMaintenanceManagementAccountActionResetRandom      PolicyAccountMaintenanceManagementAccountAction = "resetRandom"
	PolicyAccountMaintenanceManagementAccountActionFileVaultEnable  PolicyAccountMaintenanceManagementAccountAction = "fileVaultEnable"
	PolicyAccountMaintenanceManagementAccountActionFileVaultDisable PolicyAccountMaintenanceManagementAccountAction = "fileVaultDisable"
)

type PolicyAccountMaintenanceOpenFirmwareEFIPassword struct {
	OfMode     *PolicyAccountMaintenanceOpenFirmwareEFIPasswordOfMode `xml:"of_mode,omitempty"`
	OfPassword *string                                                `xml:"of_password,omitempty"`
}

type PolicyAccountMaintenanceOpenFirmwareEFIPasswordOfMode string

const (
	PolicyAccountMaintenanceOpenFirmwareEFIPasswordOfModeCommand PolicyAccountMaintenanceOpenFirmwareEFIPasswordOfMode = "command"
	PolicyAccountMaintenanceOpenFirmwareEFIPasswordOfModeNone    PolicyAccountMaintenanceOpenFirmwareEFIPasswordOfMode = "none"
)

type PolicyDiskEncryption struct {
	Action                                 *PolicyDiskEncryptionAction           `xml:"action,omitempty"`
	DiskEncryptionConfigurationID          *int                                  `xml:"disk_encryption_configuration_id,omitempty"`
	AuthRestart                            *bool                                 `xml:"auth_restart,omitempty"`
	RemediateKeyType                       *PolicyDiskEncryptionRemediateKeyType `xml:"remediate_key_type,omitempty"`
	RemediateDiskEncryptionConfigurationID *int                                  `xml:"remediate_disk_encryption_configuration_id,omitempty"`
}

type PolicyDiskEncryptionAction string

const (
	// PolicyDiskEncryptionActionNone The 'none' value is not documented, but it actually exists in the API response.
	PolicyDiskEncryptionActionNone      PolicyDiskEncryptionAction = "none"
	PolicyDiskEncryptionActionApply     PolicyDiskEncryptionAction = "apply"
	PolicyDiskEncryptionActionRemediate PolicyDiskEncryptionAction = "remediate"
)

type PolicyDiskEncryptionRemediateKeyType string

const (
	PolicyDiskEncryptionRemediateKeyTypeIndividual                 PolicyDiskEncryptionRemediateKeyType = "Individual"
	PolicyDiskEncryptionRemediateKeyTypeInstitutional              PolicyDiskEncryptionRemediateKeyType = "Institutional"
	PolicyDiskEncryptionRemediateKeyTypeIndividualAndInstitutional PolicyDiskEncryptionRemediateKeyType = "Individual And Institutional"
)

type PolicyDockItems struct {
	Size      *int              `xml:"size,omitempty"`
	DockItems *[]PolicyDockItem `xml:"dock_item,omitempty"`
}

type PolicyDockItem struct {
	ID     *int                  `xml:"id,omitempty"`
	Name   *string               `xml:"name,omitempty"`
	Action *PolicyDockItemAction `xml:"action,omitempty"`
}

type PolicyDockItemAction string

const (
	PolicyDockItemActionAddToBeginning PolicyDockItemAction = "Add To Beginning"
	PolicyDockItemActionAddToEnd       PolicyDockItemAction = "Add To End"
	PolicyDockItemActionRemove         PolicyDockItemAction = "Remove"
)

type PolicyFilesProcesses struct {
	SearchByPath         *string `xml:"search_by_path,omitempty"`
	DeleteFile           *bool   `xml:"delete_file,omitempty"`
	LocateFile           *string `xml:"locate_file,omitempty"`
	UpdateLocateDatabase *bool   `xml:"update_locate_database,omitempty"`
	SpotlightSearch      *string `xml:"spotlight_search,omitempty"`
	SearchForProcess     *string `xml:"search_for_process,omitempty"`
	KillProcess          *bool   `xml:"kill_process,omitempty"`
	RunCommand           *string `xml:"run_command,omitempty"`
}

type PolicyGeneral struct {
	ID                        *int                  `xml:"id,omitempty"`
	Name                      *string               `xml:"name,omitempty"`
	Enabled                   *bool                 `xml:"enabled,omitempty"`
	Trigger                   *PolicyGeneralTrigger `xml:"trigger,omitempty"`
	TriggerCheckin            *bool                 `xml:"trigger_checkin,omitempty"`
	TriggerEnrollmentComplete *bool                 `xml:"trigger_enrollment_complete,omitempty"`
	TriggerLogin              *bool                 `xml:"trigger_login,omitempty"`
	// The 'trigger_logout' field is documented, but it does not actually exist in the API response.
	//TriggerLogout              *bool                                 `xml:"trigger_logout,omitempty"`
	TriggerNetworkStateChanged *bool                                 `xml:"trigger_network_state_changed,omitempty"`
	TriggerStartup             *bool                                 `xml:"trigger_startup,omitempty"`
	TriggerOther               *string                               `xml:"trigger_other,omitempty"`
	Frequency                  *PolicyGeneralFrequency               `xml:"frequency,omitempty"`
	RetryEvent                 *PolicyGeneralRetryEvent              `xml:"retry_event,omitempty"`
	RetryAttempts              *int                                  `xml:"retry_attempts,omitempty"`
	NotifyOnEachFailedRetry    *bool                                 `xml:"notify_on_each_failed_retry,omitempty"`
	LocationUserOnly           *bool                                 `xml:"location_user_only,omitempty"`
	TargetDrive                *string                               `xml:"target_drive,omitempty"`
	Offline                    *bool                                 `xml:"offline,omitempty"`
	Category                   *GeneralCategory                      `xml:"category,omitempty"`
	DateTimeLimitations        *PolicyGeneralDateTimeLimitations     `xml:"date_time_limitations,omitempty"`
	NetworkLimitations         *PolicyGeneralNetworkLimitations      `xml:"network_limitations,omitempty"`
	OverrideDefaultSettings    *PolicyGeneralOverrideDefaultSettings `xml:"override_default_settings,omitempty"`
	NetworkRequirements        *PolicyGeneralNetworkRequirements     `xml:"network_requirements,omitempty"`
	Site                       *Site                                 `xml:"site,omitempty"`
}

type PolicyGeneralDateTimeLimitations struct {
	ActivationDate      *time.Time                                        `xml:"activation_date,omitempty"`
	ActivationDateEpoch *int64                                            `xml:"activation_date_epoch,omitempty"`
	ActivationDateUTC   *time.Time                                        `xml:"activation_date_utc,omitempty"`
	ExpirationDate      *time.Time                                        `xml:"expiration_date,omitempty"`
	ExpirationDateEpoch *int64                                            `xml:"expiration_date_epoch,omitempty"`
	ExpirationDateUTC   *time.Time                                        `xml:"expiration_date_utc,omitempty"`
	NoExecuteOn         *[]PolicyGeneralDateTimeLimitationsNoExecuteOnDay `xml:"no_execute_on>day,omitempty"`
	NoExecuteStart      *string                                           `xml:"no_execute_start,omitempty"`
	NoExecuteEnd        *string                                           `xml:"no_execute_end,omitempty"`
}

func (pgdtl *PolicyGeneralDateTimeLimitations) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias PolicyGeneralDateTimeLimitations
	aux := struct {
		Alias
		ActivationDate    string `xml:"activation_date"`
		ActivationDateUTC string `xml:"activation_date_utc"`
		ExpirationDate    string `xml:"expiration_date"`
		ExpirationDateUTC string `xml:"expiration_date_utc"`
	}{Alias: Alias(*pgdtl)}

	dateUTCLayout := "2006-01-02T15:04:05.000-0700"

	if pgdtl.ActivationDate != nil {
		aux.ActivationDate = pgdtl.ActivationDate.UTC().Format(time.DateTime)
	}

	if pgdtl.ActivationDateUTC != nil {
		aux.ActivationDateUTC = pgdtl.ActivationDateUTC.UTC().Format(dateUTCLayout)
	}

	if pgdtl.ExpirationDate != nil {
		aux.ExpirationDate = pgdtl.ExpirationDate.UTC().Format(time.DateTime)
	}

	if pgdtl.ExpirationDateUTC != nil {
		aux.ExpirationDateUTC = pgdtl.ExpirationDateUTC.UTC().Format(dateUTCLayout)
	}

	if err := e.EncodeElement(aux, start); err != nil {
		return fmt.Errorf("xml.Encoder#EncodeElement(): %v", err)
	}

	return nil
}

func (pgdtl *PolicyGeneralDateTimeLimitations) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias PolicyGeneralDateTimeLimitations
	aux := struct {
		Alias
		ActivationDate    string `xml:"activation_date"`
		ExpirationDate    string `xml:"expiration_date"`
		ActivationDateUTC string `xml:"activation_date_utc"`
		ExpirationDateUTC string `xml:"expiration_date_utc"`
	}{Alias: Alias(*pgdtl)}

	if err := d.DecodeElement(&aux, &start); err != nil {
		return fmt.Errorf("xml.Decoder#DecodeElement(): %v", err)
	}

	*pgdtl = PolicyGeneralDateTimeLimitations(aux.Alias)

	if aux.ActivationDate != "" {
		parsed, err := time.Parse(time.DateTime, aux.ActivationDate)
		if err != nil {
			return fmt.Errorf("time.Parse(): %v", err)
		}
		pgdtl.ActivationDate = &parsed
	}
	if aux.ExpirationDate != "" {
		parsed, err := time.Parse(time.DateTime, aux.ExpirationDate)
		if err != nil {
			return fmt.Errorf("time.Parse(): %v", err)
		}
		pgdtl.ExpirationDate = &parsed
	}

	dateUTCLayout := "2006-01-02T15:04:05.000-0700"

	if aux.ActivationDateUTC != "" {
		parsed, err := time.Parse(dateUTCLayout, aux.ActivationDateUTC)
		if err != nil {
			return fmt.Errorf("time.Parse(): %v", err)
		}
		if parsed.Location() != time.UTC {
			parsed = parsed.UTC()
		}
		pgdtl.ActivationDateUTC = &parsed
	}
	if aux.ExpirationDateUTC != "" {
		parsed, err := time.Parse(dateUTCLayout, aux.ExpirationDateUTC)
		if err != nil {
			return fmt.Errorf("time.Parse(): %v", err)
		}
		if parsed.Location() != time.UTC {
			parsed = parsed.UTC()
		}
		pgdtl.ExpirationDateUTC = &parsed
	}

	return nil
}

type PolicyGeneralDateTimeLimitationsNoExecuteOnDay string

const (
	PolicyGeneralDateTimeLimitationsNoExecuteOnDaySun PolicyGeneralDateTimeLimitationsNoExecuteOnDay = "Sun"
	PolicyGeneralDateTimeLimitationsNoExecuteOnDayMon PolicyGeneralDateTimeLimitationsNoExecuteOnDay = "Mon"
	PolicyGeneralDateTimeLimitationsNoExecuteOnDayTue PolicyGeneralDateTimeLimitationsNoExecuteOnDay = "Tue"
	PolicyGeneralDateTimeLimitationsNoExecuteOnDayWed PolicyGeneralDateTimeLimitationsNoExecuteOnDay = "Wed"
	PolicyGeneralDateTimeLimitationsNoExecuteOnDayThu PolicyGeneralDateTimeLimitationsNoExecuteOnDay = "Thu"
	PolicyGeneralDateTimeLimitationsNoExecuteOnDayFri PolicyGeneralDateTimeLimitationsNoExecuteOnDay = "Fri"
	PolicyGeneralDateTimeLimitationsNoExecuteOnDaySat PolicyGeneralDateTimeLimitationsNoExecuteOnDay = "Sat"
)

type PolicyGeneralFrequency string

const (
	PolicyGeneralFrequencyOncePerComputer        PolicyGeneralFrequency = "Once per computer"
	PolicyGeneralFrequencyOncePerUserPerComputer PolicyGeneralFrequency = "Once per user per computer"
	PolicyGeneralFrequencyOncePerUser            PolicyGeneralFrequency = "Once per user"
	PolicyGeneralFrequencyOnceEveryDay           PolicyGeneralFrequency = "Once every day"
	PolicyGeneralFrequencyOnceEveryWeek          PolicyGeneralFrequency = "Once every week"
	PolicyGeneralFrequencyOnceEveryMonth         PolicyGeneralFrequency = "Once every month"
	PolicyGeneralFrequencyOngoing                PolicyGeneralFrequency = "Ongoing"
)

type PolicyGeneralOverrideDefaultSettings struct {
	TargetDrive       *string `xml:"target_drive,omitempty"`
	DistributionPoint *string `xml:"distribution_point,omitempty"`
	ForceAfpSmb       *bool   `xml:"force_afp_smb,omitempty"`
	Sus               *string `xml:"sus,omitempty"`
}

type PolicyGeneralNetworkLimitations struct {
	MinimumNetworkConnection *PolicyGeneralNetworkLimitationsMinimumNetworkConnection `xml:"minimum_network_connection,omitempty"`
	AnyIpAddress             *bool                                                    `xml:"any_ip_address,omitempty"`
}

type PolicyGeneralNetworkLimitationsMinimumNetworkConnection string

const (
	PolicyGeneralNetworkLimitationsMinimumNetworkConnectionNoMinimum PolicyGeneralNetworkLimitationsMinimumNetworkConnection = "No Minimum"
	PolicyGeneralNetworkLimitationsMinimumNetworkConnectionEthernet  PolicyGeneralNetworkLimitationsMinimumNetworkConnection = "Ethernet"
)

type PolicyGeneralNetworkRequirements string

const (
	PolicyGeneralNetworkRequirementsAny      PolicyGeneralNetworkRequirements = "Any"
	PolicyGeneralNetworkRequirementsEthernet PolicyGeneralNetworkRequirements = "Ethernet"
)

type PolicyGeneralRetryEvent string

const (
	PolicyGeneralRetryEventNone    PolicyGeneralRetryEvent = "none"
	PolicyGeneralRetryEventTrigger PolicyGeneralRetryEvent = "trigger"
	PolicyGeneralRetryEventCheckin PolicyGeneralRetryEvent = "check-in"
)

type PolicyGeneralTrigger string

const (
	PolicyGeneralTriggerEvent         PolicyGeneralTrigger = "EVENT"
	PolicyGeneralTriggerUserInitiated PolicyGeneralTrigger = "USER_INITIATED"
)

type PolicyMaintenance struct {
	Recon                    *bool `xml:"recon,omitempty"`
	ResetName                *bool `xml:"reset_name,omitempty"`
	InstallAllCachedPackages *bool `xml:"install_all_cached_packages,omitempty"`
	Heal                     *bool `xml:"heal,omitempty"`
	Prebindings              *bool `xml:"prebindings,omitempty"`
	Permissions              *bool `xml:"permissions,omitempty"`
	Byhost                   *bool `xml:"byhost,omitempty"`
	SystemCache              *bool `xml:"system_cache,omitempty"`
	UserCache                *bool `xml:"user_cache,omitempty"`
	Verify                   *bool `xml:"verify,omitempty"`
}

type PolicyPackageConfiguration struct {
	Packages *PolicyPackageConfigurationPackages `xml:"packages,omitempty"`
	// The 'distribution_point' filed is not documented, but it actually exists in the API response.
	DistributionPoint *string `xml:"distribution_point,omitempty"`
}

type PolicyPackageConfigurationPackage struct {
	ID                              *int                                     `xml:"id,omitempty"`
	Name                            *string                                  `xml:"name,omitempty"`
	Action                          *PolicyPackageConfigurationPackageAction `xml:"action,omitempty"`
	FillUserTemplates               *bool                                    `xml:"fut,omitempty"`
	FillExistingUserHomeDirectories *bool                                    `xml:"feu,omitempty"`
	// The 'update_autorun' filed is documented, but it does not actually exist in the API response.
	//UpdateAutorun     *bool                                    `xml:"update_autorun,omitempty"`
}

type PolicyPackageConfigurationPackages struct {
	Size     *int                                 `xml:"size,omitempty"`
	Packages *[]PolicyPackageConfigurationPackage `xml:"package,omitempty"`
}

type PolicyPackageConfigurationPackageAction string

const (
	PolicyPackageConfigurationPackageActionInstall       PolicyPackageConfigurationPackageAction = "Install"
	PolicyPackageConfigurationPackageActionCache         PolicyPackageConfigurationPackageAction = "Cache"
	PolicyPackageConfigurationPackageActionInstallCached PolicyPackageConfigurationPackageAction = "Install Cached"
)

type PolicyPrinter struct {
	ID          *int                 `xml:"id,omitempty"`
	Name        *string              `xml:"name,omitempty"`
	Action      *PolicyPrinterAction `xml:"action,omitempty"`
	MakeDefault *bool                `xml:"make_default,omitempty"`
}

type PolicyPrinters struct {
	Size                 *int             `xml:"size,omitempty"`
	LeaveExistingDefault *string          `xml:"leave_existing_default,omitempty"`
	Printers             *[]PolicyPrinter `xml:"printer,omitempty"`
}

type PolicyPrinterAction string

const (
	PolicyPrinterActionInstall   PolicyPrinterAction = "install"
	PolicyPrinterActionUninstall PolicyPrinterAction = "uninstall"
)

type PolicyReboot struct {
	Message                     *string                     `xml:"message,omitempty"`
	StartupDisk                 *string                     `xml:"startup_disk,omitempty"`
	SpecifyStartup              *string                     `xml:"specify_startup,omitempty"`
	NoUserLoggedIn              *PolicyRebootNoUserLoggedIn `xml:"no_user_logged_in,omitempty"`
	UserLoggedIn                *PolicyRebootUserLoggedIn   `xml:"user_logged_in,omitempty"`
	MinutesUntilReboot          *int                        `xml:"minutes_until_reboot,omitempty"`
	StartRebootTimerImmediately *bool                       `xml:"start_reboot_timer_immediately,omitempty"`
	FileVault2Reboot            *bool                       `xml:"file_vault_2_reboot,omitempty"`
}

type PolicyRebootNoUserLoggedIn string

const (
	PolicyRebootNoUserLoggedInDoNotRestart                       PolicyRebootNoUserLoggedIn = "Do not restart"
	PolicyRebootNoUserLoggedInRestartImmediately                 PolicyRebootNoUserLoggedIn = "Restart immediately"
	PolicyRebootNoUserLoggedInRestartIfAPackageOrUpdateRequireIt PolicyRebootNoUserLoggedIn = "Restart if a package or update requires it"
)

type PolicyRebootUserLoggedIn string

const (
	PolicyRebootUserLoggedInDoNotRestart                       PolicyRebootUserLoggedIn = "Do not restart"
	PolicyRebootUserLoggedInRestart                            PolicyRebootUserLoggedIn = "Restart"
	PolicyRebootUserLoggedInRestartImmediately                 PolicyRebootUserLoggedIn = "Restart immediately"
	PolicyRebootUserLoggedInRestartIfAPackageOrUpdateRequireIt PolicyRebootUserLoggedIn = "Restart if a package or update requires it"
)

type PolicyScope struct {
	AllComputers   *bool                       `xml:"all_computers,omitempty"`
	Computers      *[]PolicyScopeComputer      `xml:"computers>computer,omitempty"`
	ComputerGroups *[]PolicyScopeComputerGroup `xml:"computer_groups>computer_group,omitempty"`
	Buildings      *[]Building                 `xml:"buildings>building,omitempty"`
	Departments    *[]Department               `xml:"departments>department,omitempty"`
	// The 'limit_to_users' will return a value in conjunction with setting the limitations. However, it is not parsed in this struct because it is unwieldy, please see user_groups of limitations.
	//LimitToUsers   interface{}              `xml:"limit_to_users,omitempty"`
	Limitations *PolicyScopeLimitations `xml:"limitations,omitempty"`
	Exclusions  *PolicyScopeExclusions  `xml:"exclusions,omitempty"`
}

type PolicyScopeExclusions struct {
	Computers       *[]PolicyScopeComputer       `xml:"computers>computer,omitempty"`
	ComputerGroups  *[]PolicyScopeComputerGroup  `xml:"computer_groups>computer_group,omitempty"`
	Buildings       *[]Building                  `xml:"buildings>building,omitempty"`
	Departments     *[]Department                `xml:"departments>department,omitempty"`
	Users           *[]PolicyScopeUser           `xml:"users>user,omitempty"`
	UserGroups      *[]PolicyScopeUserGroup      `xml:"user_groups>user_group,omitempty"`
	NetworkSegments *[]PolicyScopeNetworkSegment `xml:"network_segments>network_segment,omitempty"`
	Ibeacons        *[]PolicyScopeIbeacon        `xml:"ibeacons>ibeacon,omitempty"`
}

type PolicyScopeLimitations struct {
	Users           *[]PolicyScopeUser           `xml:"users>user,omitempty"`
	UserGroups      *[]PolicyScopeUserGroup      `xml:"user_groups>user_group,omitempty"`
	NetworkSegments *[]PolicyScopeNetworkSegment `xml:"network_segments>network_segment,omitempty"`
	Ibeacons        *[]PolicyScopeIbeacon        `xml:"ibeacons>ibeacon,omitempty"`
}

type PolicyScopeComputer struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
	UDID *string `xml:"udid,omitempty"`
}

type PolicyScopeComputerGroup struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type PolicyScopeIbeacon struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type PolicyScopeNetworkSegment struct {
	ID   *int    `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type PolicyScopeUser struct {
	// The ID field is documented, but it does not actually exist in the API response.
	//ID *int `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type PolicyScopeUserGroup struct {
	ID   *string `xml:"id,omitempty"`
	Name *string `xml:"name,omitempty"`
}

type PolicyScripts struct {
	Size    *int            `xml:"size,omitempty"`
	Scripts *[]PolicyScript `xml:"script,omitempty"`
}

type PolicyScript struct {
	ID          *int            `xml:"id,omitempty"`
	Name        *string         `xml:"name,omitempty"`
	Priority    *ScriptPriority `xml:"priority,omitempty"`
	Parameter4  *string         `xml:"parameter4,omitempty"`
	Parameter5  *string         `xml:"parameter5,omitempty"`
	Parameter6  *string         `xml:"parameter6,omitempty"`
	Parameter7  *string         `xml:"parameter7,omitempty"`
	Parameter8  *string         `xml:"parameter8,omitempty"`
	Parameter9  *string         `xml:"parameter9,omitempty"`
	Parameter10 *string         `xml:"parameter10,omitempty"`
	Parameter11 *string         `xml:"parameter11,omitempty"`
}

type PolicySelfService struct {
	UseForSelfService           *bool                              `xml:"use_for_self_service,omitempty"`
	SelfServiceDisplayName      *string                            `xml:"self_service_display_name,omitempty"`
	InstallButtonText           *string                            `xml:"install_button_text,omitempty"`
	ReinstallButtonText         *string                            `xml:"reinstall_button_text,omitempty"`
	SelfServiceDescription      *string                            `xml:"self_service_description,omitempty"`
	ForceUsersToViewDescription *bool                              `xml:"force_users_to_view_description,omitempty"`
	SelfServiceIcon             *Icon                              `xml:"self_service_icon,omitempty"`
	FeatureOnMainPage           *bool                              `xml:"feature_on_main_page,omitempty"`
	SelfServiceCategories       *[]SelfServiceCategory             `xml:"self_service_categories>category,omitempty"`
	NotificationEnabled         *bool                              `xml:"-"`
	NotificationType            *PolicySelfServiceNotificationType `xml:"-"`
	NotificationSubject         *string                            `xml:"notification_subject,omitempty"`
	NotificationMessage         *string                            `xml:"notification_message,omitempty"`
}

func (ss *PolicySelfService) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias PolicySelfService
	aux := struct {
		Alias
		Notification *[]string `xml:"notification"`
	}{Alias: Alias(*ss)}

	if ss.NotificationEnabled != nil && ss.NotificationType != nil {
		aux.Notification = &[]string{strconv.FormatBool(*ss.NotificationEnabled), string(*ss.NotificationType)}
	}

	if err := e.EncodeElement(aux, start); err != nil {
		return fmt.Errorf("xml.Encoder#EncodeElement(): %v", err)
	}

	return nil
}

func (ss *PolicySelfService) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias PolicySelfService
	aux := struct {
		Alias
		Notification *[]string `xml:"notification,omitempty"`
	}{Alias: Alias(*ss)}

	if err := d.DecodeElement(&aux, &start); err != nil {
		return fmt.Errorf("xml.Decoder#DecodeElement(): %v", err)
	}

	*ss = PolicySelfService(aux.Alias)

	if aux.Notification != nil && len(*aux.Notification) == 2 {
		notificationEnabled, err := strconv.ParseBool((*aux.Notification)[0])
		if err != nil {
			return nil
		}
		ss.NotificationEnabled = &notificationEnabled

		switch (*aux.Notification)[1] {
		case string(PolicySelfServiceNotificationTypeSelfService):
			v := PolicySelfServiceNotificationTypeSelfService
			ss.NotificationType = &v
		case string(PolicySelfServiceNotificationTypeSelfServiceAndNotificationCenter):
			v := PolicySelfServiceNotificationTypeSelfServiceAndNotificationCenter
			ss.NotificationType = &v
		default:
			return nil
		}
	}

	return nil
}

type PolicySelfServiceNotificationType string

const (
	PolicySelfServiceNotificationTypeSelfService                      PolicySelfServiceNotificationType = "Self Service"
	PolicySelfServiceNotificationTypeSelfServiceAndNotificationCenter PolicySelfServiceNotificationType = "Self Service and Notification Center"
)

type PolicyUserInteraction struct {
	MessageStart          *string    `xml:"message_start,omitempty"`
	AllowUsersToDefer     *bool      `xml:"allow_users_to_defer,omitempty"`
	AllowDeferralUntilUTC *time.Time `xml:"allow_deferral_until_utc,omitempty"`
	AllowDeferralMinutes  *int       `xml:"allow_deferral_minutes,omitempty"`
	MessageFinish         *string    `xml:"message_finish,omitempty"`
}

func (pui *PolicyUserInteraction) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias PolicyUserInteraction
	aux := struct {
		Alias
		AllowDeferralUntilUTC string `xml:"allow_deferral_until_utc"`
	}{Alias: Alias(*pui)}

	dateUTCLayout := "2006-01-02T15:04:05.000-0700"

	if pui.AllowDeferralUntilUTC != nil {
		aux.AllowDeferralUntilUTC = pui.AllowDeferralUntilUTC.UTC().Format(dateUTCLayout)
	}

	if err := e.EncodeElement(aux, start); err != nil {
		return fmt.Errorf("xml.Encoder#EncodeElement(): %v", err)
	}

	return nil
}

func (pui *PolicyUserInteraction) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias PolicyUserInteraction
	aux := struct {
		Alias
		AllowDeferralUntilUTC string `xml:"allow_deferral_until_utc"`
	}{Alias: Alias(*pui)}

	if err := d.DecodeElement(&aux, &start); err != nil {
		return fmt.Errorf("xml.Decoder#DecodeElement(): %v", err)
	}

	*pui = PolicyUserInteraction(aux.Alias)

	if aux.AllowDeferralUntilUTC != "" {
		dateUTCLayout := "2006-01-02T15:04:05.000-0700"
		parsed, err := time.Parse(dateUTCLayout, aux.AllowDeferralUntilUTC)
		if err != nil {
			return fmt.Errorf("time.Parse(): %v", err)
		}
		if parsed.Location() != time.UTC {
			parsed = parsed.UTC()
		}
		pui.AllowDeferralUntilUTC = &parsed
	}

	return nil
}

const policiesPath = "/policies"

func (s *PoliciesService) Create(ctx context.Context, policy *Policy) (*int, *jamf.Response, error) {
	if policy == nil {
		return nil, nil, errors.New("PoliciesService.Create(): cannot create nil policy")
	}
	if policy.General.Name == nil {
		return nil, nil, errors.New("PoliciesService.Create(): cannot create policy with nil Name of General")
	}

	reqBody := &struct {
		*Policy
		XMLName xml.Name `xml:"policy"`
	}{
		Policy: policy,
	}
	body, err := xml.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("xml.Marshal(): %v", err)
	}

	resp, _, err := s.client.Post(ctx, jamf.PostHttpRequestInput{
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: jamf.Uri{
			// When the ID is 0, the ID will be generated automatically at the server side.
			Entity: path.Join(policiesPath, "id", "0"),
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
		XMLName xml.Name `xml:"policy"`
		ID      *int     `xml:"id"`
	}
	if err := xml.Unmarshal(respBody, &data); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	// The create API only returns the ID of the new policy.
	return data.ID, resp, nil
}

func (s *PoliciesService) Delete(ctx context.Context, policyID int) (*jamf.Response, error) {
	resp, _, err := s.client.Delete(ctx, jamf.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(policiesPath, "id", fmt.Sprint(policyID)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("client.Delete(): %v", err)
	}

	return resp, nil
}

func (s *PoliciesService) Get(ctx context.Context, policyID int) (*Policy, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		// The API may return a 404, e.g., when a resource is just created.
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Uri: jamf.Uri{
			Entity: path.Join(policiesPath, "id", fmt.Sprint(policyID)),
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

	var policy Policy
	if err := xml.Unmarshal(respBody, &policy); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	return &policy, resp, nil
}

func (s *PoliciesService) List(ctx context.Context) (*ListPolicies, *jamf.Response, error) {
	resp, _, err := s.client.Get(ctx, jamf.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: jamf.Uri{
			Entity: policiesPath,
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

	var listPolicies ListPolicies
	if err := xml.Unmarshal(respBody, &listPolicies); err != nil {
		return nil, resp, fmt.Errorf("xml.Unmarshal(): %v", err)
	}

	return &listPolicies, resp, nil
}

func (s *PoliciesService) Update(ctx context.Context, policy *Policy) (*jamf.Response, error) {
	if policy == nil {
		return nil, errors.New("PoliciesService.Update(): cannot create nil policy")
	}
	if policy.General.ID == nil {
		return nil, errors.New("PoliciesService.Update(): cannot update policy with nil ID of General")
	}
	if policy.General.Name == nil {
		return nil, errors.New("PoliciesService.Update(): cannot update policy with nil Name of General")
	}

	reqBody := &struct {
		*Policy
		XMLName xml.Name `xml:"policy"`
	}{
		Policy: policy,
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
			Entity: path.Join(policiesPath, "id", fmt.Sprint(*policy.General.ID)),
		},
		Body: bytes.NewBuffer(body),
	})

	if err != nil {
		return nil, fmt.Errorf("client.Put(): %v", err)
	}

	return resp, nil
}
