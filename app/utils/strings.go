package utils

import "strings"

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
