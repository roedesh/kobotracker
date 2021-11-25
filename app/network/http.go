package network

import (
	"crypto/tls"
	"net/http"
	"time"
)

func GetHttpClient() http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return http.Client{
		Timeout:   time.Second * 5,
		Transport: transport,
	}
}
