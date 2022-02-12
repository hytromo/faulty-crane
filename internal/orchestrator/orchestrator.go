package orchestrator

import (
	"math"

	"github.com/cheggaaa/pb/v3"
	"github.com/hytromo/faulty-crane/internal/configuration"
	cr "github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/containerregistry/gcr"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
	log "github.com/sirupsen/logrus"
)

// Orchestrator is the entry-point of the commands and handles high level things like implementing goroutines and showing terminal output
type Orchestrator struct {
	crClient cr.ContainerRegistryClient
}

func NewOrchestrator(options configuration.AppOptions) Orchestrator {
	return Orchestrator{
		crClient: gcr.NewGCRClient(gcr.NewGCRClientParams{
			Hostname: options.ApplyPlanCommon.GoogleContainerRegistry.Host,
			Token:    options.ApplyPlanCommon.GoogleContainerRegistry.Token,
		}),
	}
}

func getNeedingDeletionInRepoCount(repo cr.Repository) int {
	repoImagesToDelete := 0
	for _, image := range repo.Images {
		if image.KeptData.Reason == keepreasons.None {
			repoImagesToDelete++
		}
	}
	return repoImagesToDelete
}

func (orchestrator Orchestrator) deleteImageFromChan(repoLink string, imagesToDeleteChan chan cr.ContainerImage, imagesDeletedChan chan error) {
	for image := range imagesToDeleteChan {
		imagesDeletedChan <- orchestrator.crClient.DeleteImage(repoLink, image, false)
	}
}

func (orchestrator Orchestrator) deleteRepoImages(repo cr.Repository, pb *pb.ProgressBar) cr.RepoDeletionResult {
	result := cr.RepoDeletionResult{
		ShouldDeleteCount:    0,
		ManagedToDeleteCount: 0,
	}

	result.ShouldDeleteCount = getNeedingDeletionInRepoCount(repo)
	deletingImagesWorkersNum := int(math.Min(8, float64(result.ShouldDeleteCount)))

	imagesToDeleteChan := make(chan cr.ContainerImage, result.ShouldDeleteCount) // jobs
	imagesDeletedChan := make(chan error, result.ShouldDeleteCount)              // results

	for i := 1; i <= deletingImagesWorkersNum; i++ {
		go orchestrator.deleteImageFromChan(repo.Link, imagesToDeleteChan, imagesDeletedChan)
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

			if managedToDeleteImage != nil {
				result.ManagedToDeleteCount++
			}
		}
	}

	return result
}

func (orchestrator *Orchestrator) deleteRepoImagesWorker(repos <-chan cr.Repository, deletionResults chan<- cr.RepoDeletionResult, pb *pb.ProgressBar) {
	for repo := range repos {
		deletionResults <- orchestrator.deleteRepoImages(repo, pb)
	}
}

func (orchestrator *Orchestrator) DeleteImagesWithNoKeepReason(repos []cr.Repository) cr.RepoDeletionResult {
	allResults := cr.RepoDeletionResult{
		ShouldDeleteCount:    0,
		ManagedToDeleteCount: 0,
	}

	reposCount := len(repos)

	// spawn max 40 goroutines, if repos are less than 40, try to list them all concurrently
	reposDeletingWorkersNum := int(math.Min(8, float64(reposCount)))

	repositoryLinksChan := make(chan cr.Repository, reposCount)         // jobs
	deletionResultsChan := make(chan cr.RepoDeletionResult, reposCount) // results

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
		go orchestrator.deleteRepoImagesWorker(repositoryLinksChan, deletionResultsChan, bar)
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

func (orchestrator Orchestrator) fetchRepoImagesWorker(repositoryLinks <-chan string, parsedRepos chan<- cr.Repository) {
	for repo := range repositoryLinks {
		parsedRepos <- orchestrator.crClient.ParseRepo(repo)
	}
}

func (orchestrator Orchestrator) GetAllRepos() []cr.Repository {
	log.Info("Getting all the repos of the registry...")

	repos := orchestrator.crClient.GetAllRepos()
	reposCount := len(repos)

	bar := pb.Full.Start(reposCount)

	repositoryLinksChan := make(chan string, reposCount)    // jobs
	parsedReposChan := make(chan cr.Repository, reposCount) // results

	// spawn max 40 goroutines, if repos are less than 40, try to list them all concurrently
	workersNum := int(math.Min(40, float64(reposCount)))

	log.Info("Fetching the images of ", reposCount, " repo(s), using ", workersNum, " routine(s)")

	for i := 1; i <= workersNum; i++ {
		go orchestrator.fetchRepoImagesWorker(repositoryLinksChan, parsedReposChan)
	}

	for _, repo := range repos {
		// feed the jobs to the workers
		repositoryLinksChan <- repo
	}

	// while the jobs are being done by the workers, we are merging all the results into one
	allRepos := []cr.Repository{}
	for range repos {
		allRepos = append(allRepos, <-parsedReposChan)
		bar.Increment()
	}

	bar.Finish()

	return allRepos
}
