package authorize

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func endpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestMissingHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()

	http.HandlerFunc(Require(endpoint)).ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status %v; got %v", http.StatusUnauthorized, rec.Result().Status)
	}
}

func TestValidHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.kK9JnTXZwzNo3BYNXJT57PGLnQk-Xyu7IBhRWFmc4C0")

	rec := httptest.NewRecorder()

	http.HandlerFunc(Require(endpoint)).ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(rec.Result().Body)

		t.Logf("Response body: %v", string(b))
		t.Errorf("Expected status %v; got %v", http.StatusOK, rec.Result().Status)
	}
}
