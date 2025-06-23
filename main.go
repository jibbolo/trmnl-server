package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	// "github.com/chromedp/chromedp"
)

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

//go:embed "all:static"
var staticFiles embed.FS

func main() {
	// generateImage()
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new middleware multiplexer
		mux := NewHumaMuxAdapter()

		// Set up static file handlers with content length middleware
		staticHandler := http.FileServer(http.FS(staticFiles))
		mux.Handle("/static/", WithContentLength(staticHandler))

		// Set up generated files handler with content length middleware
		generatedHandler := http.FileServer(http.Dir("generated"))
		mux.Handle("/generated/", http.StripPrefix("/generated/", WithContentLength(generatedHandler)))

		// Create API with the mux (MiddlewareMux implements humago.Mux interface)
		api := humago.New(mux, huma.DefaultConfig("TRMNL API", "1.0"))
		setErrorModel(api)
		addRoutes(api)

		// Chain middlewares from inside out (logging first, then CORS)
		handler := WithCORS(WithLogging(mux))

		// Configure and start the HTTP server
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", options.Port),
			Handler: handler,
		}

		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			server.ListenAndServe()
		})
		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		})
	})

	cli.Run()
}

// This section has been moved to middleware.go
