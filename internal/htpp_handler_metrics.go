package app

import (
	"appsceoncept/internal/metrics"
	"embed"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

//go:embed metrics
var metricsfs embed.FS

// HandlerMetrics handles the metrics endpoint, fetching and displaying metrics data.
func (a *App) HandlerMetrics(c *gin.Context) {
	// Fetch metrics data
	rows, maxs, err := a.metrics.Data()
	if err != nil {
		c.String(http.StatusInternalServerError, "error fetching metrics: %v", err)
		return
	}

	// Parse the template directly from embed.FS
	tmpl, err := template.ParseFS(metricsfs, "metrics/template.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "error parsing template: %v", err)
		return
	}

	// Render template with data
	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")

	data := struct {
		Rows []metrics.MetricRow
		Maxs []metrics.MetricRow
	}{
		Rows: rows,
		Maxs: maxs,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "error executing template: %v", err)
		return
	}
}
