package containerregistry

import (
	"fmt"
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

// ContainerImage contains all the data that are relevant to an image on the registry
type ContainerImage struct {
	ImageSizeBytes string
	LayerID        string `json:"layerId"`
	MediaType      string
	Tag            []string
	TimeCreatedMs  string
	TimeUploadedMs string
	Digest         string
}

// MakeGCRClient builds a new GCRClient instance, adding the missing default values e.g. http client
func MakeGCRClient(client GCRClient) GCRClient {
	gcrClient := client
	gcrClient.client = &http.Client{}
	return gcrClient
}

func (gcrClient GCRClient) imagesRepoFetchWorker(repositories <-chan Repository, containerImages chan<- []ContainerImage) {
	for repo := range repositories {
		containerImages <- gcrClient.listTags(repo)
	}
}

// GetAllImages returns all the images in a docker container registry by spawning multiple workers to make it go faster
func (gcrClient GCRClient) GetAllImages() []ContainerImage {
	log.Info("Getting all the repos of the registry...")

	repositories := gcrClient.getRepositories()
	repositoriesCount := len(repositories)

	bar := pb.Full.Start(repositoriesCount)

	repositoriesChannel := make(chan Repository, repositoriesCount)
	containerImagesChannel := make(chan []ContainerImage, repositoriesCount)

	// spawn max 40 goroutines, if repos are less than 40, try to GET them all concurrently
	workersNum := int(math.Min(40, float64(repositoriesCount)))

	log.Info("Fetching the images of ", repositoriesCount, " repo(s), using ", workersNum, " routines")

	for i := 1; i <= workersNum; i++ {
		go gcrClient.imagesRepoFetchWorker(repositoriesChannel, containerImagesChannel)
	}

	for _, repo := range repositories {
		// feed the jobs to the workers
		repositoriesChannel <- repo
	}

	// while the jobs are being done by the workers, we are merging all the results into one
	allImages := []ContainerImage{}
	for range repositories {
		allImages = append(allImages, <-containerImagesChannel...)
		bar.Increment()
	}

	bar.Finish()
	fmt.Println("All images' length is", len(allImages))

	return allImages
}
