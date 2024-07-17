package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	tests := map[string]struct {
		apiResponse   string
		expectedError bool
		statusCode    int
	}{
		"successful fetch": {
			apiResponse:   `{"query": {"search": [{"title": "Go"}]}}`,
			expectedError: false,
			statusCode:    http.StatusOK,
		},
		"bad status code": {
			apiResponse:   `{"query": {"search": [{"title": "Go"}]}}`,
			expectedError: true,
			statusCode:    http.StatusInternalServerError,
		},
		"malformed JSON response": {
			apiResponse:   `{"query": {"search": [`,
			expectedError: true,
			statusCode:    http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				w.Write([]byte(test.apiResponse))
			}))
			defer ts.Close()

			client := &Client{}
			var result WikipediaPageQueryJSON
			err := client.fetch(&result, ts.URL)

			if (err != nil) != test.expectedError {
				t.Fatalf("fetch() error = %v, expectedError %v", err, test.expectedError)
			}

			if !test.expectedError && len(result.Query.Search) == 0 {
				t.Fatalf("fetch() result is empty, expected data")
			}
		})
	}
}
