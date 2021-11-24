package cmc

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const baseApiUrl = "https://pro-api.coinmarketcap.com/v1"

type CMCMapListing struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type CMCMapResponse struct {
	Data []CMCMapListing `json:"data"`
}

type Client struct {
	apiKey string
}

func InitClient(apiKey string) (client *Client) {
	client = &Client{
		apiKey: apiKey,
	}

	return client
}

func (client *Client) DoRequest(path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", baseApiUrl, path)

	// This is dangerous, but apparently required to make web requests
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := http.Client{
		Timeout:   time.Second * 5,
		Transport: transport,
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accepts", "application/json")
	request.Header.Set("X-CMC_PRO_API_KEY", client.apiKey)

	response, getErr := httpClient.Do(request)
	if getErr != nil {
		return nil, getErr
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}

func (client *Client) GetMap() (CMCMapResponse, error) {
	body, err := client.DoRequest("cryptocurrency/map")
	if err != nil {
		return CMCMapResponse{}, err
	}

	cmcMapRes := CMCMapResponse{}
	jsonErr := json.Unmarshal(body, &cmcMapRes)
	if jsonErr != nil {
		return CMCMapResponse{}, jsonErr
	}
	return cmcMapRes, nil
}
