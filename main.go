package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.Handle("/index.js", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":8081", SetTraceID(mux)))
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func SetTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var traceID apm.TraceID
		if values := r.Header[apmhttp.TraceparentHeader]; len(values) == 1 && values[0] != "" {
			if c, err := apmhttp.ParseTraceparentHeader(values[0]); err == nil {
				traceID = c.Trace
			}
		}
		if err := traceID.Validate(); err != nil {
			uuid := uuid.New()
			var spanID apm.SpanID
			var traceOptions apm.TraceOptions
			copy(traceID[:], uuid[:])
			copy(spanID[:], traceID[8:])
			traceContext := apm.TraceContext{
				Trace:   traceID,
				Span:    spanID,
				Options: traceOptions.WithRecorded(true),
			}
			r.Header.Set(apmhttp.TraceparentHeader, apmhttp.FormatTraceparentHeader(traceContext))
		}

		w.Header().Set("traceID", traceID.String())
		next.ServeHTTP(w, r)
	})
}
