package config

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/ini.v1"
)

const (
	defaultShowNextInterval     = 10
	defaultUpdatePricesInterval = 20
)

type AppConfig struct {
	Fiat                      string
	Ids                       []string
	ShowNextInterval          int64
	UpdatePriceInterval       int64
	SkipCertificateValidation bool
	Version                   string
}

func NewAppConfigFromFile(filepath string) *AppConfig {
	config := &AppConfig{}

	iniConfig, err := ini.Load(filepath)
	if err != nil {
		panic(fmt.Sprintf("Could not load \"%s\".", filepath))
	}

	updatePriceInterval, err := iniConfig.Section("").Key("update_price_interval").Int64()
	if err != nil {
		log.Println(err.Error())
		config.UpdatePriceInterval = defaultUpdatePricesInterval
	} else if updatePriceInterval < defaultUpdatePricesInterval {
		config.UpdatePriceInterval = defaultUpdatePricesInterval
	} else {
		config.UpdatePriceInterval = updatePriceInterval
	}

	showNextInterval, err := iniConfig.Section("").Key("show_next_interval").Int64()
	if err != nil {
		log.Println(err.Error())
		config.ShowNextInterval = defaultShowNextInterval
	} else if showNextInterval < 1 {
		config.ShowNextInterval = defaultShowNextInterval
	} else {
		config.ShowNextInterval = showNextInterval
	}

	ids := iniConfig.Section("").Key("ids").String()
	config.Ids = strings.Fields(ids)
	if len(config.Ids) == 0 {
		panic("No CoinGecko ids set. Add \"ids\" to your \"config.ini\".")
	}

	fiat := iniConfig.Section("").Key("fiat").String()
	if fiat == "" {
		log.Println("No fiat currency set. Defaulting to € (Euro).")
		config.Fiat = "eur"
	} else {
		config.Fiat = fiat
	}

	config.SkipCertificateValidation = false

	return config
}
