package classic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kenchan0130/go-jamf-pro/jamfproapi"
)

func setup() (
	client *Client,
	mux *http.ServeMux,
	serverURL string, //nolint:golint,unparam // Setup for test is used in many tests, so this method returns many values
	teardown func(),
) {
	mux = http.NewServeMux()

	mux.HandleFunc("/api/v1/auth/token", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(compactJSON([]byte(fmt.Sprintf(`{
  			"token": "eyJhbGciOiJIUzUxMiJ9.eyJhdXRoZW50aWNhdGVkLWFwcCI6IkdFTkVSSUMiLCJhdXRoZW50aWNhdGlvbi10eXBlIjoiSlNTIiwiZ3JvdXBzIjpbXSwic3ViamVjdC10eXBlIjoiSlNTX1VTRVJfSUQiLCJ0b2tlbi11dWlkIjoiNzc0YWY3MGYtYWQ0Yy00N2QzLTk2MzktZjEwMjBhMTIwYzExIiwibGRhcC1zZXJ2ZXItaWQiOi0xLCJzdWIiOiIxIiwiZXhwIjoxNTM5NjE5MzQ4fQ.0t7sgYyIyA7kTTmrM8tMGE7fnXcJ1ZzQODAJp0pzg92-cBMQS0Cv8S9oWjkJD7VJS-CHA1dOppr0G_2dCPOfng",
  			"expires": "%s"
		}`, time.Now().Add(time.Minute*30).Format(time.RFC3339)))))
	})

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", http.StripPrefix("", mux))

	server := httptest.NewServer(apiHandler)

	c, _ := jamfproapi.NewClient(server.URL)
	authToken, _, _ := c.APIAuthentication.Token(context.Background(), "dummyUsername", "dummyPassword")

	client, _ = NewClient(server.URL)
	client.AuthorizationToken = authToken.Token

	return client, mux, server.URL, server.Close
}
func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want []byte) {
	t.Helper()
	got, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error reading request body: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

type ptrParam interface{}

func ptr[T ptrParam](v T) *T {
	return &v
}

func formatWithSpew(v interface{}) string {
	return spew.Sprintf("%#v", v)
}

func buildHandlePath(p ...string) string {
	return path.Join(append([]string{apiEndpointPath}, p...)...)
}

func compactJSON(b []byte) []byte {
	buf := &bytes.Buffer{}
	_ = json.Compact(buf, b)
	return buf.Bytes()
}
