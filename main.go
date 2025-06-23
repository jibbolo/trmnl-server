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

		// Apply ContentLengthMiddleware to static files
		staticHandler := http.FileServerFS(staticFiles)
		mux.Handle("/static/", ContentLengthMiddleware(staticHandler))

		generatedHandler := http.FileServerFS(os.DirFS("generated"))
		mux.Handle("/generated/", http.StripPrefix("/generated/", ContentLengthMiddleware(generatedHandler)))

		api := humago.New(mux, huma.DefaultConfig("TRMNL API", "1.0"))
		setErrorModel(api)
		addRoutes(api)

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", options.Port),
			Handler: WithCORS(mux),
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

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// WithCORS is a middleware that adds CORS headers to all responses
func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
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
	lw := NewLoggingResponseWriter(w)
	mux.ServeMux.ServeHTTP(lw, r)
	slog.Default().Info("Request",
		"method", r.Method,
		"status", lw.statusCode,
		"path", r.URL.Path,
		"proto", r.Proto,
		"remoteAddr", r.RemoteAddr,
	)
	fmt.Printf("LOG %s %d %s\nBody: %s\n\n", r.URL.Path, lw.statusCode, headers, bodyBytes)
}

// ContentLengthMiddleware ensures Content-Length header is set and disables chunked encoding
func ContentLengthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply to static files, particularly images
		if strings.HasPrefix(r.URL.Path, "/static/") &&
			(strings.HasSuffix(r.URL.Path, ".png") ||
				strings.HasSuffix(r.URL.Path, ".bmp") ||
				strings.HasSuffix(r.URL.Path, ".jpg")) {

			// Create a custom ResponseWriter that captures the response
			crw := &customResponseWriter{
				ResponseWriter: w,
				buffer:         &bytes.Buffer{},
				header:         make(http.Header),
			}

			// Call the next handler with our custom writer
			next.ServeHTTP(crw, r)

			// Set Content-Length header based on buffer size
			crw.header.Set("Content-Length", fmt.Sprintf("%d", crw.buffer.Len()))

			// Explicitly remove any Transfer-Encoding header to prevent chunked encoding
			crw.header.Del("Transfer-Encoding")

			// Copy all headers from our custom ResponseWriter to the original
			for k, vv := range crw.header {
				for _, v := range vv {
					w.Header().Add(k, v)
				}
			}

			// Set status code
			if crw.status > 0 {
				w.WriteHeader(crw.status)
			}

			// Write the buffered response body
			w.Write(crw.buffer.Bytes())
		} else {
			// For non-image requests, pass through normally
			next.ServeHTTP(w, r)
		}
	})
}

// customResponseWriter captures the response to calculate its size
type customResponseWriter struct {
	http.ResponseWriter
	buffer *bytes.Buffer
	header http.Header
	status int
}

func (crw *customResponseWriter) Header() http.Header {
	return crw.header
}

func (crw *customResponseWriter) Write(b []byte) (int, error) {
	return crw.buffer.Write(b)
}

func (crw *customResponseWriter) WriteHeader(statusCode int) {
	crw.status = statusCode
}
