package containerregistry

/**
 * Api docs at https://docs.docker.com/registry/spec/api/
 * These api methods slightly process the responses before returning
 */

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/hytromo/faulty-crane/internal/utils/stringutil"
)

// catalogDTO is the Data Transfer Object for the catalog api call
type catalogDTO struct {
	// Next is used for pagination purposes, it contains the next URL we need to GET for the next page
	Next         string
	Repositories []string
}

// listTagsDTO is the Data Transfer Object for the list tags api call
type listTagsDTO struct {
	// Manifest keys are the image digest
	Manifest map[string]ContainerImage
	Name     string
	Tags     []string
	Next     string
}

func (gcrClient GCRClient) getRepositories() []string {
	repositories := []string{}

	catalogResp := catalogDTO{
		Next: "/_catalog", // initial request
	}

	for {
		bodyBytes := gcrClient.getRequestTo(catalogResp.Next)

		catalogResp = catalogDTO{}
		err := json.Unmarshal(bodyBytes, &catalogResp)

		repositories = append(repositories, catalogResp.Repositories...)

		if err != nil {
			log.Fatalf("Invalid api call response (%v): %v", string(bodyBytes), err.Error())
		}

		if catalogResp.Next == "" { // no more pages to GET
			break
		} else { // more pages to GET
			// remove the prefix of the link as our gcr client works with suffixes
			catalogResp.Next = stringutil.TrimLeftChars(catalogResp.Next, len(gcrClient.getBaseURL()))
			if catalogResp.Next[0] != '/' {
				catalogResp.Next = "/" + catalogResp.Next
			}
		}
	}

	return repositories
}

func (gcrClient GCRClient) listTags(repositoryLink string) Repository {
	repository := Repository{
		Link:   repositoryLink,
		Images: []ContainerImage{},
	}

	listTagsResp := listTagsDTO{
		Next: "/" + repositoryLink + "/tags/list", // initial request
	}

	for {
		bodyBytes := gcrClient.getRequestTo(listTagsResp.Next)

		listTagsResp = listTagsDTO{}
		err := json.Unmarshal(bodyBytes, &listTagsResp)

		for digest, image := range listTagsResp.Manifest {
			image.Digest = digest
			image.Repo = gcrClient.Host + "/" + repositoryLink
			repository.Images = append(repository.Images, image)
		}

		if err != nil {
			log.Fatalf("Invalid api call response (%v): %v", string(bodyBytes), err.Error())
		}

		if listTagsResp.Next == "" { // no more pages to GET
			break
		} else { // more pages to GET
			// remove the prefix of the link as our gcr client works with suffixes
			listTagsResp.Next = stringutil.TrimLeftChars(listTagsResp.Next, len(gcrClient.getBaseURL()))
			if listTagsResp.Next[0] != '/' {
				listTagsResp.Next = "/" + listTagsResp.Next
			}
		}
	}

	return repository
}

func (gcrClient GCRClient) deleteImage(imageRepo string, image ContainerImage) {
	// all the tags of the image need to be deleted first
	for _, tag := range image.Tag {
		gcrClient.deleteRequestTo("/"+imageRepo+"/manifests/"+tag, true)
	}

	// after all the image tags have been deleted, we can delete the image itself
	gcrClient.deleteRequestTo("/"+imageRepo+"/manifests/"+image.Digest, false)
}
