package jamfproapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// This data is base64 encoded 1x1 black PNG image
const blackPNGImage = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

func TestIconService_Download(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	iconID := "1"

	want, _ := base64.StdEncoding.DecodeString(blackPNGImage)

	mux.HandleFunc(buildHandlePath(iconPath, "download", iconID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(want)
	})

	ctx := context.Background()
	resp, err := client.Icon.Download(ctx, iconID, IconContentOptions{})
	if err != nil {
		t.Fatalf("Icon.Download(): %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("Error closing response body: %v", err)
		}
	}()

	iconBytes, _ := io.ReadAll(resp.Body)

	if !cmp.Equal(iconBytes, want) {
		t.Fatalf("Icon.Download() returned %s, want %s", formatWithSpew(iconBytes), formatWithSpew(want))
	}
}

func TestIconService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	iconID := "1"

	mux.HandleFunc(buildHandlePath(iconPath, iconID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compactJSON([]byte(`{
 			"name": "test.png",
			"url": "https://stage-ics.services.jamfcloud.com/icon/hash_c315ef577b84505de1bfcb50b0c4b1c963da30b2a805f84b24ad09f282b7fad4",
			"id": 1
		}`)))
	})

	ctx := context.Background()
	icon, _, err := client.Icon.Get(ctx, iconID)
	if err != nil {
		t.Fatalf("Icon.Get(): %v", err)
	}

	iconIDInt, _ := strconv.Atoi(iconID)
	want := &Icon{
		ID:   ptr(iconIDInt),
		Name: ptr("test.png"),
		URL:  ptr("https://stage-ics.services.jamfcloud.com/icon/hash_c315ef577b84505de1bfcb50b0c4b1c963da30b2a805f84b24ad09f282b7fad4"),
	}
	if !cmp.Equal(icon, want) {
		t.Fatalf("Icon.Get() returned %s, want %s", formatWithSpew(icon), formatWithSpew(want))
	}
}

func TestIconService_Upload(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(buildHandlePath(iconPath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		if got := r.Header.Get("Content-Type"); !strings.Contains(got, "multipart/form-data") {
			t.Errorf("Header.Get(Content-Type) returned %q, want included multipart/form-data", got)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(compactJSON([]byte(`{
 			"name": "test.png",
			"url": "https://stage-ics.services.jamfcloud.com/icon/hash_c315ef577b84505de1bfcb50b0c4b1c963da30b2a805f84b24ad09f282b7fad4",
			"id": 1
		}`)))
	})

	file, _ := base64.StdEncoding.DecodeString(blackPNGImage)

	ctx := context.Background()
	icon, _, err := client.Icon.Upload(ctx, "test.png", bytes.NewBuffer(file))
	if err != nil {
		t.Fatalf("Icon.Upload(): %v", err)
	}

	want := &Icon{
		ID:   ptr(1),
		Name: ptr("test.png"),
		URL:  ptr("https://stage-ics.services.jamfcloud.com/icon/hash_c315ef577b84505de1bfcb50b0c4b1c963da30b2a805f84b24ad09f282b7fad4"),
	}
	if !cmp.Equal(icon, want) {
		t.Fatalf("Icon.Upload() returned %s, want %s", formatWithSpew(icon), formatWithSpew(want))
	}
}
