package imagefilters

import (
	"github.com/hytromo/faulty-crane/internal/configuration"
	"github.com/hytromo/faulty-crane/internal/containerregistry"
)

// Parse takes all the container images and the filters dictated by the user and applies the filters to the images
func Parse(repos []containerregistry.Repository, keepImages configuration.KeepImages) []containerregistry.Repository {
	parsedRepos := make([]containerregistry.Repository, len(repos))
	copy(parsedRepos, repos)

	repoFilter(parsedRepos, keepImages.Image.Repositories)
	ageFilter(parsedRepos, keepImages.YoungerThan)
	tagFilter(parsedRepos, keepImages.Image.Tags)
	digestFilter(parsedRepos, keepImages.Image.Digests)
	k8sFilter(parsedRepos, keepImages.UsedIn.KubernetesClusters)

	return parsedRepos
}
