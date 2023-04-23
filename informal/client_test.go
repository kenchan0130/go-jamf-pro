package informal

import (
	"bytes"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func setup() (
	client *Client,
	mux *http.ServeMux,
	serverURL string,
	teardown func(),
) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", http.StripPrefix("", mux))

	server := httptest.NewServer(apiHandler)

	jar, _ := cookiejar.New(nil)

	client, _ = NewClient(server.URL, "dummyUsername", "dummyPassword")
	client.HttpClient.Jar = jar

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
