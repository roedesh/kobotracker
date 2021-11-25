package cmc

import (
	"cryptokobo/app/network"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const baseApiUrl = "https://pro-api.coinmarketcap.com/v1"

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

	httpClient := network.GetHttpClient()

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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getIdParameter(cryptocurrencies []CMCCryptocurrency) string {
	ids := []string{}
	for _, cryptocurrency := range cryptocurrencies {
		ids = append(ids, strconv.Itoa(cryptocurrency.ID))
	}
	return strings.Join(ids, ",")
}

func (client *Client) GetIds(symbols []string) ([]CMCCryptocurrency, error) {
	body, err := client.DoRequest("cryptocurrency/map")
	if err != nil {
		return nil, err
	}

	cmcMapRes := CMCResponse{}
	jsonErr := json.Unmarshal(body, &cmcMapRes)
	if jsonErr != nil {
		return nil, jsonErr
	}

	lowercase_symbols := []string{}
	for _, symbol := range symbols {
		lowercase_symbols = append(lowercase_symbols, strings.ToLower(symbol))
	}

	cryptocurrencies := []CMCCryptocurrency{}
	for _, crypto := range cmcMapRes.Data {
		if stringInSlice(strings.ToLower(crypto.Symbol), lowercase_symbols) {
			cryptocurrencies = append(cryptocurrencies, crypto)
		}
	}

	return cryptocurrencies, nil
}

func (client *Client) GetPrices(cryptocurrencies []CMCCryptocurrency) ([]CMCCryptocurrency, error) {
	body, err := client.DoRequest(fmt.Sprintf("cryptocurrency/quotes/latest?id=%s&convert=EUR", getIdParameter(cryptocurrencies)))
	if err != nil {
		return nil, err
	}

	cmcMapRes := CMCMapResponse{}
	jsonErr := json.Unmarshal(body, &cmcMapRes)
	if jsonErr != nil {
		return nil, jsonErr
	}

	updated_cryptocurrencies := []CMCCryptocurrency{}
	for _, cryptocurrency := range cryptocurrencies {
		cryptocurrency.Quote = cmcMapRes.Data[strconv.Itoa(cryptocurrency.ID)].Quote
		updated_cryptocurrencies = append(updated_cryptocurrencies, cryptocurrency)
	}

	return updated_cryptocurrencies, nil
}

func (client *Client) GetLogos(cryptocurrencies []CMCCryptocurrency) ([]CMCCryptocurrency, error) {
	body, err := client.DoRequest(fmt.Sprintf("cryptocurrency/info?id=%s", getIdParameter(cryptocurrencies)))
	if err != nil {
		return nil, err
	}

	cmcMapRes := CMCMapResponse{}
	jsonErr := json.Unmarshal(body, &cmcMapRes)
	if jsonErr != nil {
		return nil, jsonErr
	}

	updated_cryptocurrencies := []CMCCryptocurrency{}
	for _, cryptocurrency := range cryptocurrencies {
		cryptocurrency.Logo = cmcMapRes.Data[strconv.Itoa(cryptocurrency.ID)].Logo
		updated_cryptocurrencies = append(updated_cryptocurrencies, cryptocurrency)
	}

	return updated_cryptocurrencies, nil
}
