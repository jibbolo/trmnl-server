package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
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
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {

		mux := &MiddlewareMux{http.NewServeMux()}
		mux.Handle("/static/", http.FileServerFS(staticFiles))
		mux.Handle("/generated/", http.StripPrefix("/generated/", http.FileServerFS(os.DirFS("generated"))))
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

	// 	// Il contenuto HTML da renderizzare
	// 	htmlContent := `
	// 	<!DOCTYPE html>
	// <html lang="en" class="bg-white text-black">
	// <head>
	//   <meta charset="UTF-8" />
	//   <meta name="viewport" content="width=device-width, initial-scale=1" />
	//   <script src="https://cdn.tailwindcss.com"></script>
	//   <style>
	//     body {
	//       font-family: 'Courier New', Courier, monospace;
	//     }
	//   </style>
	// </head>
	// <body class="flex flex-col justify-center items-center h-screen select-none">

	//   <main class="text-center space-y-3">
	//     <h1 class="text-xl leading-tight tracking-wide">DAJE BYOS</h1>
	//     <p class="text-lg">This screen was rendered by BYOS</p>
	//     <a href="#" class="underline">Giacomo Marinangeli</a>
	//   </main>

	//   <footer class="fixed bottom-4 left-4 right-4 max-w-xl mx-auto">
	//     <div class="flex items-center justify-center border border-black rounded-md px-3 py-1 text-xs font-mono leading-none whitespace-nowrap bg-white bg-opacity-90">
	//       trmnl.gmar.dev
	//     </div>
	//   </footer>

	// </body>
	// </html>
	//     `
	// ctx, cancel := chromedp.NewContext(context.Background())
	// defer cancel()

	// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	// defer cancel()

	// var buf []byte

	// // Usa url.PathEscape per codificare correttamente il contenuto HTML
	// dataURL := "data:text/html," + url.PathEscape(htmlContent)

	// err := chromedp.Run(ctx,
	// 	chromedp.Navigate(dataURL),
	// 	chromedp.WaitVisible("body", chromedp.ByQuery),
	// 	chromedp.EmulateViewport(800, 480),
	// 	chromedp.FullScreenshot(&buf, 90),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := os.WriteFile("generated/output.png", buf, 0644); err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Screenshot salvato in output.png")

	cli.Run()
}

// WithLogging wraps http.Handler with logging middleware
// It logs request method, path, remote address, response status and latency in nanoseconds
// See ExampleWithLogging for usage
func StdWithLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body.Close() //  must close
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		slog.Default().Info("Request",
			"method", r.Method,
			"path", r.URL.Path,
			"proto", r.Proto,
			"remoteAddr", r.RemoteAddr,
		)
		handler.ServeHTTP(w, r)
		fmt.Printf("LOG %s %s\n\n", r.URL.Path, bodyBytes)

	})
}

// MiddlewareMux implements multiplexer with middlewares stored in a list
type MiddlewareMux struct {
	*http.ServeMux
}

func (mux *MiddlewareMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.ServeMux.HandleFunc(pattern, handler)
}

// ServeHTTP will serve every request
func (mux *MiddlewareMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close() //  must close
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	headers, _ := json.Marshal(r.Header)
	slog.Default().Info("Request",
		"method", r.Method,
		"path", r.URL.Path,
		"proto", r.Proto,
		"remoteAddr", r.RemoteAddr,
	)
	mux.ServeMux.ServeHTTP(w, r)
	fmt.Printf("LOG %s %s\nBody: %s\n\n", r.URL.Path, headers, bodyBytes)
}
