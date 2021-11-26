package utils

import "strings"

func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SliceToLowercase(str_slice []string) []string {
	lowercase_strings := []string{}
	for _, str := range str_slice {
		lowercase_strings = append(lowercase_strings, strings.ToLower(str))
	}

	return lowercase_strings
}
