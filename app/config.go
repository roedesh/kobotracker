package app

import (
	"cryptokobo/app/utils"
	"fmt"
	"log"
	"strings"

	"gopkg.in/ini.v1"
)

type AppConfig struct {
	Fiat string
	Ids  []string
}

func getConfigFromIniFile() AppConfig {
	config := AppConfig{}

	iniConfig, err := ini.Load(utils.GetAbsolutePath("config.ini"))
	if err != nil {
		panic(fmt.Sprintf("Could not load \"%s\".", utils.GetAbsolutePath("config.ini")))
	}

	ids := iniConfig.Section("").Key("ids").String()
	config.Ids = strings.Fields(ids)
	if len(config.Ids) == 0 {
		panic("No CoinGecko ids set. Add \"ids\" to your \"config.ini\".")
	}

	fiat := iniConfig.Section("").Key("fiat").String()
	if fiat == "" {
		log.Println("No fiat currency set. Add \"fiat\" to your \"config.ini\".")
		log.Println("Defaulting to € (Euro)")
		config.Fiat = "eur"
	} else {
		config.Fiat = fiat
	}

	return config
}
