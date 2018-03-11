package gonexus

import "net/http"

// hTTPResponse is a strut for HTTP Response from Nexus
type hTTPResponse struct {
	Header http.Header
	Body   []byte
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
