package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-countryApi/config"
	"github.com/go-countryApi/models"
	"github.com/go-countryApi/utils"
	"github.com/stretchr/testify/require"
)

// TestGetCountryDataAPI : to test the getCountryData handler on multiple case scenarios
func TestGetCountryDataAPI(t *testing.T) {
	config.GetConfig()

	CreateInMemCache()
	randCtryData := utils.RandomCountryData()

	invalidPayloadD := models.ApiRes{
		Status: false,
		Result: models.Result{
			Error: "invalid payload received",
		},
	}

	noDataFoundD := invalidPayloadD
	noDataFoundD.Result.Error = "no data found for this country"

	internalServerErrorD := invalidPayloadD
	internalServerErrorD.Result.Error = "internal server error"

	testCases := []struct {
		name          string
		queryField    string
		countryName   string
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			queryField:  "name",
			countryName: randCtryData.Result.Name,
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchCountry(t, recoder.Body, randCtryData)
			},
		},
		{
			name:        "InvalidPayload",
			queryField:  "names",
			countryName: randCtryData.Result.Name,
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
				requireBodyMatchCountry(t, recoder.Body, invalidPayloadD)
			},
		},
		{
			name:        "NoDataFound",
			queryField:  "name",
			countryName: randCtryData.Result.Name + "s",
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
				requireBodyMatchCountry(t, recoder.Body, noDataFoundD)
			},
		},
		{
			name:        "InternalServerError",
			queryField:  "name",
			countryName: randCtryData.Result.Name,
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
				requireBodyMatchCountry(t, recoder.Body, internalServerErrorD)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			url := config.AppConfig.Prefix + "/search"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add(tc.queryField, tc.countryName)
			request.URL.RawQuery = q.Encode()

			GetCountryData(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

// BenchmarkSerialGetCountryData : bechmarking serail access to GetcountryData handler
func BenchmarkSerialGetCountryData(b *testing.B) {
	config.GetConfig()
	CreateInMemCache()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		recorder := httptest.NewRecorder()

		url := config.AppConfig.Prefix + "/search"
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(b, err)

		q := request.URL.Query()
		q.Add("name", utils.RandomCountryName())
		request.URL.RawQuery = q.Encode()
		GetCountryData(recorder, request)
	}
}

// BenchmarkParallelGetCountryData : bechmarking parallel access to GetcountryData handler
func BenchmarkParallelGetCountryData(b *testing.B) {
	config.GetConfig()
	CreateInMemCache()

	for i := 0; i < b.N; i++ {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				recorder := httptest.NewRecorder()

				url := config.AppConfig.Prefix + "/search"
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(b, err)

				q := request.URL.Query()
				q.Add("name", utils.RandomCountryName())
				request.URL.RawQuery = q.Encode()
				GetCountryData(recorder, request)
			}
		})
	}
}

func requireBodyMatchCountry(t *testing.T, body *bytes.Buffer, country models.ApiRes) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCountrty models.ApiRes
	err = json.Unmarshal(data, &gotCountrty)
	require.NoError(t, err)
	require.Equal(t, country, gotCountrty)
}
