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

func (*CommonAffiliateLink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vs, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ERROR parsing query string: %v", err)
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

// NewAffiliateLink create an HTTP server that handles affiliate links
func NewAffiliateLink() *AffiliateLink {
	return &AffiliateLink{
		maps: map[string]http.Handler{},
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
