package jamfproapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCategoriesService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(categoriesPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(compactJSON([]byte(`{
			"id": "1",
  			"href": "https://yourJamfProUrl.jamf/api/v1/resource/1"
		}`)))
	})

	ctx := context.Background()
	categoryID, _, err := client.Categories.Create(ctx, &Category{
		Name:     ptr("test"),
		Priority: ptr(int32(9)),
	})
	if err != nil {
		t.Fatalf("Categories.Create(): %v", err)
	}

	if categoryID == nil {
		t.Fatalf("Categories.Create() returned nil, want non-nil")
	}
}

func TestCategoriesService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	categoryID := "1"

	mux.HandleFunc(buildHandlePath(categoriesPath, categoryID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Categories.Delete(ctx, categoryID)
	if err != nil {
		t.Fatalf("Categories.Delete(): %v", err)
	}
}

func TestCategoriesService_DeleteMultiple(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(categoriesPath, "delete-multiple"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, []byte(`{"ids":["1","2"]}`))

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Categories.DeleteMultiple(ctx, []string{"1", "2"})
	if err != nil {
		t.Fatalf("Categories.DeleteMultiple(): %v", err)
	}
}

func TestCategoriesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	categoryID := "1"

	mux.HandleFunc(buildHandlePath(categoriesPath, categoryID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compactJSON([]byte(`{
  			"id": "1",
  			"name": "test",
  			"priority": 9
		}`)))
	})

	ctx := context.Background()
	category, _, err := client.Categories.Get(ctx, categoryID)
	if err != nil {
		t.Fatalf("Categories.Get(): %v", err)
	}

	want := &Category{
		ID:       ptr("1"),
		Name:     ptr("test"),
		Priority: ptr(int32(9)),
	}
	if !cmp.Equal(category, want) {
		t.Fatalf("Categories.Get() returned %s, want %s", formatWithSpew(category), formatWithSpew(want))
	}
}

func TestCategoriesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(categoriesPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compactJSON([]byte(`{
  			"totalCount": 1,
  			"results": [
    			{
      				"id": "1",
      				"name": "test",
      				"priority": 9
    			}
  			]
		}`)))
	})

	ctx := context.Background()
	list, _, err := client.Categories.List(ctx, ListOptions{})
	if err != nil {
		t.Fatalf("Categories.List(): %v", err)
	}

	wantCategories := &[]Category{
		{
			ID:       ptr("1"),
			Name:     ptr("test"),
			Priority: ptr(int32(9)),
		},
	}
	if !cmp.Equal(list.Categories, wantCategories) {
		t.Errorf("Categories.List() returned %s, want %s", formatWithSpew(list.Categories), formatWithSpew(wantCategories))
	}

	if wantTotalCount := 1; *list.TotalCount != wantTotalCount {
		t.Errorf("Categories.List() returned %s, want %s", formatWithSpew(list.TotalCount), formatWithSpew(wantTotalCount))
	}
}

func TestCategoriesService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	categoryID := "1"

	mux.HandleFunc(buildHandlePath(categoriesPath, categoryID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compactJSON([]byte(`{
      		"id": "1",
      		"name": "test updated",
      		"priority": 9
		}`)))
	})

	ctx := context.Background()
	category, _, err := client.Categories.Update(ctx, &Category{
		ID:       ptr(categoryID),
		Name:     ptr("test updated"),
		Priority: ptr(int32(9)),
	})
	if err != nil {
		t.Fatalf("Categories.Update(): %v", err)
	}

	want := &Category{
		ID:       ptr("1"),
		Name:     ptr("test updated"),
		Priority: ptr(int32(9)),
	}
	if !cmp.Equal(category, want) {
		t.Fatalf("Categories.Update() returned %s, want %s", formatWithSpew(category), formatWithSpew(want))
	}
}
