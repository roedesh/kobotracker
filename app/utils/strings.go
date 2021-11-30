package utils

import (
	"strings"

	"github.com/leekchan/accounting"
)

func GetMoneyString(currency string, number float64) string {
	localeInfo := accounting.LocaleInfo[strings.ToUpper(currency)]
	acc := accounting.Accounting{Symbol: localeInfo.ComSymbol, Precision: 2, Thousand: localeInfo.ThouSep, Decimal: localeInfo.DecSep}

	if number < 0.01 {
		acc.Precision = 8
	} else if number < 1 {
		acc.Precision = 4
	}

	return acc.FormatMoney(number)
}

func IsStringInSlice(str string, list []string) bool {
	for _, strFromList := range list {
		if strFromList == str {
			return true
		}
	}
	return false
}

func SliceToLowercase(str_slice []string) []string {
	lowercaseStrings := []string{}
	for _, str := range str_slice {
		lowercaseStrings = append(lowercaseStrings, strings.ToLower(str))
	}

	return lowercaseStrings
}
