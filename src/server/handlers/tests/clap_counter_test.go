package handlers_test

// import (
// 	"BloTils/src/server/handlers"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// )

// func TestGetClaps(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		method         string
// 		url            string
// 		referer        string
// 		contentType    string
// 		expectedStatus int
// 		expectedCount  int
// 		expectedOK     bool
// 	}{
// 		{
// 			name:           "GET request success",
// 			method:         "GET",
// 			url:            "/claps",
// 			referer:        "http://example.com/page",
// 			contentType:    "application/json",
// 			expectedStatus: http.StatusOK,
// 			expectedCount:  0,
// 			expectedOK:     true,
// 		},
// 		{
// 			name:           "POST request invalid content type",
// 			method:         "POST",
// 			url:            "/claps",
// 			referer:        "http://example.com/page",
// 			contentType:    "text/plain",
// 			expectedStatus: http.StatusUnsupportedMediaType,
// 			expectedCount:  0,
// 			expectedOK:     false,
// 		},
// 		{
// 			name:           "Invalid HTTP method",
// 			method:         "PUT",
// 			url:            "/claps",
// 			referer:        "http://example.com/page",
// 			contentType:    "application/json",
// 			expectedStatus: http.StatusMethodNotAllowed,
// 			expectedCount:  0,
// 			expectedOK:     false,
// 		},
// 		{
// 			name:           "Missing referer",
// 			method:         "GET",
// 			url:            "/claps",
// 			referer:        "",
// 			contentType:    "application/json",
// 			expectedStatus: http.StatusForbidden,
// 			expectedCount:  0,
// 			expectedOK:     false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req := httptest.NewRequest(tt.method, tt.url, nil)
// 			req.Header.Set("Content-Type", tt.contentType)
// 			req.Header.Set("Referer", tt.referer)

// 			w := httptest.NewRecorder()
// 			handlers.GetClaps(w, req)

// 			if w.Code != tt.expectedStatus {
// 				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
// 			}

// 			if w.Code == http.StatusOK {
// 				var response handlers.ClapCounter
// 				err := json.NewDecoder(w.Body).Decode(&response)
// 				if err != nil {
// 					t.Fatalf("failed to decode response: %v", err)
// 				}

// 				if response.Count != tt.expectedCount {
// 					t.Errorf("expected count %d, got %d", tt.expectedCount, response.Count)
// 				}

// 			}
// 		})
// 	}
// }

// func TestGetClapsWithBody(t *testing.T) {
// 	body := strings.NewReader(`{"URL": "http://example.com", "Page": "/test"}`)
// 	req := httptest.NewRequest(http.MethodPost, "/claps", body)
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	handlers.GetClaps(w, req)

// 	if w.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", w.Code)
// 	}

// 	var response handlers.ClapCounter
// 	err := json.NewDecoder(w.Body).Decode(&response)
// 	if err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	if response.URL != "http://example.com" {
// 		t.Errorf("expected URL http://example.com, got %s", response.URL)
// 	}
// }
