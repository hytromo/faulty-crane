package imagefilters

import (
	"sort"
	"strconv"

	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
	log "github.com/sirupsen/logrus"
)

func numberFilter(repos []containerregistry.Repository, _keepAtLeast int) {
	if _keepAtLeast == 0 {
		return
	}

	for repoIndex, repo := range repos {
		keptInRepoNumber := 0
		for _, parsedImage := range repo.Images {
			if parsedImage.KeptData.Reason != keepreasons.None {
				// image already kept for some other reason
				continue
			}
		}

		keepAtLeast := _keepAtLeast
		repoImagesCount := len(repo.Images)
		if keepAtLeast > repoImagesCount {
			keepAtLeast = repoImagesCount // we cannot keep more than the repo images count
		}

		needToKeepAdditionalToReachAtLeast := keepAtLeast - keptInRepoNumber

		if needToKeepAdditionalToReachAtLeast <= 0 {
			// this repo already has enough images, so we can move on to the next repo
			continue
		}

		// largest age (= more recent) first
		sort.SliceStable(repo.Images, func(i, j int) bool {
			imageI := repo.Images[i]
			imageJ := repo.Images[j]
			uploadedMsI, err := strconv.ParseInt(imageI.TimeUploadedMs, 10, 64)

			if err != nil {
				log.Fatalf("Image %v contains invalid time uploaded field: %v", imageI.Digest, imageI.TimeUploadedMs)
			}

			uploadedMsJ, err := strconv.ParseInt(imageJ.TimeUploadedMs, 10, 64)

			if err != nil {
				log.Fatalf("Image %v contains invalid time uploaded field: %v", imageJ.Digest, imageJ.TimeUploadedMs)
			}

			return uploadedMsI > uploadedMsJ
		})

		markedAsKeptNumber := 0
		for imageIndex, image := range repo.Images {
			if image.KeptData.Reason == keepreasons.None {
				repos[repoIndex].Images[imageIndex].KeptData.Reason = keepreasons.OneOfFew
				markedAsKeptNumber++
				if markedAsKeptNumber >= needToKeepAdditionalToReachAtLeast {
					break // move on to next repo
				}
			}
		}
	}
}
