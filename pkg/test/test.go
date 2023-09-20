package test_utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

const (
	ValidResponse   = "valid response"
	InvalidResponse = "invalid response"
)

func NewServer(status int, write any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprint(write)))
	}))
}

func ReadUnmarshal[T any](resp *http.Response) *T {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var str *T
	if err = json.Unmarshal(body, &str); err != nil {
		return nil
	}

	return str
}
