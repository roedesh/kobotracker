package network

import (
	"crypto/tls"
	"net/http"
	"time"
)

func GetHttpClient() *http.Client {
	// In order to make web requests on the Kobo, we need to disable certificate validation.
	// Normally you shouldn't do this, but it should be fine in this case, because the
	// app only displays some JSON data.
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{
		Timeout:   time.Second * 5,
		Transport: transport,
	}
}
