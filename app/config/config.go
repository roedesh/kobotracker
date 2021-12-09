package config

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/ini.v1"
)

const (
	defaultShowNextInterval     = 10
	defaultUpdatePricesInterval = 30
)

type AppConfig struct {
	DarkMode                  bool
	Fiat                      string
	Ids                       []string
	DaysChart                 int
	ShowNextInterval          int64
	UpdatePriceInterval       int64
	SkipCertificateValidation bool
	TakeScreenshotOnExit      bool
	Version                   string
}

func NewAppConfigFromFile(filepath string) *AppConfig {
	config := &AppConfig{}

	iniConfig, err := ini.Load(filepath)
	if err != nil {
		panic(fmt.Sprintf("Could not load \"%s\".", filepath))
	}

	ids := iniConfig.Section("").Key("ids").String()
	config.Ids = strings.Fields(ids)
	if len(config.Ids) == 0 {
		panic("No CoinGecko ids set. Add \"ids\" to your \"config.ini\".")
	}

	daysChart, err := iniConfig.Section("").Key("days").Int()
	if err != nil || daysChart < 1 {
		daysChart = 1
	}
	config.DaysChart = daysChart

	darkmode := iniConfig.Section("").Key("darkmode").String()
	config.DarkMode = strings.ToLower(darkmode) == "true"

	takeScreenshot := iniConfig.Section("").Key("take_screenshot").String()
	config.TakeScreenshotOnExit = strings.ToLower(takeScreenshot) == "true"

	fiat := iniConfig.Section("").Key("fiat").String()
	if fiat == "" {
		log.Println("No fiat currency set. Defaulting to € (Euro).")
		config.Fiat = "eur"
	} else {
		config.Fiat = fiat
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

	updatePriceInterval, err := iniConfig.Section("").Key("update_price_interval").Int64()
	if err != nil {
		log.Println(err.Error())
		config.UpdatePriceInterval = defaultUpdatePricesInterval
	} else if updatePriceInterval < defaultUpdatePricesInterval {
		config.UpdatePriceInterval = defaultUpdatePricesInterval
	} else {
		config.UpdatePriceInterval = updatePriceInterval
	}

	config.SkipCertificateValidation = false

	return config
}
