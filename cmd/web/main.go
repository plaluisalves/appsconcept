package main

import (
	app "appsceoncept/internal"
	"appsceoncept/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main entry point for the application
func main() {

	log.SetFlags(0)
	log.SetOutput(utils.LogWriter{})
	defer utils.RecoverPanic()() // recover from panic and log it

	app := app.New()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":8080",
		Handler: app.Server,
	}

	go func() {

		log.Println("server up and listening...")

		log.Println("server handling requests at " + server.Addr)
		log.Printf("server running at http://localhost%s", server.Addr)

		// start http server
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Println(fmt.Errorf("HTTP server error: %v", err))

		}

		log.Println("stopped serving new connections.")

	}()

	<-stop

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println(fmt.Errorf("http shutdown error: %v", err))
	}
	log.Println("graceful shutdown complete.")

}
