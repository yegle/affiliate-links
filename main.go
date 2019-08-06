package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// CommonAffiliateLink redirect to the URL in the query string.
type CommonAffiliateLink struct {
	S string
}

func (c *CommonAffiliateLink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vs, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ERROR parsing query string: %v", err)
		return
	}
	url := vs.Get(c.S)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// AffiliateLink is an http.HandleFunc
type AffiliateLink struct {
	maps map[string]http.Handler
}

// NewAffiliateLink create an HTTP server that handles affiliate links
func NewAffiliateLink() *AffiliateLink {
	return &AffiliateLink{
		maps: map[string]http.Handler{
			"go.redirectingat.com":  &CommonAffiliateLink{S: "url"},
			"click.linksynergy.com": &CommonAffiliateLink{S: "murl"},
			"www.jdoqocy.com":       &CommonAffiliateLink{S: "url"},
		},
	}
}

func (a *AffiliateLink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Header.Get("X-Forwarded-Host")
	handler, ok := a.maps[host]
	if !ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ERROR: don't know how to handle %q: %v", host, r)
		return
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
