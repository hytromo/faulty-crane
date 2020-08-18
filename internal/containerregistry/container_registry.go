package containerregistry

import (
	"fmt"
	"net/http"
)

// GCRClient is used for creating Google Container Registry clients
type GCRClient struct {
	Link      string
	AccessKey string
	client    *http.Client
}

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

	fmt.Println("Catalog contains", len(repositories), "repos")

	if len(repositories) == 0 {
		return nil
	}

	return nil
}
