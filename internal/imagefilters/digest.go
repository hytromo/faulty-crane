package imagefilters

import (
	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
)

func digestFilter(repos []containerregistry.Repository, digestsToKeep []string) {
	if len(digestsToKeep) == 0 {
		return
	}

	// let's create a map for the digests so we don't do O(n) every time we are searching to see if a digest is whitelisted
	digestsToKeepMap := make(map[string]bool)
	for _, digest := range digestsToKeep {
		digestsToKeepMap[digest] = true
	}

	for repoIndex := range repos {
		for imageIndex := range repos[repoIndex].Images {
			parsedImage := repos[repoIndex].Images[imageIndex]
			if parsedImage.KeptData.Reason != keepreasons.None {
				// image already kept for some other reason
				continue
			}

			_, exists := digestsToKeepMap[parsedImage.Digest]
			if exists {
				repos[repoIndex].Images[imageIndex].KeptData.Reason = keepreasons.WhitelistedDigest
				break
			}
		}
	}
}
