package informal

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"testing"
)

func TestSessionService_Create(t *testing.T) {
	client, mux, serverURL, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, []byte(url.Values{"username": {client.common.username}, "password": {client.common.password}}.Encode()))

		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Errorf("Header.Get(Content-Type) returned %q, want \"application/x-www-form-urlencoded\"", got)
		}

		w.WriteHeader(http.StatusFound)
	})

	ctx := context.Background()
	loginURL := path.Join(serverURL, "/?failover=0123456789ABCDEF")
	_, err := client.Session.Create(
		ctx,
		loginURL,
	)
	if err != nil {
		t.Fatalf("Session.Create(): %v", err)
	}
}
