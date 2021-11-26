package network

import (
	"cryptokobo/app/utils"
	"net/http"
	"strings"

	coingecko "github.com/superoo7/go-gecko/v3"
)

type CoinGeckoClient struct {
	httpClient *http.Client
	client     *coingecko.Client
}

type CoinGeckoCoin struct {
	ID     string
	Name   string
	Symbol string
	Price  float32
}

func InitCoinGecko() (cg *CoinGeckoClient) {
	cg = &CoinGeckoClient{}
	cg.httpClient = GetHttpClient()
	cg.client = coingecko.NewClient(cg.httpClient)

	return cg
}

func getIds(coins []CoinGeckoCoin) []string {
	ids := []string{}
	for _, coin := range coins {
		ids = append(ids, coin.ID)
	}
	return ids
}

func (cg CoinGeckoClient) GetCoinsForIds(ids []string) []CoinGeckoCoin {
	filtered_coins := []CoinGeckoCoin{}
	coins, err := cg.client.CoinsList()
	if err != nil {
		panic(err.Error())
	}

	lowercase_ids := utils.SliceToLowercase(ids)

	for _, coin := range *coins {
		if utils.IsStringInSlice(strings.ToLower(coin.ID), lowercase_ids) {
			filtered_coins = append(filtered_coins, CoinGeckoCoin{ID: coin.ID, Name: coin.Name, Symbol: coin.Symbol})
		}
	}

	if len(filtered_coins) == 0 {
		panic("No coins found! Make sure you set the correct CoinGecko ids.")
	}

	return filtered_coins
}

func (cg CoinGeckoClient) ApplyPricesToCoins(coins []CoinGeckoCoin) ([]CoinGeckoCoin, error) {
	updated_coins := []CoinGeckoCoin{}
	prices, err := cg.client.SimplePrice(getIds(coins), []string{"eur"})
	if err != nil {
		return nil, err
	}

	for _, coin := range coins {
		coin.Price = (*prices)[coin.ID]["eur"]
		updated_coins = append(updated_coins, coin)
	}

	return updated_coins, nil
}
