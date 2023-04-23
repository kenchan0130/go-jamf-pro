package jamfproapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSSOFailoverService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(ssoFailoverPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compactJSON([]byte(`{
		"failoverUrl": "https://jamf.jamfcloud.com/?failover=0123456789ABCDEF",
		"generationTime": 1674133253000
		}`)))
	})

	ctx := context.Background()
	ssoFailover, _, err := client.SSOFailover.Get(ctx)
	if err != nil {
		t.Fatalf("SSOFailover.Get() returned error: %v", err)
	}

	want := &SSOFailover{
		FailoverURL:    ptr("https://jamf.jamfcloud.com/?failover=0123456789ABCDEF"),
		GenerationTime: ptr(int64(1674133253000)),
	}
	if !cmp.Equal(ssoFailover, want) {
		t.Fatalf("SSOFailover.Get() returned %s, want %s", formatWithSpew(ssoFailover), formatWithSpew(want))
	}
}

func TestSSOFailoverService_Generate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(ssoFailoverPath, "generate"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compactJSON([]byte(`{
		"failoverUrl": "https://jamf.jamfcloud.com/?failover=0123456789ABCDEF",
		"generationTime": 1674133253000
		}`)))
	})

	ctx := context.Background()
	ssoFailover, _, err := client.SSOFailover.Generate(ctx)
	if err != nil {
		t.Fatalf("SSOFailover.Generate(): %v", err)
	}

	want := &SSOFailover{
		FailoverURL:    ptr("https://jamf.jamfcloud.com/?failover=0123456789ABCDEF"),
		GenerationTime: ptr(int64(1674133253000)),
	}
	if !cmp.Equal(ssoFailover, want) {
		t.Fatalf("SSOFailover.Generate() returned %s, want %s", formatWithSpew(ssoFailover), formatWithSpew(want))
	}
}
