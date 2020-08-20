package containerregistry

import (
	"math"
	"net/http"

	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
)

// GCRClient is used for creating Google Container Registry clients
type GCRClient struct {
	Host      string
	AccessKey string
	client    *http.Client
}

// Repository is a struct that holds information about a container registry's repository
type Repository struct {
	Link   string
	Images []ContainerImage
}

// ContainerImage contains all the data that are relevant to an image on the registry
type ContainerImage struct {
	ImageSizeBytes string
	LayerID        string `json:"layerId"`
	MediaType      string
	Tag            []string
	TimeCreatedMs  string
	TimeUploadedMs string
	Digest         string
	Repo           string // Repo is the name of the image's repository without the tag in the form e.g. eu.gcr.io/faulty-crane-project/faulty-crane-test
}

// MakeGCRClient builds a new GCRClient instance, adding the missing default values e.g. http client
func MakeGCRClient(client GCRClient) GCRClient {
	gcrClient := client
	gcrClient.client = &http.Client{}
	return gcrClient
}

func (gcrClient GCRClient) fetchRepoImagesWorker(repositoryLinks <-chan string, parsedRepositories chan<- Repository) {
	for repo := range repositoryLinks {
		parsedRepositories <- gcrClient.listTags(repo)
	}
}

// GetAllRepos finds all the repositories and parses them finding the images they contain
func (gcrClient GCRClient) GetAllRepos() []Repository {
	log.Info("Getting all the repos of the registry...")

	repositories := gcrClient.getRepositories()
	repositoriesCount := len(repositories)

	bar := pb.Full.Start(repositoriesCount)

	repositoryLinksChan := make(chan string, repositoriesCount) // jobs
	parsedReposChan := make(chan Repository, repositoriesCount) // results

	// spawn max 40 goroutines, if repos are less than 40, try to list them all concurrently
	workersNum := int(math.Min(40, float64(repositoriesCount)))

	log.Info("Fetching the images of ", repositoriesCount, " repo(s), using ", workersNum, " routine(s)")

	for i := 1; i <= workersNum; i++ {
		go gcrClient.fetchRepoImagesWorker(repositoryLinksChan, parsedReposChan)
	}

	for _, repo := range repositories {
		// feed the jobs to the workers
		repositoryLinksChan <- repo
	}

	// while the jobs are being done by the workers, we are merging all the results into one
	allRepos := []Repository{}
	for range repositories {
		allRepos = append(allRepos, <-parsedReposChan)
		bar.Increment()
	}

	bar.Finish()

	return allRepos
}
