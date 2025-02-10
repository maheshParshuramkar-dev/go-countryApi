package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-countryApi/models"
)

func HttpReq(ctx context.Context, method, http_url, reqBody, proxyUrl string, headers map[string]string, ctxTimeout int) (int, []byte) {
	var (
		req  io.Reader
		resp *http.Request
		err  error
	)

	if method == "GET" {
		req = nil
	} else {
		req = bytes.NewBuffer([]byte(reqBody))
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(ctxTimeout)*time.Millisecond)
	defer cancel()

	resp, err = http.NewRequestWithContext(ctxWithTimeout, method, http_url, req)

	if err != nil {
		log.Println("error while executing curl req ", err.Error())
		return 500, nil
	}

	// setting headers in the resp
	for k, v := range headers {
		resp.Header.Set(k, v)
	}

	client := &http.Client{}

	if proxyUrl != "" {
		log.Println("proxy url: ", proxyUrl)

		proxyUrlP, err := url.Parse(proxyUrl)
		if err != nil {
			log.Println("error while parsing porxy url ", err.Error())
			return 500, nil
		}

		transport := http.Transport{
			Proxy:             http.ProxyURL(proxyUrlP),
			DisableKeepAlives: true,
		}
		client.Transport = &transport
	}

	response, err := client.Do(resp)
	if err != nil {
		log.Println("error while executing client req ", err.Error())
		return 500, nil
	}
	defer response.Body.Close()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("error while reading resp body received")
	}
	return response.StatusCode, body
}

func ApiResToRet(res models.ApiRes) string {
	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Printf("[ ApiResToRet ] err while marshalling %v", err.Error())
		return "internal server error"
	}

	return string(resBytes)
}
