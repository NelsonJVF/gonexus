package gonexus

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"time"
	"errors"
)

// hTTPResponse is a strut for HTTP Response from Nexus
type hTTPResponse struct {
	Header 	http.Header
	Body 		[]byte
}

// Configuration is a struct for Nexus access information
type Configuration struct {
	Lable string `yaml:"lable"` // Some projects have more than one Nexus, so just lable as you wish
	User string  `yaml:"user"` // Username for Nexus
	Pass string  `yaml:"pass"` // Password from Nexus Username
	URL string   `yaml:"url"` // Url to Nexus hostname + port
	Timeout int  `yaml:"timeout"` // Timeout to call Nexus
}

// SearchResponse is a struct of the Nexus Search Response
type SearchResponse struct {
	TotalCount     int  `json:"totalCount"`
	From           int  `json:"from"`
	Count          int  `json:"count"`
	TooManyResults bool `json:"tooManyResults"`
	Collapsed      bool `json:"collapsed"`
	RepoDetails    []struct {
		RepositoryID           string `json:"repositoryId"`
		RepositoryName         string `json:"repositoryName"`
		RepositoryContentClass string `json:"repositoryContentClass"`
		RepositoryKind         string `json:"repositoryKind"`
		RepositoryPolicy       string `json:"repositoryPolicy"`
		RepositoryURL          string `json:"repositoryURL"`
	} `json:"repoDetails"`
	Data []struct {
		GroupID                   string `json:"groupId"`
		ArtifactID                string `json:"artifactId"`
		Version                   string `json:"version"`
		LatestRelease             string `json:"latestRelease"`
		LatestReleaseRepositoryID string `json:"latestReleaseRepositoryId"`
		HighlightedFragment       string `json:"highlightedFragment"`
		ArtifactHits              []struct {
			RepositoryID  string `json:"repositoryId"`
			ArtifactLinks []struct {
				Extension string `json:"extension"`
			} `json:"artifactLinks"`
		} `json:"artifactHits"`
		LatestSnapshot             string `json:"latestSnapshot,omitempty"`
		LatestSnapshotRepositoryID string `json:"latestSnapshotRepositoryId,omitempty"`
	} `json:"data"`
}

// Config variable of type Configuration stores the information of Nexus server(s)
var Config []Configuration

// HTTPRequest represents a a HTTP request for Nexus REST API
func hTTPRequest(project string, urlPath string, jsonBody string) (hTTPResponse, error) {
	var hTTPResp hTTPResponse
	var user string
	var pass string
	var url string
	var timeout int
	var httpMethod = "GET"

	for _, c := range Config {
		if c.Lable == project {
			user = c.User
			pass = c.Pass
			url = c.URL
			timeout = c.Timeout
		}
	}

	if(len(url) == 0) {
		err := errors.New("Nexus configuration is missing, for project " + project)
		return hTTPResp, err
	}

	url = fmt.Sprintf("%s%s", url, urlPath)

	if(len(jsonBody) > 0) {
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

	r.SetBasicAuth(user, pass)
	r.Header.Add("Accept", "application/json")

	resp, errDo := client.Do(r)
	if errDo != nil {
		err := errors.New("Error connecting to Nexus server (" + project + ")." + errDo.Error())
		return hTTPResp, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := errors.New("Error reading Body: " + err.Error())
		return hTTPResp, err
	}

	hTTPResp.Header = resp.Header
	hTTPResp.Body = body

	return hTTPResp, nil
}

// RequestSearch represents a search feature in Nexus rest api, we should specify the project from that item
func RequestSearch(project string, query string) (SearchResponse, error) {
	var urlSearchPath string
	var data SearchResponse

	urlSearchPath = fmt.Sprintf("nexus/service/local/lucene/search?q=%s", query)

	response, error := hTTPRequest(project, urlSearchPath, "")
	if error != nil {
		return data, error
	}

	err := json.Unmarshal(response.Body, &data)
	if err != nil {
		log.Printf("gonexus.RequestSearch - json.Unmarshal err   #%v ", err)
	}

	return data, nil
}
