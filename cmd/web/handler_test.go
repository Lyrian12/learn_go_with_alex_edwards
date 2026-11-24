package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSnippetCreatePostRejectsInvalidForm(t *testing.T) {
	tests := []struct {
		name   string
		values url.Values
	}{
		{
			name:   "missing title",
			values: url.Values{"content": {"A snippet"}, "expires": {"7"}},
		},
		{
			name:   "missing content",
			values: url.Values{"title": {"A title"}, "expires": {"7"}},
		},
		{
			name:   "blank content",
			values: url.Values{"title": {"A title"}, "content": {"  \n "}, "expires": {"7"}},
		},
		{
			name:   "invalid expiration",
			values: url.Values{"title": {"A title"}, "content": {"A snippet"}, "expires": {"not-a-number"}},
		},
		{
			name:   "unsupported expiration",
			values: url.Values{"title": {"A title"}, "content": {"A snippet"}, "expires": {"30"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(test.values.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			response := httptest.NewRecorder()

			app := &application{}
			app.snippetCreatePost(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
			}
		})
	}
}
