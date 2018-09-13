package api

import (
	"github.com/codetaming/indy-ingest/internal/persistence"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandlers_Handler(t *testing.T) {
	tests := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{
			in:             httptest.NewRequest("GET", "/", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			logger := log.New(os.Stdout, "ingest-test ", log.LstdFlags|log.Lshortfile)
			dataStore, err := persistence.NewInMemoryDataStore(logger)
			if err != nil {
				logger.Fatalf("failed to create data store: %v", err)
			}

			fileStore, err := persistence.NewLocalFileStore(logger, "/tmp")
			if err != nil {
				logger.Fatalf("failed to create file store: %v", err)
			}
			a := NewAPI(logger, dataStore, fileStore)
			a.CreateDataset(test.out, test.in)
			if test.out.Code != test.expectedStatus {
				t.Logf("expected: %d\ngot: %d\n", test.expectedStatus, test.out.Code)
				t.Fail()
			}
			/*
				body := test.out.Body.String()
				if body != test.expectedBody {
					t.Logf("expected: %s\ngot: %s\n", test.expectedBody, body)
					t.Fail()
				}
			*/
		})
	}
}
