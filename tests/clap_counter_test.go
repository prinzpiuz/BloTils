package handlers_test

import (
	"BloTils/src/server"
	test_setup "BloTils/tests"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

const testDomain = "example.com"

// Helper function to create mock DB with cleanup function
func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	cleanup := func() {
		_ = mockDB.Close()
	}
	return mockDB, mock, cleanup
}

// Helper function to setup common domain query expectation
func expectDomainQuery(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`SELECT d\.id, d\.settings_id, d\.domain, d\.created_time, ds\.id, ds\.likes, ds\.comments, ds\.created_time FROM Domain d JOIN DomainSettings ds ON d\.id = ds\.id WHERE d\.domain = \?`).
		WithArgs(testDomain).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "settings_id", "domain", "created_time",
				"id", "likes", "comments", "created_time",
			}).AddRow(
				1, 1, testDomain,
				time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
				1, 1, 0,
				time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			))
}

// Helper for GetLikedIP query
func expectGetLikedIPQuery(mock sqlmock.Sqlmock, ip, path string, found bool) {
	query := mock.ExpectQuery(`SELECT \* FROM Liked_IPs\s+WHERE Liked_IPs\.ip = \?\s+AND Liked_IPs\.domain = \?\s+AND Liked_IPs\.path = \?`).
		WithArgs(ip, testDomain, path)

	if found {
		query.WillReturnRows(
			sqlmock.NewRows([]string{"id", "ip", "count", "domain", "path", "created_time"}).
				AddRow(1, ip, 1, testDomain, path, time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)),
		)
	} else {
		query.WillReturnError(sql.ErrNoRows)
	}
}

// Helper for UpdateLikeCount query
func expectUpdateLikeQuery(mock sqlmock.Sqlmock, uri string, domainID int) {
	mock.ExpectExec(`INSERT INTO Likes \(uri, domain_id, count, created_time\)\s+VALUES \(\?, \?, 1, datetime\(\)\)\s+ON CONFLICT\(uri, domain_id\)\s+DO UPDATE SET count = count \+ 1`).
		WithArgs(uri, domainID).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

// Helper for UpdateIPLikeCount query
func expectUpdateIPLikeQuery(mock sqlmock.Sqlmock, ip, path string) {
	mock.ExpectExec(`INSERT INTO Liked_IPs \(ip, domain, path, count, created_time\)\s+VALUES \(\?, \?, \?, 1, datetime\(\)\)\s+ON CONFLICT\(ip, domain, path\)\s+DO UPDATE SET count = count \+ 1`).
		WithArgs(ip, testDomain, path).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

// Helper for GetLikes query
func expectGetLikesQuery(mock sqlmock.Sqlmock, uri string, count int) {
	mock.ExpectQuery(`SELECT id, uri, count, domain_id FROM Likes WHERE uri = \? AND domain_id = \(SELECT id FROM Domain WHERE domain = \?\)`).
		WithArgs(uri, testDomain).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "uri", "count", "domain_id"}).
				AddRow(1, uri, count, 1),
		)
}

// Helper for GetLikes returning no rows
func expectGetLikesQueryNoRows(mock sqlmock.Sqlmock, uri string) {
	mock.ExpectQuery(`SELECT id, uri, count, domain_id FROM Likes WHERE uri = \? AND domain_id = \(SELECT id FROM Domain WHERE domain = \?\)`).
		WithArgs(uri, testDomain).
		WillReturnError(sql.ErrNoRows)
}

func TestGetClaps(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		referer        string
		contentType    string
		expectedStatus int
		expectedCount  int
		expectedOK     bool
		setupMock      func(sqlmock.Sqlmock)
	}{
		{
			name:           "GET request success",
			method:         "GET",
			url:            "/api/v1/count_like?page=/blog/post",
			referer:        "https://example.com/",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
			expectedOK:     true,
			setupMock: func(mock sqlmock.Sqlmock) {
				expectDomainQuery(mock)
				expectGetLikesQueryNoRows(mock, "/blog/post")
			},
		},
		{
			name:           "POST request invalid content type",
			method:         "POST",
			url:            "/api/v1/count_like",
			referer:        "https://example.com/",
			contentType:    "text/plain",
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedCount:  0,
			expectedOK:     false,
			setupMock: func(mock sqlmock.Sqlmock) {
				// No DB calls expected
			},
		},
		{
			name:           "Missing referer",
			method:         "GET",
			url:            "/api/v1/count_like",
			referer:        "",
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
			expectedOK:     false,
			setupMock: func(mock sqlmock.Sqlmock) {
				// No DB calls expected
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, cleanup := newMockDB(t)
			defer cleanup()

			test_setup.App.DBConfig.Connection = mockDB

			req := httptest.NewRequest(tt.method, tt.url, nil)
			ctx := context.WithValue(req.Context(), server.AppContext, *test_setup.App)
			req = req.WithContext(ctx)
			req.Header.Set("Content-Type", tt.contentType)
			req.Header.Set("Referer", tt.referer)

			tt.setupMock(mock)

			w := httptest.NewRecorder()
			server.GetClaps(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Code == http.StatusOK {
				var response server.ClapResponse
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if response.Count != tt.expectedCount {
					t.Errorf("expected count %d, got %d", tt.expectedCount, response.Count)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetClapsWithBody(t *testing.T) {
	mockDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	test_setup.App.DBConfig.Connection = mockDB

	body := strings.NewReader(`{"page": "/blog/post"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/count_like", body)
	ctx := context.WithValue(req.Context(), server.AppContext, *test_setup.App)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://example.com/")

	expectDomainQuery(mock)
	expectGetLikedIPQuery(mock, "192.0.2.1", "/blog/post", false)
	expectUpdateLikeQuery(mock, "/blog/post", 1)
	expectUpdateIPLikeQuery(mock, "192.0.2.1", "/blog/post")
	expectGetLikesQuery(mock, "/blog/post", 1)

	w := httptest.NewRecorder()
	server.GetClaps(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response server.ClapResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Domain != testDomain {
		t.Errorf("expected %s, got %s", testDomain, response.Domain)
	}

	if !response.Success {
		t.Errorf("Request not successful")
	}

	if response.Message != "Clap Counted Successfully" {
		t.Errorf("Message is not expected, got: %s", response.Message)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestGetClapsWithBodyForAlreadyLiked(t *testing.T) {
	mockDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	test_setup.App.DBConfig.Connection = mockDB

	body := strings.NewReader(`{"page": "/blog/post"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/count_like", body)
	ctx := context.WithValue(req.Context(), server.AppContext, *test_setup.App)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://example.com/")

	expectDomainQuery(mock)
	expectGetLikedIPQuery(mock, "192.0.2.1", "/blog/post", true)
	expectGetLikesQuery(mock, "/blog/post", 5)

	w := httptest.NewRecorder()
	server.GetClaps(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response server.ClapResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Domain != testDomain {
		t.Errorf("expected %s, got %s", testDomain, response.Domain)
	}

	if response.Success {
		t.Errorf("Like should not be counted for already liked IP")
	}

	if response.Message != "Clap Already Counted" {
		t.Errorf("Message is not expected, got: %s", response.Message)
	}

	if response.Count != 5 {
		t.Errorf("expected count 5, got %d", response.Count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestGetClapsNewLikeNoExistingRecord(t *testing.T) {
	mockDB, mock, cleanup := newMockDB(t)
	defer cleanup()

	test_setup.App.DBConfig.Connection = mockDB

	body := strings.NewReader(`{"page": "/blog/new-post"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/count_like", body)
	ctx := context.WithValue(req.Context(), server.AppContext, *test_setup.App)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://example.com/")

	expectDomainQuery(mock)
	expectGetLikedIPQuery(mock, "192.0.2.1", "/blog/new-post", false)
	expectUpdateLikeQuery(mock, "/blog/new-post", 1)
	expectUpdateIPLikeQuery(mock, "192.0.2.1", "/blog/new-post")
	expectGetLikesQuery(mock, "/blog/new-post", 1)

	w := httptest.NewRecorder()
	server.GetClaps(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response server.ClapResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Errorf("Request should be successful")
	}

	if response.Count != 1 {
		t.Errorf("expected count 1, got %d", response.Count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
