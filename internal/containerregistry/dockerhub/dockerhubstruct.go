package dockerhub

import myhttp "github.com/hytromo/faulty-crane/internal/http"

// RegistryClient is a dockerhub registry client
type RegistryClient struct {
	httpClient myhttp.Client
	token      string
	namespace  string
}

// UsersLoginDTO is the Data Transfer Object for the /users/login api call
type UsersLoginDTO struct {
	Token string
}

// RepositoryResultDTO is the DTO of the corresponding API call
type RepositoryResultDTO struct {
	User              string
	Name              string
	Namespace         string
	RepositoryType    string
	Status            int
	Description       string
	IsPrivate         bool
	IsAutomated       bool
	CanEdit           bool
	StarCount         int
	PullCount         int
	LastUpdated       string
	IsMigrated        bool
	CollaboratorCount int
	Affiliation       string
	HubUser           string
}

// RepositoryDTO is the Data Transfer Object for the /repositories/{namespace} api call
type RepositoryDTO struct {
	Count    int
	Next     string
	Previous string
	Results  []RepositoryResultDTO
}

// TagResultImageDTO is the DTO of the corresponding API call
type TagResultImageDTO struct {
	Architecture string
	Features     string
	Variant      string
	Digest       string
	Os           string
	OsFeatures   string
	OsVersion    string
	Size         int64
	Status       string
	LastPulled   string
	LastPushed   string
}

// TagResultDTO is the DTO of the corresponding API call
type TagResultDTO struct {
	Creator             int64
	ID                  int64
	ImageID             string
	Images              []TagResultImageDTO
	LastUpdated         string
	LastUpdater         int64
	LastUpdaterUsername string
	Name                string
	Repository          int64
	FullSize            int64
	V2                  bool
	TagStatus           string
	TagLastPulled       string
	TagLastPushed       string
}

// TagsDTO is the Data Transfer Object for the /repositories/{namespace}/{repo}/tags api call
type TagsDTO struct {
	Count    int
	Next     string
	Previous string
	Results  []TagResultDTO
}
