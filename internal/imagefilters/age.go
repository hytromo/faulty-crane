package imagefilters

import (
	"strconv"
	"time"

	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
	log "github.com/sirupsen/logrus"
	"maze.io/x/duration"
)

func getMsTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func getStringDurationInMs(stringDuration string) int64 {
	parsedDuration, err := duration.ParseDuration(stringDuration)

	if err != nil {
		log.Fatalf("Could not parse duration '%v'. Please check your configuration.", stringDuration)
	}

	return int64(parsedDuration.Seconds() * 1000)
}

func ageFilter(repos []containerregistry.Repository, keepYoungerThan string) {
	if keepYoungerThan == "" {
		return
	}

	nowMs := getMsTime()
	youngerDurationMs := getStringDurationInMs(keepYoungerThan)

	for repoIndex := range repos {
		for imageIndex := range repos[repoIndex].Images {
			parsedImage := repos[repoIndex].Images[imageIndex]

			if parsedImage.KeptData.Reason != keepreasons.None {
				// image already kept for some other reason
				continue
			}

			uploadedMs, err := strconv.ParseInt(parsedImage.TimeUploadedMs, 10, 64)

			if err != nil {
				log.Errorf("Image %v contains invalid time uploaded field: %v", parsedImage.Digest, parsedImage.TimeUploadedMs)
				continue
			}

			ageMs := nowMs - uploadedMs

			if ageMs < youngerDurationMs {
				// image young enough, needs to be kept
				repos[repoIndex].Images[imageIndex].KeptData.Reason = keepreasons.Young
			}
		}
	}
}
