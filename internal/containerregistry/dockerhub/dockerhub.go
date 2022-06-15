package dockerhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Rican7/conjson"
	"github.com/Rican7/conjson/transform"
	cr "github.com/hytromo/faulty-crane/internal/containerregistry"
	myhttp "github.com/hytromo/faulty-crane/internal/http"
	log "github.com/sirupsen/logrus"
)

var baseURL = "https://hub.docker.com/v2"

// Login logs in into dockerhub
func (client *RegistryClient) Login(username string, password string) error {
	jsonPayload, _ := json.Marshal(map[string]interface{}{
		"username": username,
		"password": password,
	})

	bodyBytes, err := client.httpClient.PostRequestTo("/users/login", jsonPayload, false, false)

	if err != nil {
		log.Fatalf("Failed to login to dockerhub with error: %v", err.Error())
	}

	loginResp := UsersLoginDTO{}
	err = json.Unmarshal(bodyBytes, &loginResp)

	if err != nil {
		log.Fatalf("Failed to parse dockerhub login response with error: %v", err.Error())
	}

	client.token = loginResp.Token
	client.httpClient.InjectAuthInRequest = func(req *http.Request) {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.token))
	}

	return err
}

// DeleteImage delets an image from dockerhub
func (client *RegistryClient) DeleteImage(imageRepo string, image cr.ContainerImage, silentErrors bool) error {
	// all the tags of the image need to be deleted first
	var err error

	for _, tag := range image.Tag {
		err = client.httpClient.DeleteRequestTo("/repositories/"+imageRepo+"/tags/"+tag+"/", true, silentErrors)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAllRepos parses all dockerhub repositories
func (client *RegistryClient) GetAllRepos() []string {
	repositories := []string{}

	repositoryResp := RepositoryDTO{
		Next: fmt.Sprintf("/repositories/%s?page_size=100", client.namespace), // initial request
	}

	for {
		bodyBytes, err := client.httpClient.GetRequestTo(repositoryResp.Next)

		if err != nil {
			log.Fatalf("Error on api call: %v", err.Error())
		}

		repositoryResp = RepositoryDTO{}

		err = json.Unmarshal(
			[]byte(bodyBytes),
			conjson.NewUnmarshaler(&repositoryResp, transform.ConventionalKeys()),
		)

		if err != nil {
			log.Fatalf("Invalid api call response (%v): %v", string(bodyBytes), err.Error())
		}

		err = json.Unmarshal(bodyBytes, &repositoryResp)

		if err != nil {
			log.Fatalf("Invalid api call response (%v): %v", string(bodyBytes), err.Error())
		}

		for _, result := range repositoryResp.Results {
			repositories = append(repositories, fmt.Sprintf("%s/%s", client.namespace, result.Name))
		}

		if repositoryResp.Next == "" { // no more pages to GET
			break
		}
	}

	return repositories
}

// ParseRepo parses a specific repository
func (client *RegistryClient) ParseRepo(repositoryLink string) cr.Repository {
	repository := cr.Repository{
		Link:   repositoryLink,
		Images: []cr.ContainerImage{},
	}

	listTagsResp := TagsDTO{
		Next: "/repositories/" + repositoryLink + "/tags?page_size=100", // initial request
	}

	for {
		bodyBytes, err := client.httpClient.GetRequestTo(listTagsResp.Next)

		if err != nil {
			log.Fatalf("Error on api call: %v", err.Error())
		}

		listTagsResp = TagsDTO{}

		err = json.Unmarshal(
			[]byte(bodyBytes),
			conjson.NewUnmarshaler(&listTagsResp, transform.ConventionalKeys()),
		)

		if err != nil {
			log.Fatalf("Error on api call: %v", err.Error())
		}

		timeLayout := "2006-01-02T15:04:05.999999Z"
		for _, result := range listTagsResp.Results {
			repoImage := cr.ContainerImage{}
			repoImage.Digest = []string{}
			repoImage.Tag = []string{result.Name}
			var totalImageSize int64 = 0
			t, err := time.Parse(timeLayout, result.TagLastPushed)

			if err == nil {
				updatedMs := strconv.FormatInt(t.UTC().UnixMilli(), 10)
				repoImage.TimeCreatedMs = updatedMs
				repoImage.TimeUploadedMs = updatedMs
			}

			repoImage.LayerID = strconv.FormatInt(result.ID, 10)
			repoImage.MediaType = "application/vnd.docker.distribution.manifest.v2+json"
			repoImage.Repo = repositoryLink

			for _, image := range result.Images {
				totalImageSize += int64(image.Size)
				repoImage.Digest = append(repoImage.Digest, image.Digest)
			}

			repoImage.ImageSizeBytes = strconv.FormatInt(totalImageSize, 10)

			repository.Images = append(repository.Images, repoImage)
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

// NewHubClientParams is the required parameters to build a new client
type NewHubClientParams struct {
	Namespace string
}

// NewHubClient builds a new client
func NewHubClient(params NewHubClientParams) cr.Client {
	return &RegistryClient{
		httpClient: myhttp.NewClient(myhttp.NewClientParams{
			BaseURL:             baseURL,
			InjectAuthInRequest: nil, // we set this once we have logged in
		}),
		namespace: params.Namespace,
	}
}
