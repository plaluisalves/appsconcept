package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerFizzBuzz(t *testing.T) {

	a := New()

	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=10&str1=fizz&str2=buzz", nil)
		w := httptest.NewRecorder()

		a.Server.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Contains(t, w.Body.String(), "fizz")
		assert.Equal(t, `["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz"]`, w.Body.String())
	})

	t.Run("error", func(t *testing.T) {
		// missing int1 parameter
		req, _ := http.NewRequest(http.MethodGet, "/fizzbuzz?int2=5&limit=10&str1=fizz&str2=buzz", nil)
		w := httptest.NewRecorder()

		a.Server.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid parameters")
	})

	t.Run("limits", func(t *testing.T) {
		tests := []struct {
			name string
			url  string
		}{
			{"negative_int1", "/fizzbuzz?int1=-1&int2=5&limit=10&str1=fizz&str2=buzz"},
			{"zero_int2", "/fizzbuzz?int1=3&int2=0&limit=10&str1=fizz&str2=buzz"},
			{"limit_too_large", "/fizzbuzz?int1=3&int2=5&limit=100000&str1=fizz&str2=buzz"},
			{"str1_too_long", "/fizzbuzz?int1=3&int2=5&limit=10&str1=" + strings.Repeat("a", 51) + "&str2=buzz"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
				w := httptest.NewRecorder()

				a.Server.ServeHTTP(w, req)

				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Contains(t, w.Body.String(), "invalid parameters")
			})
		}
	})

}
