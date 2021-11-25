package network

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"time"
)

func GetHttpClient() http.Client {
	// This is dangerous, but required to make web requests from the Kobo
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return http.Client{
		Timeout:   time.Second * 5,
		Transport: transport,
	}
}

func DownloadFile(filepath string, url string) error {
	client := GetHttpClient()

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
