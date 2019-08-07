package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// CommonAffiliateLink redirect to the URL in the query string.
type CommonAffiliateLink struct {
}

func qsParsingError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ERROR parsing query string: %v", err)
}

func (*CommonAffiliateLink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vs, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		qsParsingError(w, err)
		return
	}
	for _, v := range vs {
		if url := v[0]; strings.HasPrefix(url, "http") {
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ERROR: don't know how to handle %v", r)
	return
}

// AffiliateLink is an http.HandleFunc
type AffiliateLink struct {
	maps map[string]http.Handler
}

func debug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ERROR: don't know how to handle %v", r)
	return
}

// OJRQ handles www.ojrq.net
type OJRQ struct{}

func (*OJRQ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vs, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		qsParsingError(w, err)
		return
	}
	vs, err = url.ParseQuery(vs.Get("return"))
	if err != nil {
		qsParsingError(w, err)
		return
	}
	http.Redirect(w, r, vs.Get("u"), http.StatusSeeOther)
}

// NewAffiliateLink create an HTTP server that handles affiliate links
func NewAffiliateLink() *AffiliateLink {
	return &AffiliateLink{
		maps: map[string]http.Handler{
			"www.ojrq.net": &OJRQ{},
		},
	}
}

func (a *AffiliateLink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Header.Get("X-Forwarded-Host")
	handler, ok := a.maps[host]
	if !ok {
		handler = &CommonAffiliateLink{}
	}
	handler.ServeHTTP(w, r)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	log.Fatal(http.ListenAndServe(":"+port, NewAffiliateLink()))
}
