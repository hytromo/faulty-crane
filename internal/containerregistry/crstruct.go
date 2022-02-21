package containerregistry

import "github.com/hytromo/faulty-crane/internal/keepreasons"

// Repository is a struct that holds information about a container registry's repository
type Repository struct {
	// Link is the relative link, also refered to as "image name" on the documentation, each repository can contain a lot of images with different tags and manifests
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
	Digest         []string
	Repo           string               // Repo is the name of the image's repository without the tag in the form e.g. eu.gcr.io/faulty-crane-project/faulty-crane-test
	KeptData       keepreasons.KeptData `json:",omitempty"`
}

// RepoDeletionResult is the repository deletion result
type RepoDeletionResult struct {
	ShouldDeleteCount    int
	ManagedToDeleteCount int
}

// CatalogDTO is the Data Transfer Object for the catalog api call
type CatalogDTO struct {
	// Next is used for pagination purposes, it contains the next URL we need to GET for the next page
	Next         string
	Repositories []string
}

// ListTagsDTO is the Data Transfer Object for the list tags api call
type ListTagsDTO struct {
	// Manifest keys are the image digest
	Manifest map[string]ContainerImage
	Name     string
	Tags     []string
	Next     string
}

// Client is used for implementing container registry clients
type Client interface {
	Login(username string, password string) error
	DeleteImage(imageRepo string, image ContainerImage, silentErrors bool) error
	GetAllRepos() []string
	ParseRepo(repositoryLink string) Repository
}
