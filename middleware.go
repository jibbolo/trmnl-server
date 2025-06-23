package main

import (
	"log/slog"
	"maps"
	"net/http"
)

// EnhancedResponseWriter is a consolidated response writer that tracks status code.
// It can be used by multiple middleware components.
type EnhancedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewEnhancedResponseWriter creates a new instance of EnhancedResponseWriter.
func NewEnhancedResponseWriter(w http.ResponseWriter) *EnhancedResponseWriter {
	return &EnhancedResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
	}
}

// WriteHeader captures the status code and passes it to the underlying ResponseWriter
func (erw *EnhancedResponseWriter) WriteHeader(code int) {
	erw.statusCode = code
	erw.ResponseWriter.WriteHeader(code)
}

// StatusCode returns the HTTP status code that was set
func (erw *EnhancedResponseWriter) StatusCode() int {
	return erw.statusCode
}

// ContentLengthResponseWriter wraps an http.ResponseWriter to ensure
// Content-Length header is set properly and prevents chunked encoding.
type ContentLengthResponseWriter struct {
	*EnhancedResponseWriter
	header http.Header
}

// NewContentLengthResponseWriter creates a new wrapper for setting Content-Length.
func NewContentLengthResponseWriter(w http.ResponseWriter) *ContentLengthResponseWriter {
	// Create a copy of the original headers
	headers := make(http.Header)
	maps.Copy(headers, w.Header())
	return &ContentLengthResponseWriter{
		EnhancedResponseWriter: NewEnhancedResponseWriter(w),
		header:                 headers,
	}
}

// Header returns the header map for the writer
func (crw *ContentLengthResponseWriter) Header() http.Header {
	return crw.header
}

// Write passes the data to the underlying ResponseWriter after ensuring headers are set
func (crw *ContentLengthResponseWriter) Write(b []byte) (int, error) {
	// Apply all headers to the original response writer
	for k, vv := range crw.header {
		for _, v := range vv {
			crw.ResponseWriter.Header().Set(k, v)
		}
	}

	// Set the status code if not default
	if crw.EnhancedResponseWriter.StatusCode() != http.StatusOK {
		crw.ResponseWriter.WriteHeader(crw.EnhancedResponseWriter.StatusCode())
	}

	// Write the data to the original response writer
	return crw.ResponseWriter.Write(b)
}

// HumaMuxAdapter implements a multiplexer with enhanced logging that can be used
// with the humago adapter. It wraps a standard http.ServeMux.
type HumaMuxAdapter struct {
	*http.ServeMux
}

// NewHumaMuxAdapter creates a new MiddlewareMux instance
func NewHumaMuxAdapter() *HumaMuxAdapter {
	return &HumaMuxAdapter{
		ServeMux: http.NewServeMux(),
	}
}

// HandleFunc registers a function to handle a pattern
func (mux *HumaMuxAdapter) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.ServeMux.HandleFunc(pattern, handler)
}

// ServeHTTP implements the http.Handler interface
func (mux *HumaMuxAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Simply delegate to the standard ServeMux
	mux.ServeMux.ServeHTTP(w, r)
}

// WithCORS is a middleware that adds CORS headers to allow requests from any origin.
// It handles preflight OPTIONS requests and passes other requests to the next handler.
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

// WithLogging is a middleware that logs request details before and after processing.
// It captures request details including method, path, headers, body, and response status code.
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Create enhanced response writer to track status code
		erw := NewEnhancedResponseWriter(w)

		// Process the request
		next.ServeHTTP(erw, r)

		// Log request and response details
		slog.Default().Info("Request",
			"method", r.Method,
			"status", erw.StatusCode(),
			"path", r.URL.Path,
			"proto", r.Proto,
			"remoteAddr", r.RemoteAddr,
		)
	})
}

// WithContentLength is a middleware that ensures the Content-Length header is set
// for static files. It leverages the file system to get the file size without
// unnecessary buffering, ensuring efficient handling of files from embedded FS.
func WithContentLength(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer wrapper
		crw := NewContentLengthResponseWriter(w)

		// Add a hook to handle filesystem files
		_, isHeadRequest := r.Context().Value(http.LocalAddrContextKey).(bool)

		// If this is a GET or HEAD request (likely for static content)
		if r.Method == "GET" || r.Method == "HEAD" || isHeadRequest {
			// Explicitly remove Transfer-Encoding to prevent chunked encoding
			crw.Header().Del("Transfer-Encoding")
		}

		// Call the next handler with our wrapper
		next.ServeHTTP(crw, r)
	})
}
