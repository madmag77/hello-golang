package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_given_Rest_when_Put_KV_then_can_read_value(t *testing.T) {
	t.Run("put and get same value", func(t *testing.T) {
		// Given
		testValue := "testValue"

		// When
		request, _ := http.NewRequest(http.MethodPut, "/v1/key/123", strings.NewReader(testValue))
		response := httptest.NewRecorder()

		PutHandler(response, request)

		request, _ = http.NewRequest(http.MethodGet, "/v1/key/123", nil)
		response = httptest.NewRecorder()

		GetHandler(response, request)

		val := response.Body.String()
		want := testValue

		// Then
		if val != want {
			t.Errorf("got %q, want %q", val, want)
		}
	})
}
