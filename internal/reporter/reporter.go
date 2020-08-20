package reporter

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hytromo/faulty-crane/internal/imagefilters"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
	"github.com/hytromo/faulty-crane/internal/utils/stringutil"
	color "github.com/logrusorgru/aurora"
	tablewriter "github.com/olekukonko/tablewriter"
)

// ReportRepositoriesStatus prints out in a nice way the status of the repositories, e.g. what needs to be deleted and for what reason
func ReportRepositoriesStatus(repos []imagefilters.ParsedRepo) {
	sort.SliceStable(repos, func(i int, j int) bool {
		return repos[i].Repo.Link < repos[j].Repo.Link
	})

	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"Kept", "Repo", "Tags", "Digest", "Size", "Uploaded"}
	headersCount := len(headers)
	table.SetHeader(headers)

	keepCount := 0
	deleteCount := 0
	var deleteTotalSizeBytes int64 = 0
	var keepTotalSizeBytes int64 = 0

	for _, parsedRepo := range repos {
		for _, parsedImage := range parsedRepo.Images {
			image := parsedImage.Image
			keptReason := parsedImage.KeptData.Reason

			tableValues := make([]string, headersCount)
			tableColors := make([]tablewriter.Colors, headersCount)

			imageSizeBytes, err := strconv.ParseInt(image.ImageSizeBytes, 10, 64)
			if err != nil {
				imageSizeBytes = 0 // we will not crash the app for this reason
			}

			if keptReason == keepreasons.None {
				// needs to be deleted
				tableValues[0] = "✗"
				tableColors[0] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor}
				deleteCount++
				deleteTotalSizeBytes = deleteTotalSizeBytes + imageSizeBytes
			} else {
				tableValues[0] = "✔"
				tableColors[0] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor}
				keepCount++
				keepTotalSizeBytes = keepTotalSizeBytes + imageSizeBytes
			}

			tableValues[1] = stringutil.KeepAtMost(parsedRepo.Repo.Link, 80)
			if keptReason == keepreasons.WhitelistedRepository {
				tableColors[1] = tablewriter.Colors{tablewriter.FgGreenColor}
			} else {
				tableColors[1] = tablewriter.Colors{}
			}

			tableValues[2] = stringutil.KeepAtMost(strings.Join(image.Tag, ","), 50)
			if keptReason == keepreasons.WhitelistedTag {
				tableColors[2] = tablewriter.Colors{tablewriter.FgGreenColor}
			} else {
				tableColors[2] = tablewriter.Colors{}
			}

			digestClean := strings.Replace(image.Digest, "sha256:", "", 1)
			tableValues[3] = stringutil.TrimRightChars(digestClean, len(digestClean)-12) // keep only the first few chars
			if keptReason == keepreasons.WhitelistedDigest {
				tableColors[3] = tablewriter.Colors{tablewriter.FgGreenColor}
			} else {
				tableColors[3] = tablewriter.Colors{}
			}

			tableColors[4] = tablewriter.Colors{}

			tableValues[4] = stringutil.HumanFriendlySize(imageSizeBytes)

			uploadedMs, err := strconv.ParseInt(image.TimeUploadedMs, 10, 64)
			if err != nil {
				log.Fatalf("Invalid uploaded timestamp %v", image.TimeUploadedMs)
			}

			tableValues[5] = time.Unix(uploadedMs/1000, 0).Format(time.RFC822)
			if keptReason == keepreasons.Young {
				tableColors[5] = tablewriter.Colors{tablewriter.FgGreenColor}
			} else {
				tableColors[5] = tablewriter.Colors{}
			}

			table.Rich(tableValues, tableColors)
		}
	}

	table.Render()

	fmt.Println(
		deleteCount,
		"image(s) will be deleted,",
		color.Red(fmt.Sprintf("or %v", stringutil.HumanFriendlySize(deleteTotalSizeBytes))),
	)

	fmt.Println(
		keepCount,
		"image(s) will be kept,",
		color.Green(fmt.Sprintf("or %v", stringutil.HumanFriendlySize(keepTotalSizeBytes))),
	)
}
