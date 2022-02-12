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

func (gcrClient GCRClient) getRepositories() []string {
	repositories := []string{}

	catalogResp := CatalogDTO{
		Next: "/_catalog", // initial request
	}

	for {
		bodyBytes := gcrClient.getRequestTo(catalogResp.Next)

		catalogResp = CatalogDTO{}
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

	listTagsResp := ListTagsDTO{
		Next: "/" + repositoryLink + "/tags/list", // initial request
	}

	for {
		bodyBytes := gcrClient.getRequestTo(listTagsResp.Next)

		listTagsResp = ListTagsDTO{}
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

func (gcrClient GCRClient) DeleteImage(imageRepo string, image ContainerImage, silentErrors bool) bool {
	// all the tags of the image need to be deleted first
	atLeastOneTagFailed := false
	for _, tag := range image.Tag {
		atLeastOneTagFailed = !gcrClient.deleteRequestTo("/"+imageRepo+"/manifests/"+tag, true, silentErrors)
		if atLeastOneTagFailed {
			break
		}
	}

	if atLeastOneTagFailed {
		return false
	}

	// after all the image tags have been deleted, we can delete the image itself
	return gcrClient.deleteRequestTo("/"+imageRepo+"/manifests/"+image.Digest, true, silentErrors)
}
