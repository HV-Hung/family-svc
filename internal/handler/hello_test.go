package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HV-Hung/family-svc/internal/handler"
)

func TestHelloHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "default greeting",
			query:          "",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Hello, World!"}`,
		},
		{
			name:           "greeting with name",
			query:          "?name=Alice",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Hello, Alice!"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/hello"+tc.query, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := handler.HelloHandler()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// Encode adds a newline, so we expect a newline at the end of the body
			expected := tc.expectedBody + "\n"
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		})
	}
}
