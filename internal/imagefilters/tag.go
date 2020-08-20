package imagefilters

import "github.com/hytromo/faulty-crane/internal/keepreasons"

func tagFilter(repos []ParsedRepo, tagsToKeep []string) {
	if len(tagsToKeep) == 0 {
		return
	}
	// let's create a map for the tags so we don't do O(n) every time we are searching to see if a tag is whitelisted
	tagsToKeepMap := make(map[string]bool)

	for _, tag := range tagsToKeep {
		tagsToKeepMap[tag] = true
	}

	for repoIndex := range repos {
		for imageIndex := range repos[repoIndex].Images {
			parsedImage := repos[repoIndex].Images[imageIndex]

			if parsedImage.KeptData.Reason != keepreasons.None {
				// image already kept for some other reason
				continue
			}

			image := parsedImage.Image
			for _, tag := range image.Tag {
				_, exists := tagsToKeepMap[tag]
				if exists {
					repos[repoIndex].Images[imageIndex].KeptData.Reason = keepreasons.WhitelistedTag
					break
				}
			}
		}
	}
}
