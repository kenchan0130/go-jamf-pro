package classic

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestComputerGroupsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(computerGroupsPath, "id", "0"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, []byte("<computer_group><name>Test Computer Group</name><is_smart>true</is_smart></computer_group>"))

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<computer_group>
  <id>1</id>
</computer_group>`))
	})

	ctx := context.Background()
	computerGroupID, _, err := client.ComputerGroups.Create(ctx, &ComputerGroup{
		Name:    ptr("Test Computer Group"),
		IsSmart: ptr(true),
	})
	if err != nil {
		t.Fatalf("ComputerGroups.Create(): %v", err)
	}

	want := ptr(1)
	if !cmp.Equal(computerGroupID, want) {
		t.Errorf("ComputerGroups.Create() returned %s, want %s", formatWithSpew(computerGroupID), formatWithSpew(want))
	}
}

func TestComputerGroupsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(computerGroupsPath, "id", "1"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<computer_group>
  <id>1</id>
</computer_group>`))
	})

	ctx := context.Background()
	computerGroupID := 1
	_, err := client.ComputerGroups.Delete(ctx, computerGroupID)
	if err != nil {
		t.Errorf("ComputerGroups.Delete(): %v", err)
	}
}

func TestComputerGroupsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	computerGroupID := 1
	mux.HandleFunc(buildHandlePath(computerGroupsPath, "id", fmt.Sprint(computerGroupID)), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<computer_group>
  <id>1</id>
  <name>Test Computer Group</name>
  <is_smart>true</is_smart>
  <site>
    <id>-1</id>
    <name>None</name>
  </site>
  <criteria>
    <size>1</size>
    <criterion>
      <name>Operating System Version</name>
      <priority>0</priority>
      <and_or>and</and_or>
      <search_type>greater than or equal</search_type>
      <value>12</value>
      <opening_paren>false</opening_paren>
      <closing_paren>false</closing_paren>
    </criterion>
  </criteria>
  <computers>
    <size>1</size>
    <computer>
      <id>1</id>
      <name>Test Computer</name>
      <mac_address>52:42:00:3D:ED:44</mac_address>
      <alt_mac_address>82:33:00:11:67:33</alt_mac_address>
      <serial_number>FFFF123XJK7L</serial_number>
    </computer>
  </computers>
</computer_group>`))
	})

	ctx := context.Background()
	computerGroup, _, err := client.ComputerGroups.Get(ctx, computerGroupID)
	if err != nil {
		t.Fatalf("ComputerGroups.Get(): %v", err)
	}

	want := &ComputerGroup{
		ID:      ptr(1),
		Name:    ptr("Test Computer Group"),
		IsSmart: ptr(true),
		Site: &Site{
			ID:   ptr(-1),
			Name: ptr("None"),
		},
		Criteria: &ComputerGroupCriteria{
			Size: ptr(1),
			Criterion: &[]ComputerGroupCriteriaCriterion{{
				Name:         ptr("Operating System Version"),
				Priority:     ptr(0),
				AndOr:        ptr(ComputerGroupCriteriaCriterionAndOrAnd),
				SearchType:   ptr(ComputerGroupCriteriaCriterionSearchTypeGreaterThanOrEqual),
				Value:        ptr("12"),
				OpeningParen: ptr(false),
				ClosingParen: ptr(false),
			}},
		},
		Computers: &ComputerGroupComputers{
			Size: ptr(1),
			Computers: &[]ComputerGroupComputer{{
				ID:            ptr(1),
				Name:          ptr("Test Computer"),
				MacAddress:    ptr("52:42:00:3D:ED:44"),
				AltMacAddress: ptr("82:33:00:11:67:33"),
				SerialNumber:  ptr("FFFF123XJK7L"),
			}},
		},
	}
	if !cmp.Equal(computerGroup, want) {
		t.Errorf("ComputerGroups.Get()returned %s, want %s", formatWithSpew(computerGroup), formatWithSpew(want))
	}
}

func TestComputerGroupsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(computerGroupsPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<computer_groups>
  <size>1</size>
  <computer_group>
    <id>1</id>
    <name>Test Computer Group</name>
    <is_smart>true</is_smart>
  </computer_group>
</computer_groups>`))
	})

	ctx := context.Background()
	computerGroups, _, err := client.ComputerGroups.List(ctx)
	if err != nil {
		t.Fatalf("ComputerGroups.List(): %v", err)
	}

	want := &ListComputerGroups{
		Size: ptr(1),
		ComputerGroups: &[]ListComputerGroup{{
			ID:      ptr(1),
			Name:    ptr("Test Computer Group"),
			IsSmart: ptr(true),
		}},
	}
	if !cmp.Equal(computerGroups, want) {
		t.Errorf("ComputerGroups.List() returned %s, want %s", formatWithSpew(computerGroups), formatWithSpew(want))
	}
}

func TestComputerGroupsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	computerGroupID := 1

	mux.HandleFunc(buildHandlePath(computerGroupsPath, "id", fmt.Sprint(computerGroupID)), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, []byte("<computer_group><id>1</id><name>Test Computer Group Updated</name><is_smart>true</is_smart></computer_group>"))

		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	computerGroup := &ComputerGroup{
		ID:      ptr(computerGroupID),
		Name:    ptr("Test Computer Group Updated"),
		IsSmart: ptr(true),
	}
	_, err := client.ComputerGroups.Update(ctx, computerGroup)
	if err != nil {
		t.Errorf("ComputerGroups.Update(): %v", err)
	}
}
