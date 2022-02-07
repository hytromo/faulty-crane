package containerregistry

import (
	"math"
	"net/http"

	"github.com/cheggaaa/pb/v3"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
	log "github.com/sirupsen/logrus"
)

// GCRClient is used for creating Google Container Registry clients
type GCRClient struct {
	Host      string
	AccessKey string
	client    *http.Client
}

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
	Digest         string
	Repo           string               // Repo is the name of the image's repository without the tag in the form e.g. eu.gcr.io/faulty-crane-project/faulty-crane-test
	KeptData       keepreasons.KeptData `json:",omitempty"`
}

type RepoDeletionResult struct {
	ShouldDeleteCount    int
	ManagedToDeleteCount int
}

// MakeGCRClient builds a new GCRClient instance, adding the missing default values e.g. http client
func MakeGCRClient(client GCRClient) GCRClient {
	gcrClient := client
	gcrClient.client = &http.Client{}
	return gcrClient
}

func (gcrClient GCRClient) deleteImageFromChan(repoLink string, imagesToDeleteChan chan ContainerImage, imagesDeletedChan chan bool) {
	for image := range imagesToDeleteChan {
		imagesDeletedChan <- gcrClient.DeleteImage(repoLink, image, false)
	}
}

func (gcrClient GCRClient) deleteRepoImages(repo Repository, pb *pb.ProgressBar) RepoDeletionResult {
	result := RepoDeletionResult{
		ShouldDeleteCount:    0,
		ManagedToDeleteCount: 0,
	}

	result.ShouldDeleteCount = getNeedingDeletionInRepoCount(repo)
	deletingImagesWorkersNum := int(math.Min(8, float64(result.ShouldDeleteCount)))

	imagesToDeleteChan := make(chan ContainerImage, result.ShouldDeleteCount) // jobs
	imagesDeletedChan := make(chan bool, result.ShouldDeleteCount)            // results

	for i := 1; i <= deletingImagesWorkersNum; i++ {
		go gcrClient.deleteImageFromChan(repo.Link, imagesToDeleteChan, imagesDeletedChan)
	}

	for _, image := range repo.Images {
		if image.KeptData.Reason == keepreasons.None {
			// feed the jobs to the workers
			imagesToDeleteChan <- image
		}
	}

	// while the jobs are being done by the workers, we are counting them
	for _, image := range repo.Images {
		if image.KeptData.Reason == keepreasons.None {
			managedToDeleteImage := <-imagesDeletedChan

			pb.Increment()

			if managedToDeleteImage {
				result.ManagedToDeleteCount++
			}
		}
	}

	return result
}

func (gcrClient GCRClient) fetchRepoImagesWorker(repositoryLinks <-chan string, parsedRepos chan<- Repository) {
	for repo := range repositoryLinks {
		parsedRepos <- gcrClient.listTags(repo)
	}
}

func (gcrClient GCRClient) deleteRepoImagesWorker(repos <-chan Repository, deletionResults chan<- RepoDeletionResult, pb *pb.ProgressBar) {
	for repo := range repos {
		deletionResults <- gcrClient.deleteRepoImages(repo, pb)
	}
}

func getNeedingDeletionInRepoCount(repo Repository) int {
	repoImagesToDelete := 0
	for _, image := range repo.Images {
		if image.KeptData.Reason == keepreasons.None {
			repoImagesToDelete++
		}
	}
	return repoImagesToDelete
}

func (gcrClient GCRClient) DeleteImagesWithNoKeepReason(repos []Repository) RepoDeletionResult {
	allResults := RepoDeletionResult{
		ShouldDeleteCount:    0,
		ManagedToDeleteCount: 0,
	}

	reposCount := len(repos)

	// spawn max 40 goroutines, if repos are less than 40, try to list them all concurrently
	reposDeletingWorkersNum := int(math.Min(8, float64(reposCount)))

	repositoryLinksChan := make(chan Repository, reposCount)         // jobs
	deletionResultsChan := make(chan RepoDeletionResult, reposCount) // results

	totalImagesToDelete := 0

	for _, repo := range repos {
		totalImagesToDelete += getNeedingDeletionInRepoCount(repo)
	}

	if totalImagesToDelete == 0 {
		return allResults
	}

	log.Info("Deleting the images of ", reposCount, " repo(s), using ", reposDeletingWorkersNum, " routine(s)")

	bar := pb.Full.Start(totalImagesToDelete)

	for i := 1; i <= reposDeletingWorkersNum; i++ {
		go gcrClient.deleteRepoImagesWorker(repositoryLinksChan, deletionResultsChan, bar)
	}

	for _, repo := range repos {
		// feed the jobs to the workers
		repositoryLinksChan <- repo
	}

	// while the jobs are being done by the workers, we are merging all the results into one
	for range repos {
		thisResult := <-deletionResultsChan
		allResults.ShouldDeleteCount += thisResult.ShouldDeleteCount
		allResults.ManagedToDeleteCount += thisResult.ManagedToDeleteCount
	}

	bar.Finish()

	return allResults
}

// GetAllRepos finds all the repositories and parses them finding the images they contain
func (gcrClient GCRClient) GetAllRepos() []Repository {
	log.Info("Getting all the repos of the registry...")

	repos := gcrClient.getRepositories()
	reposCount := len(repos)

	bar := pb.Full.Start(reposCount)

	repositoryLinksChan := make(chan string, reposCount) // jobs
	parsedReposChan := make(chan Repository, reposCount) // results

	// spawn max 40 goroutines, if repos are less than 40, try to list them all concurrently
	workersNum := int(math.Min(40, float64(reposCount)))

	log.Info("Fetching the images of ", reposCount, " repo(s), using ", workersNum, " routine(s)")

	for i := 1; i <= workersNum; i++ {
		go gcrClient.fetchRepoImagesWorker(repositoryLinksChan, parsedReposChan)
	}

	for _, repo := range repos {
		// feed the jobs to the workers
		repositoryLinksChan <- repo
	}

	// while the jobs are being done by the workers, we are merging all the results into one
	allRepos := []Repository{}
	for range repos {
		allRepos = append(allRepos, <-parsedReposChan)
		bar.Increment()
	}

	bar.Finish()

	return allRepos
}
