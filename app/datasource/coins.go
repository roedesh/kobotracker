package datasource

import (
	"cryptokobo/app/network"
	"cryptokobo/app/utils"
	"net/http"
	"strings"

	coingecko "github.com/superoo7/go-gecko/v3"
)

type Coin struct {
	ID     string
	Name   string
	Symbol string
	Price  float32
}

type CoinsDataSource struct {
	httpClient *http.Client
	client     *coingecko.Client

	Coins []Coin
}

func getIds(coins []Coin) []string {
	ids := []string{}
	for _, coin := range coins {
		ids = append(ids, coin.ID)
	}

	return ids
}

func InitCoinsDataSource() (cds *CoinsDataSource) {
	cds = &CoinsDataSource{}
	cds.httpClient = network.GetHttpClient()
	cds.client = coingecko.NewClient(cds.httpClient)
	cds.Coins = []Coin{}

	return cds
}

func (cds *CoinsDataSource) LoadCoinsForIds(ids []string) {
	filteredCoins := []Coin{}
	coins, err := cds.client.CoinsList()
	if err != nil {
		panic(err.Error())
	}

	lowercaseIds := utils.SliceToLowercase(ids)

	for _, coin := range *coins {
		if utils.IsStringInSlice(strings.ToLower(coin.ID), lowercaseIds) {
			filteredCoins = append(filteredCoins, Coin{ID: coin.ID, Name: coin.Name, Symbol: strings.ToUpper(coin.Symbol)})
		}
	}

	if len(filteredCoins) == 0 {
		panic("No coins found! Make sure you set the correct CoinGecko ids.")
	}

	cds.Coins = filteredCoins
}

func (cds *CoinsDataSource) ApplyPricesToCoins(fiat string) error {
	updatedCoins := []Coin{}
	prices, err := cds.client.SimplePrice(getIds(cds.Coins), []string{fiat})
	if err != nil {
		return err
	}

	for _, coin := range cds.Coins {
		coin.Price = (*prices)[coin.ID][fiat]
		updatedCoins = append(updatedCoins, coin)
	}

	cds.Coins = updatedCoins
	return nil
}
