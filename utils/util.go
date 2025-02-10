package utils

import (
	"math/rand"

	"github.com/go-countryApi/models"
)

func RandomCountryData() models.ApiRes {
	countryNames := []string{"India", "Qatar", "South Georgia", "Grenada", "Switzerland"}
	countryCapital := []string{"New Delhi", "Doha", "King Edward Point", "St. George's", "Bern"}
	countryCurrency := []string{"₹", "ر.ق", "£", "$", "Fr."}
	countryPopulation := []int64{1380004385, 2881060, 30, 112519, 8654622}

	randNo := rand.Intn(len(countryNames))

	return models.ApiRes{
		Status: true,
		Result: models.Result{
			Name:       countryNames[randNo],
			Capital:    countryCapital[randNo],
			Currency:   countryCurrency[randNo],
			Population: countryPopulation[randNo],
		},
	}
}

func RandomCountryName() string {
	countryNames := []string{"India", "Qatar", "South Georgia", "Grenada", "Switzerland"}
	n := len(countryNames)

	return countryNames[rand.Intn(n)]
}

func RandomCountryCapital() string {
	countryCapital := []string{"New Delhi", "King Edward Point", "St. George's", "Bern"}
	n := len(countryCapital)

	return countryCapital[rand.Intn(n)]
}

func RandomCountryCurrency() string {
	countryCurrency := []string{"₹", "£", "$", "Fr."}
	n := len(countryCurrency)

	return countryCurrency[rand.Intn(n)]
}

func RandomCountryPopulation() int64 {
	countryPopulation := []int64{1380004385, 1380004, 104385, 827322}
	n := len(countryPopulation)

	return countryPopulation[rand.Intn(n)]
}
