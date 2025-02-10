package cache

import (
	"sync"
	"testing"

	"github.com/go-countryApi/utils"
	"github.com/stretchr/testify/require"
)

type CountryData struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Currency   string `json:"currency"`
	Population int64  `json:"population"`
}

func randomCountry() CountryData {
	return CountryData{
		Name:       utils.RandomCountryName(),
		Capital:    utils.RandomCountryCapital(),
		Currency:   utils.RandomCountryCurrency(),
		Population: utils.RandomCountryPopulation(),
	}
}

// to test get and set cache methods from in-memv cache
func TestCache(t *testing.T) {
	tCache := NewCache()

	setData1 := randomCountry()
	setData2 := randomCountry()

	tCache.Set("India", setData1)
	tCache.Set("Qatar", setData2)

	setD1, exists := tCache.Get("India")
	require.True(t, exists)
	require.NotEmpty(t, setD1)

	setData1Res := setD1.(CountryData)
	require.Equal(t, setData1.Name, setData1Res.Name)
	require.Equal(t, setData1.Capital, setData1Res.Capital)
	require.Equal(t, setData1.Currency, setData1Res.Currency)
	require.Equal(t, setData1.Population, setData1Res.Population)

	setD2, exists := tCache.Get("Qatar")
	require.True(t, exists)
	require.NotEmpty(t, setD2)

	setData2Res := setD2.(CountryData)
	require.Equal(t, setData2.Name, setData2Res.Name)
	require.Equal(t, setData2.Capital, setData2Res.Capital)
	require.Equal(t, setData2.Currency, setData2Res.Currency)
	require.Equal(t, setData2.Population, setData2Res.Population)
}

// to test parallel cache access
func TestMultipleCache(t *testing.T) {
	countryCh := make(chan CountryData)
	wg := &sync.WaitGroup{}
	tCache := NewCache()
	countries := []string{"India", "Qatar", "South Georgia", "Grenada", "Switzerland"}

	go getDataFromCache(tCache, countryCh, t)

	for i := 0; i < len(countries); i++ {
		wg.Add(1)
		go setDataIntoCache(tCache, countries[i], countryCh, wg)
	}
	wg.Wait()

	close(countryCh)
}

func getDataFromCache(tc *InMemCache, countryCh chan CountryData, t *testing.T) {
	for country := range countryCh {
		data, exists := tc.Get(country.Name)
		require.True(t, exists)
		require.NotEmpty(t, data)

		dataJson := data.(CountryData)

		require.Equal(t, dataJson.Name, country.Name)
		require.Equal(t, dataJson.Capital, country.Capital)
		require.Equal(t, dataJson.Currency, country.Currency)
		require.Equal(t, dataJson.Population, country.Population)
	}
}

func setDataIntoCache(tc *InMemCache, countryName string, countryCh chan CountryData, wg *sync.WaitGroup) {
	defer wg.Done()
	countryD := randomCountry()
	countryD.Name = countryName
	tc.Set(countryD.Name, countryD)
	countryCh <- countryD
}
