package gcr

import (
	"encoding/json"
	"fmt"
	"net/http"

	cr "github.com/hytromo/faulty-crane/internal/containerregistry"
	myhttp "github.com/hytromo/faulty-crane/internal/http"
	"github.com/hytromo/faulty-crane/internal/utils/stringutil"
	log "github.com/sirupsen/logrus"
)

// GoogleContainerRegistryClient is a GCR client
type GoogleContainerRegistryClient struct {
	httpClient myhttp.Client
}

// Login logins to GCR
func (client *GoogleContainerRegistryClient) Login(username string, password string) error {
	// GCR client does not need to login to get any kind of token, it just needs to specify the token in each request
	// regardless, we need to specify this function to comply with the CRClient interface
	return nil
}

// DeleteImage deletes an image from GCR
func (client *GoogleContainerRegistryClient) DeleteImage(imageRepo string, image cr.ContainerImage, silentErrors bool) error {
	// all the tags of the image need to be deleted first
	var err error

	for _, tag := range image.Tag {
		err = client.httpClient.DeleteRequestTo("/"+imageRepo+"/manifests/"+tag, true, silentErrors)
		if err != nil {
			return err
		}
	}

	for _, digest := range image.Digest {
		// after all the image tags have been deleted, we can delete the image itself
		err = client.httpClient.DeleteRequestTo("/"+imageRepo+"/manifests/"+digest, true, silentErrors)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAllRepos parses or repos of GCR
func (client *GoogleContainerRegistryClient) GetAllRepos() []string {
	repositories := []string{}

	catalogResp := cr.CatalogDTO{
		Next: "/_catalog", // initial request
	}

	for {
		bodyBytes, err := client.httpClient.GetRequestTo(catalogResp.Next)

		if err != nil {
			log.Fatalf("Error on api call: %v", err.Error())
		}

		catalogResp = cr.CatalogDTO{}
		err = json.Unmarshal(bodyBytes, &catalogResp)

		repositories = append(repositories, catalogResp.Repositories...)

		if err != nil {
			log.Fatalf("Invalid api call response (%v): %v", string(bodyBytes), err.Error())
		}

		if catalogResp.Next == "" { // no more pages to GET
			break
		} else { // more pages to GET
			// remove the prefix of the link as our gcr client works with suffixes
			catalogResp.Next = stringutil.TrimLeftChars(catalogResp.Next, len(client.httpClient.BaseURL))
			if catalogResp.Next[0] != '/' {
				catalogResp.Next = "/" + catalogResp.Next
			}
		}
	}

	return repositories
}

// ParseRepo parses a specific repo
func (client *GoogleContainerRegistryClient) ParseRepo(repositoryLink string) cr.Repository {
	repository := cr.Repository{
		Link:   repositoryLink,
		Images: []cr.ContainerImage{},
	}

	listTagsResp := cr.ListTagsDTO{
		Next: "/" + repositoryLink + "/tags/list", // initial request
	}

	for {
		bodyBytes, err := client.httpClient.GetRequestTo(listTagsResp.Next)

		if err != nil {
			log.Fatalf("Error on api call: %v", err.Error())
		}

		listTagsResp = cr.ListTagsDTO{}
		err = json.Unmarshal(bodyBytes, &listTagsResp)

		for digest, image := range listTagsResp.Manifest {
			image.Digest = []string{digest}
			image.Repo = client.httpClient.BaseURL + "/" + repositoryLink
			repository.Images = append(repository.Images, image)
		}

		if err != nil {
			log.Fatalf("Invalid api call response (%v): %v", string(bodyBytes), err.Error())
		}

		if listTagsResp.Next == "" { // no more pages to GET
			break
		}
	}

	return repository
}

// NewGCRClientParams are the required parameters to build a GCR client
type NewGCRClientParams struct {
	// one of gcr.io, us.gcr.io, eu.gcr.io, asia.gcr.io https://cloud.google.com/container-registry/docs/overview#registries
	Hostname string
	// e.g. the result of `gcloud auth print-access-token`
	Token string
}

// NewGCRClient builds a new GCR client
func NewGCRClient(params NewGCRClientParams) cr.Client {
	return &GoogleContainerRegistryClient{
		httpClient: myhttp.NewClient(myhttp.NewClientParams{
			BaseURL: fmt.Sprintf("https://%s/v2", params.Hostname),
			InjectAuthInRequest: func(req *http.Request) {
				req.SetBasicAuth("_token", params.Token)
			},
		}),
	}
}
