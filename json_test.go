package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name   string
		code   int
		msg    string
		err    error
		expect string
	}{
		{
			name:   "400 Bad Request",
			code:   http.StatusBadRequest,
			msg:    "Bad Request",
			err:    nil,
			expect: `{"error":"Bad Request"}`,
		},
		{
			name:   "500 Internal Server Error",
			code:   http.StatusInternalServerError,
			msg:    "Internal Server Error",
			err:    nil,
			expect: `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			respondWithError(rr, tc.code, tc.msg, tc.err)

			if status := rr.Code; status != tc.code {
				t.Errorf("expected status %d, got %d", tc.code, status)
			}

			if strings.TrimSpace(rr.Body.String()) != tc.expect {
				t.Errorf("expected body %s, got %s", tc.expect, rr.Body.String())
			}
		})
	}
}

func TestRespondWithJSON(t *testing.T) {
	type Response struct {
		Message string `json:"message"`
	}

	rr := httptest.NewRecorder()
	expected := Response{Message: "Success"}
	respondWithJSON(rr, http.StatusOK, expected)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	var actual Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if actual.Message != expected.Message {
		t.Errorf("expected message %s, got %s", expected.Message, actual.Message)
	}
}
