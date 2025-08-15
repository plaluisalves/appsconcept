package app

import (
	"appsceoncept/internal/metrics"
	"appsceoncept/tests/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsEndpoint(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		a := New()

		req, _ := http.NewRequest(http.MethodGet, "/"+MetricsEndpoint, nil)
		w := httptest.NewRecorder()
		a.Server.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		a := New()

		// mock metrics service returning an error
		mockService := &mocks.MockMetricsService{
			DataFunc: func() ([]metrics.MetricRow, []metrics.MetricRow, error) {
				return nil, nil, errors.New("something went wrong with prometheus")
			},
		}

		a.metrics = mockService

		req, _ := http.NewRequest(http.MethodGet, "/"+MetricsEndpoint, nil)
		w := httptest.NewRecorder()
		a.Server.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error fetching metrics")
		assert.Contains(t, w.Body.String(), "something went wrong with prometheus")
	})

}
