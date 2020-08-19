package containerregistry

import (
	"fmt"
	"math"
	"net/http"
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
	repositories := gcrClient.getRepositories()

	repositoriesChannel := make(chan Repository, len(repositories))
	containerImagesChannel := make(chan []ContainerImage, len(repositories))

	// spawn max 40 goroutines, if repos are less than 40, try to GET them all concurrently
	workersNum := int(math.Min(40, float64(len(repositories))))

	for i := 1; i <= workersNum; i++ {
		go gcrClient.imagesRepoFetchWorker(repositoriesChannel, containerImagesChannel)
	}

	for _, repo := range repositories {
		repositoriesChannel <- repo
	}

	// after all the channels have been filled, it's time to read from them now
	allImages := []ContainerImage{}
	for range repositories {
		allImages = append(allImages, <-containerImagesChannel...)
	}

	fmt.Println("All images' length is", len(allImages))

	return allImages
}
