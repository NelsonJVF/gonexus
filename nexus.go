package gonexus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPRequest represents a a HTTP request for Nexus REST API
func hTTPRequest(URL string, urlPath string, username string, password string, jsonBody string) (hTTPResponse, error) {
	var hTTPResp hTTPResponse
	var timeout int
	var httpMethod = "GET"

	if len(URL) == 0 {
		return hTTPResp, errors.New("Nexus configuration is missing, for project ")
	}

	url := fmt.Sprintf("%s%s", URL, urlPath)

	if len(jsonBody) > 0 {
		httpMethod = "POST"
	}

	timeoutVal := time.Duration(time.Duration(timeout) * time.Second)
	client := &http.Client{
		Timeout: timeoutVal,
	}
	r, erroHTTP := http.NewRequest(httpMethod, url, bytes.NewBufferString(jsonBody))
	if erroHTTP != nil {
		return hTTPResp, erroHTTP
	}

	r.SetBasicAuth(username, password)
	r.Header.Add("Accept", "application/json")

	resp, errDo := client.Do(r)
	if errDo != nil {
		return hTTPResp, errors.New("Error connecting to Nexus server." + errDo.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return hTTPResp, errors.New("Error reading Body: " + err.Error())
	}

	hTTPResp.Header = resp.Header
	hTTPResp.Body = body

	return hTTPResp, nil
}

// RequestSearch represents a search feature in Nexus rest api, we should specify the project from that item
func RequestSearch(URL string, urlPath string, username string, password string, query string) (SearchResponse, error) {
	var urlSearchPath string
	var data SearchResponse

	urlSearchPath = fmt.Sprintf("nexus/service/local/lucene/search?q=%s", query)

	response, err := hTTPRequest(URL, urlSearchPath, username, password, "")
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(response.Body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
