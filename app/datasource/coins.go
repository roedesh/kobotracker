package datasource

import (
	"crypto/tls"
	"cryptokobo/app/utils"
	"net/http"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
	coingecko "github.com/superoo7/go-gecko/v3"
)

type Coin struct {
	ID          string
	Name        string
	Symbol      string
	Price       float32
	PricePoints []float64
}

func (coin *Coin) GetBaselinePrices() (float64, float64) {
	var min, max float64

	min, _ = stats.Min(coin.PricePoints)
	max, _ = stats.Max(coin.PricePoints)

	return min, max
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

func NewCoinsDataSource(insecure bool) (cds *CoinsDataSource) {
	cds = &CoinsDataSource{}
	cds.httpClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}},
	}
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
		marketChart, err := cds.client.CoinsIDMarketChart(coin.ID, fiat, "1")
		if err == nil {
			pricePoints := []float64{}
			for _, chartPoint := range *marketChart.Prices {
				pricePoint := float64(chartPoint[1])
				pricePoints = append(pricePoints, pricePoint)
			}
			coin.PricePoints = pricePoints
		}
		updatedCoins = append(updatedCoins, coin)
	}

	cds.Coins = updatedCoins
	return nil
}
