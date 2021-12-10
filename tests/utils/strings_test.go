package main

import (
	"kobotracker/app/utils"
	"testing"
)

func TestGetMoneyString(t *testing.T) {
	highAmountStr := utils.GetMoneyString("eur", 1.20)
	middleAmountStr := utils.GetMoneyString("eur", 0.12)
	lowAmountStr := utils.GetMoneyString("eur", 0.0012)

	if highAmountStr != "€1,20" {
		t.Errorf("highAmountStr was incorrect, got: %s, want: €1,20.", highAmountStr)
	}

	if middleAmountStr != "€0,1200" {
		t.Errorf("middleAmountStr was incorrect, got: %s, want: €0,1200.", middleAmountStr)
	}

	if lowAmountStr != "€0,00120000" {
		t.Errorf("lowAmountStr was incorrect, got: %s, want: €0,00120000.", lowAmountStr)
	}
}

func TestIsStringInSlice(t *testing.T) {
	slice := []string{"apple", "banana", "pineapple"}
	isInSlice := utils.IsStringInSlice("apple", slice)

	if isInSlice == false {
		t.Errorf("isInSlice was incorrect, got: %v, want: true.", isInSlice)
	}

	isInSlice = utils.IsStringInSlice("orange", slice)
	if isInSlice {
		t.Errorf("isInSlice was incorrect, got: %v, want: false.", isInSlice)
	}
}

func TestSliceToLowercase(t *testing.T) {
	slice := []string{"Apple", "BANANA", "pineAPPLE"}
	slice = utils.SliceToLowercase(slice)

	apple := slice[0]
	banana := slice[1]
	pineapple := slice[2]

	if apple != "apple" {
		t.Errorf("apple was incorrect, got: %s, want: apple.", apple)
	}

	if banana != "banana" {
		t.Errorf("banana was incorrect, got: %s, want: banana.", banana)
	}

	if pineapple != "pineapple" {
		t.Errorf("pineapple was incorrect, got: %s, want: pineapple.", apple)
	}
}
