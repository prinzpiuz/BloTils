package handlers_test

import (
	"BloTils/src/server"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	test_setup "BloTils/tests"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPing(t *testing.T) {

	tests := []struct {
		name           string
		dbAvailable    bool
		dbPingError    bool
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Successful ping with database",
			dbAvailable:    true,
			dbPingError:    false,
			expectedStatus: http.StatusOK,
			expectedMsg:    "pong",
		},
		{
			name:           "Database not found",
			dbAvailable:    false,
			dbPingError:    false,
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Database Not Found",
		},
		{
			name:           "Database ping error",
			dbAvailable:    true,
			dbPingError:    true,
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Database Ping Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			var mockDB *sql.DB
			var mock sqlmock.Sqlmock
			var err error
			if tt.dbAvailable {
				// Create a mock database
				mockDB, mock, err = sqlmock.New(sqlmock.MonitorPingsOption(true))
				if err != nil {
					t.Fatalf("failed to create mock database: %v", err)
				}

				// Set up expectations
				if tt.dbPingError {
					mock.ExpectPing().WillReturnError(errors.New("Database Ping Error"))
				} else {
					mock.ExpectPing()
				}
				test_setup.App.DBConfig.Connection = mockDB
				ctx := context.WithValue(req.Context(), server.AppContext, *test_setup.App)
				req = req.WithContext(ctx)

			}
			w := httptest.NewRecorder()
			server.Ping(w, req)

			if tt.dbAvailable {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("unfulfilled expectations: %v", err)
				}
			}

			// Check response status
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]string
			err = json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if tt.expectedStatus == http.StatusOK {
				if msg := response["message"]; msg != tt.expectedMsg {
					t.Errorf("expected message %q, got %q", tt.expectedMsg, msg)
				}
			} else {
				// Assuming error messages are in "message" or "error" fields
				if msg := response["message"]; msg != tt.expectedMsg {
					t.Errorf("expected error message %q, got %q", tt.expectedMsg, msg)
				}
			}
		})
	}
}
