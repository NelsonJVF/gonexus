// Copyright 2017 NelsonJVF. All rights reserved.
package gonexus

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"time"
)

type hTTPResponse struct {
	Header 	http.Header
	Body 		[]byte
}

/*
	Struct for Nexus access information
 */
type Configuration struct {
	Lable string `yaml:"lable"` // Some projects have more than one Jira, so just lable as you wish
	User string  `yaml:"user"` // Username for Jira
	Pass string  `yaml:"pass"` // Password from Jira Username
	Url string   `yaml:"url"` // URL to Jira hostname + port
	Timeout int  `yaml:"timeout"` // URL to Jira hostname + port
}

/*
	Nexus Search Response Struct
 */
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

var Config []Configuration

/*
	Generic HTTP caller
 */
func HTTPRequest(project string, urlPath string, jsonBody string) hTTPResponse {
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
			url = c.Url
			timeout = c.Timeout
		}
	}

	if(len(url) == 0) {
		log.Printf(" ---------- Nexus configuration is missing  ---------- ")
		log.Printf("\t For project " + project)
		return hTTPResp
	}

	url = fmt.Sprintf("%s%s", url, urlPath)

	if(len(jsonBody) > 0) {
		httpMethod = "POST"
	}

	timeoutVal := time.Duration(time.Duration(timeout) * time.Second)
	client := &http.Client{
		Timeout: timeoutVal,
	}
	r, _ := http.NewRequest(httpMethod, url, bytes.NewBufferString(jsonBody))

	r.SetBasicAuth(user, pass)
	r.Header.Add("Accept", "application/json")

	resp, errDo := client.Do(r)
	if errDo != nil {
		log.Println("Error connecting to Nexus server (" + project + ").")
		return hTTPResp
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err   #%v ", err)
	}

	hTTPResp.Header = resp.Header
	hTTPResp.Body = body

	return hTTPResp
}

/*
	Search in Jira, we should specify the project from that item
 */
func RequestSearch(project string, query string) SearchResponse {
	var urlSearchPath string
	var data SearchResponse

	urlSearchPath = fmt.Sprintf("nexus/service/local/lucene/search?q=%s", query)

	response := HTTPRequest(project, urlSearchPath, "")

	err := json.Unmarshal(response.Body, &data)
	if err != nil {
		log.Printf("json.Unmarshal err   #%v ", err)
	}

	return data
}