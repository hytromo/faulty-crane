package reporter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hytromo/faulty-crane/internal/imagefilters"
	color "github.com/logrusorgru/aurora"
)

// ReportRepositoriesStatus prints out in a nice way the status of the repositories, e.g. what needs to be deleted and for what reason
func ReportRepositoriesStatus(repos []imagefilters.ParsedRepo) {
	sort.SliceStable(repos, func(i int, j int) bool {
		return repos[i].Repo.Link < repos[j].Repo.Link
	})

	for _, parsedRepo := range repos {
		fmt.Printf("---- %v ----\n", color.Yellow(parsedRepo.Repo.Link))
		for _, parsedImage := range parsedRepo.Images {
			image := parsedImage.Image
			if parsedImage.KeptReason == "" {
				// needs to be deleted
				myStr := fmt.Sprintf("%v / %v", image.Digest, strings.Join(image.Tag, ","))
				fmt.Println(color.Red(myStr))
			} else {
				// needs to be kept

				myStr := fmt.Sprintf("%v / %v -> %v", image.Digest, strings.Join(image.Tag, ","), parsedImage.KeptReason)
				fmt.Println(color.Green(myStr))
			}
		}
		fmt.Println()
	}
}
