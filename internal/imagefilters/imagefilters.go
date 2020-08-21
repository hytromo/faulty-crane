package imagefilters

import (
	"github.com/hytromo/faulty-crane/internal/configuration"
	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
)

// ParsedImage is a container image that has gone through filters and has metadata about the results of the filtering process
type ParsedImage struct {
	Image    containerregistry.ContainerImage
	KeptData keepreasons.KeptData
}

// ParsedRepo is a repository that contains parsed/filtered images
type ParsedRepo struct {
	Repo   containerregistry.Repository
	Images []ParsedImage
}

func (parsedImage ParsedImage) parsedImageWithKeptReason(reason keepreasons.KeptReason) ParsedImage {
	parsedImage.KeptData.Reason = reason
	parsedImage.KeptData.Metadata = ""
	return parsedImage
}

func reposToParsedRepos(repos []containerregistry.Repository) []ParsedRepo {
	parsedRepos := make([]ParsedRepo, len(repos))

	for i, repo := range repos {
		parsedRepos[i] = ParsedRepo{
			Repo:   repo,
			Images: imagesToParsedImages(repo.Images),
		}
	}

	return parsedRepos
}

func imagesToParsedImages(images []containerregistry.ContainerImage) []ParsedImage {
	parsedImages := make([]ParsedImage, len(images))
	for i, image := range images {
		parsedImages[i] = ParsedImage{
			Image: image,
			KeptData: keepreasons.KeptData{
				Reason:   keepreasons.None,
				Metadata: "",
			},
		}
	}
	return parsedImages
}

// Parse takes all the container images and the filters dictated by the user and applies the filters to the images
func Parse(repos []containerregistry.Repository, keepImages configuration.KeepImages) []ParsedRepo {
	parsedRepos := reposToParsedRepos(repos)

	repoFilter(parsedRepos, keepImages.Image.Repositories)
	ageFilter(parsedRepos, keepImages.YoungerThan)
	tagFilter(parsedRepos, keepImages.Image.Tags)
	digestFilter(parsedRepos, keepImages.Image.Digests)
	k8sFilter(parsedRepos, keepImages.UsedIn.KubernetesClusters)

	return parsedRepos
}
