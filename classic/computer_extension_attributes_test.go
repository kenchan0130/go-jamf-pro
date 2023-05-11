package classic

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestComputerExtensionAttributesService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(computerExtensionAttributesPath, "id", "0"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, []byte("<computer_extension_attribute><name>Test Computer Extension Attribute</name></computer_extension_attribute>"))

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<computer_extension_attribute>
  <id>1</id>
</computer_extension_attribute>`))
	})

	ctx := context.Background()
	computerExtensionAttributeID, _, err := client.ComputerExtensionAttributes.Create(ctx, &ComputerExtensionAttribute{
		Name: ptr("Test Computer Extension Attribute"),
	})
	if err != nil {
		t.Fatalf("ComputerExtensionAttributes.Create(): %v", err)
	}

	want := ptr(1)
	if !cmp.Equal(computerExtensionAttributeID, want) {
		t.Errorf("ComputerExtensionAttributes.Create() returned %s, want %s", formatWithSpew(computerExtensionAttributeID), formatWithSpew(want))
	}
}

func TestComputerExtensionAttributesService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(computerExtensionAttributesPath, "id", "1"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<computer_extension_attribute>
  <id>1</id>
</computer_extension_attribute>`))
	})

	ctx := context.Background()
	computerExtensionAttributeID := 1
	_, err := client.ComputerExtensionAttributes.Delete(ctx, computerExtensionAttributeID)
	if err != nil {
		t.Errorf("ComputerExtensionAttributes.Delete(): %v", err)
	}
}

func TestComputerExtensionAttributesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	computerExtensionAttributeID := 1
	mux.HandleFunc(buildHandlePath(computerExtensionAttributesPath, "id", fmt.Sprint(computerExtensionAttributeID)), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<computer_extension_attribute>
  <id>1</id>
  <name>Test Computer Extension Attribute</name>
  <enabled>true</enabled>
  <description>Test Description</description>
  <data_type>String</data_type>
  <input_type>
    <type>script</type>
    <platform>Mac</platform>
    <script>#!/bin/bash
&lt;result&gt;test&lt;/result&gt;</script>
  </input_type>
  <inventory_display>General</inventory_display>
</computer_extension_attribute>`))
	})

	ctx := context.Background()
	computerExtensionAttribute, _, err := client.ComputerExtensionAttributes.Get(ctx, computerExtensionAttributeID)
	if err != nil {
		t.Fatalf("ComputerExtensionAttributes.Get(): %v", err)
	}

	want := &ComputerExtensionAttribute{
		ID:          ptr(1),
		Name:        ptr("Test Computer Extension Attribute"),
		Enabled:     ptr(true),
		Description: ptr("Test Description"),
		DataType:    ptr(ComputerExtensionAttributeDataTypeString),
		InputType: &ComputerExtensionAttributeInputType{
			Type:     ptr(ComputerExtensionAttributeInputTypeTypeScript),
			Platform: ptr(ComputerExtensionAttributeInputTypePlatformMac),
			Script:   ptr("#!/bin/bash\n<result>test</result>"),
		},
		InventoryDisplay: ptr(ComputerExtensionAttributeInventoryDisplayGeneral),
	}
	if !cmp.Equal(computerExtensionAttribute, want) {
		t.Errorf("ComputerExtensionAttributes.Get() returned %s, want %s", formatWithSpew(computerExtensionAttribute), formatWithSpew(want))
	}
}

func TestComputerExtensionAttributesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(computerExtensionAttributesPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<computer_extension_attributes>
  <size>1</size>
  <computer_extension_attribute>
    <id>1</id>
    <name>Test Computer Extension Attribute</name>
    <enabled>true</enabled>
  </computer_extension_attribute>
</computer_extension_attributes>`))
	})

	ctx := context.Background()
	computerExtensionAttributes, _, err := client.ComputerExtensionAttributes.List(ctx)
	if err != nil {
		t.Fatalf("ComputerExtensionAttributes.List(): %v", err)
	}

	want := &ListComputerExtensionAttributes{
		Size: ptr(1),
		ComputerExtensionAttributes: &[]ListComputerExtensionAttribute{{
			ID:      ptr(1),
			Name:    ptr("Test Computer Extension Attribute"),
			Enabled: ptr(true),
		}},
	}
	if !cmp.Equal(computerExtensionAttributes, want) {
		t.Errorf("ComputerExtensionAttributes.List() returned %s, want %s", formatWithSpew(computerExtensionAttributes), formatWithSpew(want))
	}
}

func TestComputerExtensionAttributesService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	computerExtensionAttributeID := 1

	mux.HandleFunc(buildHandlePath(computerExtensionAttributesPath, "id", fmt.Sprint(computerExtensionAttributeID)), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testBody(t, r, []byte("<computer_extension_attribute><id>1</id><name>Test Computer Extension Attribute Updated</name><enabled>true</enabled></computer_extension_attribute>"))

		w.WriteHeader(http.StatusCreated)
	})

	ctx := context.Background()
	computerExtensionAttribute := &ComputerExtensionAttribute{
		ID:      ptr(computerExtensionAttributeID),
		Name:    ptr("Test Computer Extension Attribute Updated"),
		Enabled: ptr(true),
	}
	_, err := client.ComputerExtensionAttributes.Update(ctx, computerExtensionAttribute)
	if err != nil {
		t.Errorf("ComputerExtensionAttributes.Update(): %v", err)
	}
}
