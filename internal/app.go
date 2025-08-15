package app

import (
	"appsceoncept/internal/metrics"
	"appsceoncept/utils"
	"net/http"

	"github.com/TickLabVN/tonic"
	"github.com/gin-gonic/gin"
)

// App we are using this struct to the main services
type App struct {
	Server  *gin.Engine
	metrics metrics.Service
}

// New initializes the app
func New() *App {
	app := &App{}

	app.setupMetrics()
	app.setupRoutes()

	return app

}

// setupMetrics initializes the metrics service
func (a *App) setupMetrics() {
	// Initialize metrics
	a.metrics = metrics.New(metrics.WithCustomSolution()) // or metrics.WithPrometheus() for Prometheus metrics
}

// setupRoutes initializes the HTTP routes for the application and defines the openapi spec for the application
func (a *App) setupRoutes() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(utils.Logger())

	// Initialize Tonic with updated project info
	tonic.Init(&tonic.Config{
		OpenAPIVersion: "3.0.0",
		Info: map[string]interface{}{
			"title":       "FizzBuzz API",
			"description": "A simple FizzBuzz API with metrics",
			"version":     "1.0.0",
		},
	})

	// Redirect root to Swagger UI
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger")
	})

	tonic.CreateRoutes("", []tonic.Route{
		{
			Method: tonic.Get,
			Url:    FizzBuzzEndpoint, // "fizzbuzz"
			HandlerRegister: func(path string) {
				r.GET(path, a.HandlerFizzBuzz)
			},
			Schema: &tonic.RouteSchema{
				Summary:     "Generate a FizzBuzz sequence",
				Description: "Returns a sequence where certain numbers are replaced by custom strings. In this example, the replacement string is 'joaquim'.",
				Querystring: FizzBuzzParams{},
				Response: map[int]any{
					200: []string{
						"1", "joaquim", "3", "joaquim", "5", "joaquim", "7", "joaquim", "9", "joaquim",
					},
					400: map[string]string{"error": "invalid parameters"},
				},
			},

			Tags: []string{"FizzBuzz"},
		},
		{
			Method: tonic.Get,
			Url:    MetricsEndpoint, // "metrics"
			HandlerRegister: func(path string) {
				r.GET(path, a.HandlerMetrics)
			},
			Schema: &tonic.RouteSchema{
				Summary: "Incremental metric counter",
				Description: "Returns the total count for a specific combination of query parameters. If the same combination is sent again, the total increments by one. <br>" +
					"For example:<br>" +
					`Request: "int1=&int2=&limit=10&str1=joaquim&str2=ana", Total: 12345` + "<br>" +
					`Request: "int1=&int2=&limit=10&str1=joaquim&str2=joaquim", Total: 1`,
			},
			Tags: []string{"Metrics"},
		},
	})

	// Serve Swagger documentation
	r.GET("/swagger/*w", gin.WrapH(http.StripPrefix("/swagger", tonic.GetHandler())))

	a.Server = r
}
