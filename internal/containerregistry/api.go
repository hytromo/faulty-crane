package containerregistry

/**
 * Api docs at https://docs.docker.com/registry/spec/api/
 */

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/hytromo/faulty-crane/internal/utils/stringutil"
)

// catalogDTO is the Data Transfer Object for the catalog api call
type catalogDTO struct {
	// Next is used for pagination purposes, it contains the next URL we need to GET for the next page
	Next         string `json:",omitempty"`
	Repositories []Repository
}

// Repository describes a url where the registry holds docker images
type Repository string

func (gcrClient GCRClient) getRepositories() []Repository {
	repositories := []Repository{}

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
			catalogResp.Next = stringutil.TrimLeftChars(catalogResp.Next, len(gcrClient.Link))
			if catalogResp.Next[0] != '/' {
				catalogResp.Next = "/" + catalogResp.Next
			}
		}
	}

	return repositories
}
