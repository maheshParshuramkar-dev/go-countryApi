package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-countryApi/cache"
	"github.com/go-countryApi/config"
	"github.com/go-countryApi/models"
	"github.com/go-countryApi/utils"
)

var gcache *cache.InMemCache

// creating global in mem-cache
func CreateInMemCache() {
	gcache = cache.NewCache()
}

// GetCountryData : to get country data of country name passed in query params of url
// if data found for country name in cache then return it
// else get the data from api, store it in cache and return it
func GetCountryData(w http.ResponseWriter, r *http.Request) {
	hasSearch := r.URL.Query().Has("name")
	name := r.URL.Query().Get("name")

	w.Header().Set("Content-Type", "application/json")

	var apiRes models.ApiRes
	apiRes.Status = false

	if !hasSearch {
		w.WriteHeader(http.StatusBadRequest)
		apiRes.Result.Error = "invalid payload received"
		io.WriteString(w, utils.ApiResToRet(apiRes))
		return
	}

	var countryData interface{}
	var ok bool
	if countryData, ok = gcache.Get(strings.ToLower(name)); !ok {
		// check from api and store in cache
		log.Printf(":::::: [ GetCountryData ] no data found in cache getting from api ::::::::")
		var err error

		apiRes, err = getCtyDataFrmUrl(r.Context(), name)
		if err != nil {
			apiRes.Result.Error = err.Error()
			if err.Error() == "internal server error" {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, utils.ApiResToRet(apiRes))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, utils.ApiResToRet(apiRes))
			return
		}

		io.WriteString(w, utils.ApiResToRet(apiRes))
		return
	}

	apiRes = countryData.(models.ApiRes)
	apiRes.Status = true
	log.Printf("[ GetCountryData ] data found from cache :: %v", apiRes.Result)
	io.WriteString(w, utils.ApiResToRet(apiRes))
}

// geCtyDataFrmUrl : to get country data from the api url of restcountries
// store the same in cache first then return
func getCtyDataFrmUrl(ctx context.Context, name string) (models.ApiRes, error) {

	var apiRes models.ApiRes
	apiRes.Status = false
	var retErr error = fmt.Errorf("internal server error")

	fetchUrl := config.AppConfig.ExtUrls.CountiresFetchUrl + name

	// fetch specfic daata only
	base, err := url.Parse(fetchUrl)
	if err != nil {
		// return err or do something
		log.Printf(" [ getCtyDataFrmUrl ] error while parsing url:: %v", err.Error())
		return apiRes, retErr
	}

	params := url.Values{}
	params.Add("fields", "name,currencies,capital,population")
	base.RawQuery = params.Encode()

	statusCode, res := utils.HttpReq(ctx, "GET", base.String(), "", "", nil, config.AppConfig.ExtUrls.CtxTimeoutUrl)

	if statusCode == 200 {
		var resJson []models.CountryData
		err = json.Unmarshal(res, &resJson)
		if err != nil {
			log.Printf("[ getCtyDataFrmUrl ] err while unmarshalling %v", err.Error())
			return apiRes, retErr
		}
		// log.Printf("data recieved from api : %v", resJson)

		for i := 0; i < len(resJson); i++ {
			if strings.EqualFold(resJson[i].Name.Common, name) {
				apiRes.Result.Name = name
				apiRes.Result.Capital = resJson[i].Capital[0]

				for _, v := range resJson[i].Currencies {
					apiRes.Result.Currency = v.Symbol
				}

				apiRes.Result.Population = resJson[i].Population
				apiRes.Status = true

				gcache.Set(strings.ToLower(name), apiRes)
				log.Printf("[ getCtyDataFrmUrl ] data found from api :: %v", resJson[i])
				retErr = nil
				return apiRes, retErr
			}
		}
		retErr = fmt.Errorf("no data found for this country")
	} else if statusCode == 404 {
		retErr = fmt.Errorf("no data found for this country")
	}

	return apiRes, retErr
}
