package imagefilters

func digestFilter(repos []ParsedRepo, digestsToKeep []string) {
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
			if parsedImage.KeptReason != "" {
				// image already kept for some other reason
				continue
			}

			_, exists := digestsToKeepMap[parsedImage.Image.Digest]
			if exists {
				repos[repoIndex].Images[imageIndex].KeptReason = "Whitelisted digest"
				break
			}
		}
	}
}
