package informal

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
)

// This data is base64 encoded 1x1 black PNG image
const blackPNGImage = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

func TestDistributionFileUploadService_Upload(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	packageID := 1
	fileName := "test.dmg"
	fileType := DistributionFileUploadFileTypePackage
	destination := DistributionFileUploadDestinationDefault
	file, _ := base64.StdEncoding.DecodeString(blackPNGImage)

	mux.HandleFunc(distributionFileUploadPath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, file)

		if got := r.Header.Get("Content-Type"); got != "application/octet-stream" {
			t.Errorf("Header.Get(Content-Type) returned %q, want \"application/octet-stream\"", got)
		}
		if got := r.Header.Get("OBJECT_ID"); got != fmt.Sprint(packageID) {
			t.Errorf("Header.Get(OBJECT_ID) returned %q, want %q", got, packageID)
		}
		if got := r.Header.Get("DESTINATION"); got != string(destination) {
			t.Errorf("Header.Get(DESTINATION) returned %q, want %q", got, destination)
		}
		if got := r.Header.Get("FILE_TYPE"); got != string(fileType) {
			t.Errorf("Header.Get(FILE_TYPE) returned %q, want %q", got, fileType)
		}
		if got := r.Header.Get("FILE_NAME"); got != fileName {
			t.Errorf("Header.Get(FILE_NAME) returned %q, want %q", got, fileName)
		}

		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	_, err := client.DistributionFileUpload.Upload(
		ctx,
		packageID,
		fileName,
		fileType,
		destination,
		bytes.NewBuffer(file),
	)
	if err != nil {
		t.Fatalf("DistributionFileUpload.Upload(): %v", err)
	}
}
