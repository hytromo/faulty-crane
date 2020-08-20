package imagefilters

import "github.com/hytromo/faulty-crane/internal/keepreasons"

func repoFilter(repos []ParsedRepo, reposToKeep []string) {
	if len(reposToKeep) == 0 {
		return
	}

	// let's create a map for the repos so we don't do O(n) every time we are searching to see if a repo is whitelisted
	reposToKeepMap := make(map[string]bool)
	for _, repo := range reposToKeep {
		reposToKeepMap[repo] = true
	}

	for repoIndex := range repos {
		if _, exists := reposToKeepMap[repos[repoIndex].Repo.Link]; exists {
			for imageIndex := range repos[repoIndex].Images {
				parsedImage := repos[repoIndex].Images[imageIndex]
				if parsedImage.KeptData.Reason != keepreasons.None {
					// image already kept for some other reason
					continue
				}

				repos[repoIndex].Images[imageIndex].KeptData.Reason = keepreasons.WhitelistedRepository
			}
		}

	}
}
