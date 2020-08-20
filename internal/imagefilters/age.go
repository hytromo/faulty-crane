package imagefilters

import (
	"fmt"
	"strconv"
	"time"

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

func ageFilter(repos []ParsedRepo, keepYoungerThan string) {
	if keepYoungerThan == "" {
		return
	}

	nowMs := getMsTime()
	youngerDurationMs := getStringDurationInMs(keepYoungerThan)

	fmt.Println("Younger duration ms", youngerDurationMs)

	for repoIndex := range repos {
		for imageIndex := range repos[repoIndex].Images {
			parsedImage := repos[repoIndex].Images[imageIndex]

			if parsedImage.KeptReason != "" {
				// image already kept for some other reason
				continue
			}

			image := parsedImage.Image

			uploadedMs, err := strconv.ParseInt(image.TimeUploadedMs, 10, 64)

			if err != nil {
				log.Fatalf("Image %v contains invalid time uploaded field: %v", image.Digest, image.TimeUploadedMs)
			}

			ageMs := nowMs - uploadedMs

			fmt.Println("Age ms", ageMs)

			if ageMs < youngerDurationMs {
				// image young enough, needs to be kept
				repos[repoIndex].Images[imageIndex].KeptReason = "Young"
			}
		}
	}
}
