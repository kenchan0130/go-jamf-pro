package classic

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPackagesService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(packagesPath, "id", "0"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, []byte("<package><name>Test Package</name><filename>Test.pkg</filename></package>"))

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<package>
  <id>1</id>
</package>`))
	})

	ctx := context.Background()
	packageID, _, err := client.Packages.Create(ctx, &Package{
		Name:     ptr("Test Package"),
		Filename: ptr("Test.pkg"),
	})
	if err != nil {
		t.Fatalf("Packages.Create(): %v", err)
	}

	want := ptr(1)
	if !cmp.Equal(packageID, want) {
		t.Errorf("Packages.Create() returned %s, want %s", formatWithSpew(packageID), formatWithSpew(want))
	}
}

func TestPackagesService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(packagesPath, "id", "1"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<package>
  <id>1</id>
</package>`))
	})

	ctx := context.Background()
	packageID := 1
	_, err := client.Packages.Delete(ctx, packageID)
	if err != nil {
		t.Errorf("Packages.Delete(): %v", err)
	}
}

func TestPackagesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	packageID := 1
	mux.HandleFunc(buildHandlePath(packagesPath, "id", fmt.Sprint(packageID)), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<package>
  <id>1</id>
  <name>Test Package</name>
  <category>No category assigned</category>
  <filename>test.pkg</filename>
  <info/>
  <notes/>
  <priority>10</priority>
  <reboot_required>false</reboot_required>
  <fill_user_template>false</fill_user_template>
  <fill_existing_users>false</fill_existing_users>
  <allow_uninstalled>false</allow_uninstalled>
  <os_requirements/>
  <required_processor>None</required_processor>
  <hash_type>MD5</hash_type>
  <hash_value>05133a2170431a3cb50d76607ab0a1cc</hash_value>
  <switch_with_package>Do Not Install</switch_with_package>
  <install_if_reported_available>false</install_if_reported_available>
  <reinstall_option>Do Not Reinstall</reinstall_option>
  <triggering_files/>
  <send_notification>false</send_notification>
</package>`))
	})

	ctx := context.Background()
	pkg, _, err := client.Packages.Get(ctx, packageID)
	if err != nil {
		t.Fatalf("Packages.Get(): %v", err)
	}

	want := &Package{
		ID:                         ptr(1),
		Name:                       ptr("Test Package"),
		Category:                   ptr("No category assigned"),
		Filename:                   ptr("test.pkg"),
		Info:                       ptr(""),
		Notes:                      ptr(""),
		Priority:                   ptr(10),
		RebootRequired:             ptr(false),
		FillUserTemplate:           ptr(false),
		FillExistingUsers:          ptr(false),
		AllowUninstalled:           ptr(false),
		OSRequirements:             ptr(""),
		RequiredProcessor:          ptr("None"),
		HashType:                   ptr("MD5"),
		HashValue:                  ptr("05133a2170431a3cb50d76607ab0a1cc"),
		SwitchWithPackage:          ptr("Do Not Install"),
		InstallIfReportedAvailable: ptr(false),
		ReinstallOption:            ptr("Do Not Reinstall"),
		TriggeringFiles:            ptr(""),
		SendNotification:           ptr(false),
	}
	if !cmp.Equal(pkg, want) {
		t.Errorf("Packages.Get() returned %s, want %s", formatWithSpew(pkg), formatWithSpew(want))
	}
}

func TestPackagesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(packagesPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<packages>
  <size>1</size>
  <package>
    <id>1</id>
    <name>Test Package</name>
  </package>
</packages>`))
	})

	ctx := context.Background()
	packages, _, err := client.Packages.List(ctx)
	if err != nil {
		t.Fatalf("Packages.List(): %v", err)
	}

	want := &ListPackages{
		Size: ptr(1),
		Packages: &[]ListPackage{{
			ID:   ptr(1),
			Name: ptr("Test Package"),
		}},
	}
	if !cmp.Equal(packages, want) {
		t.Errorf("Packages.List() returned %s, want %s", formatWithSpew(packages), formatWithSpew(want))
	}
}

func TestPackagesService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	packageID := 1

	mux.HandleFunc(buildHandlePath(packagesPath, "id", fmt.Sprint(packageID)), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, []byte("<package><id>1</id><name>Test Package Updated</name><filename>test.pkg</filename></package>"))

		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	pkg := &Package{
		ID:       ptr(packageID),
		Name:     ptr("Test Package Updated"),
		Filename: ptr("test.pkg"),
	}
	_, err := client.Packages.Update(ctx, pkg)
	if err != nil {
		t.Errorf("Packages.Update(): %v", err)
	}
}
