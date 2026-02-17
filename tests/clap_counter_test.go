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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockDB *sql.DB
			var mock sqlmock.Sqlmock
			var err error
			var req *http.Request
			mockDB, mock, err = sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock database in test %s : %v", t.Name(), err)
			}
			test_setup.App.DBConfig.Connection = mockDB
			req = httptest.NewRequest(tt.method, tt.url, nil)
			ctx := context.WithValue(req.Context(), server.AppContext, *test_setup.App)
			req = req.WithContext(ctx)
			req.Header.Set("Content-Type", tt.contentType)
			req.Header.Set("Referer", tt.referer)
			const expectedQuery = `SELECT d.id, d.settings_id, d.domain, d.created_time, ds.id, ds.likes, ds.comments, ds.created_time FROM Domain d JOIN DomainSettings ds ON d.id = ds.id WHERE d.domain = ?`
			mock.ExpectQuery(expectedQuery).
				WithArgs("example.com").
				WillReturnRows(
					sqlmock.NewRows([]string{
						"id",
						"settings_id",
						"domain",
						"created_time",
						"id",
						"likes",
						"comments",
						"created_time",
					}).AddRow(
						1,
						1,
						"example.com",
						time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
						1,
						1,
						0,
						time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)))
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
		})
	}
}

func TestGetClapsWithBody(t *testing.T) {
	var mockDB *sql.DB
	var mock sqlmock.Sqlmock
	var err error
	mockDB, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database in test %s : %v", t.Name(), err)
	}
	test_setup.App.DBConfig.Connection = mockDB
	body := strings.NewReader(`{"page": "/blog/post"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/count_like", body)
	ctx := context.WithValue(req.Context(), server.AppContext, *test_setup.App)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://example.com/")

	// 1. First expectation: Domain query (checkForDomain)
	expectedQuery := `SELECT d.id, d.settings_id, d.domain, d.created_time, ds.id, ds.likes, ds.comments, ds.created_time FROM Domain d JOIN DomainSettings ds ON d.id = ds.id WHERE d.domain = ?`
	mock.ExpectQuery(expectedQuery).
		WithArgs("example.com").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"settings_id",
				"domain",
				"created_time",
				"id",
				"likes",
				"comments",
				"created_time",
			}).AddRow(
				1,
				1,
				"example.com",
				time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
				1,
				1,
				0,
				time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)))

	// 2. Second expectation: Check if IP already liked (hasAlreadyLikedIP -> GetLikedIP)
	// Note: extractClientIP strips the port, so "192.0.2.1:1234" becomes "192.0.2.1"
	mock.ExpectQuery(
		`^SELECT \* FROM Liked_IPs\s+WHERE Liked_IPs\.ip = \?\s+AND Liked_IPs\.domain = \?\s+AND Liked_IPs\.path = \?$`,
	).WithArgs("192.0.2.1", "example.com", "/blog/post").WillReturnError(sql.ErrNoRows)

	// 3. Third expectation: Update like count (UpdateLikeCount)
	mock.ExpectExec(
		`^INSERT INTO Likes\(uri, count, domain_id\)\s+`+
			`VALUES \(\?, 1, \?\)\s+`+
			`ON CONFLICT\(uri\)\s+`+
			`DO UPDATE\s+`+
			`SET count = count \+ 1$`,
	).WithArgs("/blog/post", 1). // uri, domain_id
					WillReturnResult(sqlmock.NewResult(1, 1))

	// 4. Fourth expectation: Record IP like (UpdateIPLikeCount)
	mock.ExpectExec(
		`^INSERT INTO Liked_IPs\(domain, path, ip, count, created_time\)\s+`+
			`VALUES\(\?, \?, \?, 1, datetime\(\)\)\s+`+
			`ON CONFLICT\(ip\)\s+`+
			`DO UPDATE\s+`+
			`SET count = count \+ 1$`,
	).WithArgs("example.com", "/blog/post", "192.0.2.1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 5. Fifth expectation: Get updated likes count (GetLikes)
	mockRows := sqlmock.NewRows([]string{
		"Likes.id",
		"Likes.uri",
		"Likes.domain_id",
		"Likes.count",
		"Domain.id",
		"Domain.settings_id",
		"Domain.domain",
		"Domain.created_time",
	}).AddRow(
		1,             // Likes.id
		"/blog/post",  // Likes.uri
		1,             // Likes.domain_id
		1,             // Likes.count (incremented)
		1,             // Domain.id
		1,             // Domain.settings_id
		"example.com", // Domain.domain
		time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), // Domain.created_time
	)
	mock.ExpectQuery(
		`^SELECT\s+Likes\.id,\s+Likes\.uri,\s+Likes\.domain_id,\s+Likes\.count,\s+Domain\.id,\s+Domain\.settings_id,\s+Domain\.domain,\s+Domain\.created_time\s+`+
			`FROM\s+Likes\s+JOIN Domain ON Likes\.domain_id = Domain\.id\s+`+
			`WHERE\s+Domain\.domain = \?\s+AND Likes\.uri = \?$`,
	).WithArgs("example.com", "/blog/post").
		WillReturnRows(mockRows)

	w := httptest.NewRecorder()
	server.GetClaps(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response server.ClapResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.URL != "example.com" {
		t.Errorf("expected example.com, got %s", response.URL)
	}

	if !response.Success {
		t.Errorf("Request Not Succesfull")
	}

	if response.Message != "Clap Counted Successfully" {
		t.Errorf("Message Is Not Expected")
	}
}

func TestGetClapsWithBodyForAlreadyLiked(t *testing.T) {
	var mockDB *sql.DB
	var mock sqlmock.Sqlmock
	var err error
	mockDB, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database in test %s : %v", t.Name(), err)
	}
	test_setup.App.DBConfig.Connection = mockDB
	body := strings.NewReader(`{"page": "/blog/post"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/count_like", body)
	ctx := context.WithValue(req.Context(), server.AppContext, *test_setup.App)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://example.com/")

	// 1. First expectation: Domain query (checkForDomain)
	expectedQuery := `SELECT d.id, d.settings_id, d.domain, d.created_time, ds.id, ds.likes, ds.comments, ds.created_time FROM Domain d JOIN DomainSettings ds ON d.id = ds.id WHERE d.domain = ?`
	mock.ExpectQuery(expectedQuery).
		WithArgs("example.com").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"settings_id",
				"domain",
				"created_time",
				"id",
				"likes",
				"comments",
				"created_time",
			}).AddRow(
				1,
				1,
				"example.com",
				time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
				1,
				1,
				0,
				time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)))

	// 2. Second expectation: Check if IP already liked (hasAlreadyLikedIP -> GetLikedIP)
	// Returns a row indicating IP has already liked this page
	// Note: extractClientIP strips the port, so "192.0.2.1:1234" becomes "192.0.2.1"
	likedIPscolumns := []string{
		"id",
		"ip",         // Liked_IPs.ip
		"count",      // Liked_IPs.count
		"domain",     // Liked_IPs.domain
		"path",       // Liked_IPs.path
		"created_at", // Liked_IPs.created_at (timestamp)
	}
	likedipMockRows := sqlmock.NewRows(likedIPscolumns).
		AddRow(
			1,
			"192.0.2.1",
			1,
			"example.com",
			"/blog/post",
			time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		)
	mock.ExpectQuery(
		`^SELECT \* FROM Liked_IPs\s+WHERE Liked_IPs\.ip = \?\s+AND Liked_IPs\.domain = \?\s+AND Liked_IPs\.path = \?$`,
	).WithArgs("192.0.2.1", "example.com", "/blog/post").WillReturnRows(likedipMockRows)

	// 3. Third expectation: Get current likes count (GetLikes) - called when already liked
	mockRows := sqlmock.NewRows([]string{
		"Likes.id",
		"Likes.uri",
		"Likes.domain_id",
		"Likes.count",
		"Domain.id",
		"Domain.settings_id",
		"Domain.domain",
		"Domain.created_time",
	}).AddRow(
		1,             // Likes.id
		"/blog/post",  // Likes.uri
		1,             // Likes.domain_id
		5,             // Likes.count (existing count)
		1,             // Domain.id
		1,             // Domain.settings_id
		"example.com", // Domain.domain
		time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), // Domain.created_time
	)
	mock.ExpectQuery(
		`^SELECT\s+Likes\.id,\s+Likes\.uri,\s+Likes\.domain_id,\s+Likes\.count,\s+Domain\.id,\s+Domain\.settings_id,\s+Domain\.domain,\s+Domain\.created_time\s+`+
			`FROM\s+Likes\s+JOIN Domain ON Likes\.domain_id = Domain\.id\s+`+
			`WHERE\s+Domain\.domain = \?\s+AND Likes\.uri = \?$`,
	).WithArgs("example.com", "/blog/post").
		WillReturnRows(mockRows)

	w := httptest.NewRecorder()
	server.GetClaps(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response server.ClapResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.URL != "example.com" {
		t.Errorf("expected example.com, got %s", response.URL)
	}

	if response.Success {
		t.Errorf("Like Counted For Already Liked IP")
	}

	if response.Message != "Clap Already Counted" {
		t.Errorf("Message Is Not Expected")
	}
}
