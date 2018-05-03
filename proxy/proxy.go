package proxy

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/drmjo/heimdall/spec"
)

type Config struct {
	NoGzip bool
	Addr   string
}

var Transport http.RoundTripper = &http.Transport{
	Proxy: nil,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	DisableCompression:    true}

func Serve(api *spec.Api, config *Config) {
	router := NewRouter(api, config.NoGzip)
	srv := &http.Server{
		Handler: router,
		Addr:    config.Addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on %s", config.Addr)
	log.Fatal(srv.ListenAndServe())
}
