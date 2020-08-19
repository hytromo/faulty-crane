package containerregistry

import (
	"net/http"

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
	Link string
}

// MakeGCRClient builds a new GCRClient instance, adding the missing default values e.g. http client
func MakeGCRClient(client GCRClient) GCRClient {
	gcrClient := client
	gcrClient.client = &http.Client{}
	return gcrClient
}

// GetAllImages returns all the images in a docker container registry
func (gcrClient GCRClient) GetAllImages() []ContainerImage {
	repositories := gcrClient.getRepositories()

	log.Info("Catalog contains", len(repositories), "repos")

	if len(repositories) == 0 {
		return nil
	}

	return nil
}
