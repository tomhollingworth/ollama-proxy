package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func main() {
	address := flag.String("address", "", "Proxy address in the format http://host:port?token=... (required)")
	listen := flag.String("listen", ":11434", "Address to listen on (default :11434)")
	flag.Parse()

	if strings.TrimSpace(*address) == "" {
		log.Fatal("-address flag is required")
	}

	addr := strings.TrimSpace(*address)
	parsed, err := url.Parse(addr)
	if err != nil {
		log.Fatalf("Failed to parse address: %v", err)
	}

	token := parsed.Query().Get("token")
	if token == "" {
		log.Fatal("Token not found in address query string")
	}

	// Remove token from query for target host
	q := parsed.Query()
	q.Del("token")
	parsed.RawQuery = q.Encode()
	parsed.Fragment = ""
	target := &url.URL{
		Scheme:   parsed.Scheme,
		Host:     parsed.Host,
		Path:     parsed.Path,
		RawQuery: parsed.RawQuery,
	}

	// Create a cookie jar and custom client
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}
	client := &http.Client{Jar: jar}

	// Initial request to get the cookie
	initURL := target.String() + "?token=" + token
	req, err := http.NewRequest("GET", initURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to do initial request: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("Initial request status:", resp.StatusCode)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxyRequest(w, r, target, client)
	})

	log.Printf("Proxy listening on %s ...", *listen)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func proxyRequest(w http.ResponseWriter, r *http.Request, target *url.URL, client *http.Client) {
	// Create the proxy request
	proxyURL := *target
	proxyURL.Path = r.URL.Path
	proxyURL.RawQuery = r.URL.RawQuery

	fmt.Printf("%s Proxying to %s\n", time.Now().Format(time.RFC3339), proxyURL.String())

	proxyReq, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}
	proxyReq.Header = r.Header.Clone()

	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to reach backend server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
