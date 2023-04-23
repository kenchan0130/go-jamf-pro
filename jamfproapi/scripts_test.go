package jamfproapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestScriptsService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(scriptsPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(compactJSON([]byte(`{
			"id": "1",
  			"href": "https://yourJamfProUrl.jamf/api/v1/resource/1"
		}`)))
	})

	ctx := context.Background()
	scriptID, _, err := client.Scripts.Create(ctx, &Script{
		Name: ptr("test"),
	})
	if err != nil {
		t.Fatalf("Scripts.Create(): %v", err)
	}

	if scriptID == nil {
		t.Fatalf("Scripts.Create() returned nil, want non-nil")
	}
}

func TestScriptsService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scriptID := "1"

	mux.HandleFunc(buildHandlePath(scriptsPath, scriptID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Scripts.Delete(ctx, scriptID)
	if err != nil {
		t.Fatalf("Scripts.Delete(): %v", err)
	}
}

func TestScriptsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scriptID := "1"

	mux.HandleFunc(buildHandlePath(scriptsPath, scriptID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compactJSON([]byte(`{
			"id": "1",
			"name": "Test script name",
            "info": "Test script info",
            "notes": "Test script notes",
            "priority": "AFTER",
            "categoryId": "1",
            "categoryName": "Test category name",
            "parameter4": "4",
            "parameter5": "5",
            "parameter6": "6",
            "parameter7": "7",
            "parameter8": "8",
            "parameter9": "9",
            "parameter10": "10",
            "parameter11": "11",
            "osRequirements": "10.10.x",
            "scriptContents": "echo \"test contents.\""
		}`)))
	})

	ctx := context.Background()
	script, _, err := client.Scripts.Get(ctx, scriptID)
	if err != nil {
		t.Fatalf("Scripts.Get() returned error: %v", err)
	}

	want := &Script{
		ID:             ptr("1"),
		Name:           ptr("Test script name"),
		Info:           ptr("Test script info"),
		Notes:          ptr("Test script notes"),
		Priority:       ptr(ScriptPriorityAfter),
		CategoryID:     ptr("1"),
		CategoryName:   ptr("Test category name"),
		Parameter4:     ptr("4"),
		Parameter5:     ptr("5"),
		Parameter6:     ptr("6"),
		Parameter7:     ptr("7"),
		Parameter8:     ptr("8"),
		Parameter9:     ptr("9"),
		Parameter10:    ptr("10"),
		Parameter11:    ptr("11"),
		OSRequirements: ptr("10.10.x"),
		ScriptContents: ptr("echo \"test contents.\""),
	}
	if !cmp.Equal(script, want) {
		t.Fatalf("Scripts.Get() returned %s, want %s", formatWithSpew(script), formatWithSpew(want))
	}
}

func TestScriptsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(scriptsPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compactJSON([]byte(`{
			"totalCount": 1,
			"results": [
				{
			       "id": "1",
			       "name": "Test script name",
                   "info": "Test script info",
                   "notes": "Test script notes",
                   "priority": "AFTER",
                   "categoryId": "1",
                   "categoryName": "Test category name",
                   "parameter4": "4",
                   "parameter5": "5",
                   "parameter6": "6",
                   "parameter7": "7",
                   "parameter8": "8",
                   "parameter9": "9",
                   "parameter10": "10",
                   "parameter11": "11",
                   "osRequirements": "10.10.x",
                   "scriptContents": "echo \"test contents.\""
				}
			]
		}`)))
	})

	ctx := context.Background()
	list, _, err := client.Scripts.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("Scripts.List(): %v", err)
	}

	wantScripts := &[]Script{{
		ID:             ptr("1"),
		Name:           ptr("Test script name"),
		Info:           ptr("Test script info"),
		Notes:          ptr("Test script notes"),
		Priority:       ptr(ScriptPriorityAfter),
		CategoryID:     ptr("1"),
		CategoryName:   ptr("Test category name"),
		Parameter4:     ptr("4"),
		Parameter5:     ptr("5"),
		Parameter6:     ptr("6"),
		Parameter7:     ptr("7"),
		Parameter8:     ptr("8"),
		Parameter9:     ptr("9"),
		Parameter10:    ptr("10"),
		Parameter11:    ptr("11"),
		OSRequirements: ptr("10.10.x"),
		ScriptContents: ptr("echo \"test contents.\""),
	}}
	if !cmp.Equal(list.Scripts, wantScripts) {
		t.Errorf("Scripts.List() returned %s, want %s", formatWithSpew(list.Scripts), formatWithSpew(wantScripts))
	}

	if wantTotalCount := 1; *list.TotalCount != wantTotalCount {
		t.Errorf("Scripts.List() returned %s, want %s", formatWithSpew(list.TotalCount), formatWithSpew(wantTotalCount))
	}
}

func TestScriptsService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scriptID := "1"

	mux.HandleFunc(buildHandlePath(scriptsPath, scriptID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compactJSON([]byte(`{
			"id": "1",
			"name": "Test script name",
            "info": "Test script info updated",
            "notes": "Test script notes",
            "priority": "AFTER",
            "categoryId": "1",
            "categoryName": "Test category name",
            "parameter4": "4",
            "parameter5": "5",
            "parameter6": "6",
            "parameter7": "7",
            "parameter8": "8",
            "parameter9": "9",
            "parameter10": "10",
            "parameter11": "11",
            "osRequirements": "10.10.x",
            "scriptContents": "echo \"test contents.\""
		}`)))
	})

	ctx := context.Background()
	script, _, err := client.Scripts.Update(ctx, &Script{
		ID:   ptr(scriptID),
		Name: ptr("Test script name"),
		Info: ptr("Test script info updated"),
	})
	if err != nil {
		t.Fatalf("Scripts.Update(): %v", err)
	}

	want := &Script{
		ID:             ptr("1"),
		Name:           ptr("Test script name"),
		Info:           ptr("Test script info updated"),
		Notes:          ptr("Test script notes"),
		Priority:       ptr(ScriptPriorityAfter),
		CategoryID:     ptr("1"),
		CategoryName:   ptr("Test category name"),
		Parameter4:     ptr("4"),
		Parameter5:     ptr("5"),
		Parameter6:     ptr("6"),
		Parameter7:     ptr("7"),
		Parameter8:     ptr("8"),
		Parameter9:     ptr("9"),
		Parameter10:    ptr("10"),
		Parameter11:    ptr("11"),
		OSRequirements: ptr("10.10.x"),
		ScriptContents: ptr("echo \"test contents.\""),
	}
	if !cmp.Equal(script, want) {
		t.Fatalf("Scripts.Update() returned %s, want %s", formatWithSpew(script), formatWithSpew(want))
	}
}
