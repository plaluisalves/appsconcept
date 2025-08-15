package app

import (
	"appsceoncept/tests/fake"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	vegeta "github.com/tsenart/vegeta/lib"
)

// TestApp will test the performance of the app, specifically the FizzBuzz endpoint. At the same time we can also check the metrics endpoint at http://localhost:8080/metrics
func TestApp(t *testing.T) {

	t.Skip("skipping performance test")

	requests := 100
	rate := vegeta.Rate{Freq: requests, Per: time.Second} // Rate of requests per second
	duration := 5 * time.Second                           // duration of the test

	// initialize the app
	app := New()

	// test server
	var ts *httptest.Server
	ready := make(chan string)

	go func() {
		// create a listener with the desired port. This way we can force a port for the test server.
		l, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatal(err)
		}
		ts = httptest.NewUnstartedServer(http.HandlerFunc(app.Server.ServeHTTP))
		ts.Listener.Close()
		ts.Listener = l

		ts.Start()
		ready <- ts.URL

	}()

	<-ready

	base, err := url.Parse(ts.URL + "/fizzbuzz")
	assert.NoError(t, err, "failed to parse base URL")

	t.Run("fizzbuzz", func(t *testing.T) {
		f := func(tgt *vegeta.Target) error {
			q := url.Values{}
			q.Set("int1", fmt.Sprint(fake.Int1()))
			q.Set("int2", fmt.Sprint(fake.Int2()))
			q.Set("limit", fmt.Sprint(fake.Limit()))
			q.Set("str1", fake.Str1())
			q.Set("str2", fake.Str2())

			base.RawQuery = q.Encode()

			tgt.Method = "GET"
			tgt.URL = base.String()
			return nil
		}
		attacker := vegeta.NewAttacker()
		_ = attacker.Attack(f, rate, duration, "fizzbuzz")
	})
	select {}

}
