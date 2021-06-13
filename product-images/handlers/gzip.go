package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

type WrappedResponseWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(w)
	return &WrappedResponseWriter{w, gw}
}

func (ww *WrappedResponseWriter) Header() http.Header {
	return ww.w.Header()
}

func (ww *WrappedResponseWriter) Write(d []byte) (int, error) {
	return ww.gw.Write(d)
}

func (ww *WrappedResponseWriter) WriteHeader(statuscode int) {
	ww.w.WriteHeader(statuscode)
}

func (ww *WrappedResponseWriter) Flush() {
	ww.gw.Flush()
	ww.gw.Close()
}

func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// create a gziped response
			wrw := NewWrappedResponseWriter(rw)
			wrw.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(wrw, r)
			defer wrw.Flush()
			return
		}

		// handle normal
		next.ServeHTTP(rw, r)
	})
}
