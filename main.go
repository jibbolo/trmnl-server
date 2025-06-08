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
)

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

//go:embed "static"
var staticFiles embed.FS

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {

		mux := http.NewServeMux()
		mux.Handle("/static/", http.FileServerFS(staticFiles))
		api := humago.New(mux, huma.DefaultConfig("TRMNL API", "1.0"))

		setErrorModel(api)
		addRoutes(api)

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", options.Port),
			Handler: mux,
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
