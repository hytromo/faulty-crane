package reporter

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	timeago "github.com/caarlos0/timea.go"
	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
	"github.com/hytromo/faulty-crane/internal/utils/stringutil"
	color "github.com/logrusorgru/aurora"
	tablewriter "github.com/olekukonko/tablewriter"
)

// ReportRepositoriesStatus prints out in a nice way the status of the repositories, e.g. what needs to be deleted and for what reason
func ReportRepositoriesStatus(repos []containerregistry.Repository, showAnalyticalPlan bool) {
	sort.SliceStable(repos, func(i int, j int) bool {
		return repos[i].Link < repos[j].Link
	})

	table := tablewriter.NewWriter(os.Stdout)

	keepCount := 0
	deleteCount := 0
	var deleteTotalSizeBytes int64 = 0
	var keepTotalSizeBytes int64 = 0

	if showAnalyticalPlan {
		headers := []string{"Kept", "Repo", "Tags", "Digest", "Size", "Cluster", "Uploaded"}
		headersCount := len(headers)
		table.SetHeader(headers)
		for _, parsedRepo := range repos {
			for _, parsedImage := range parsedRepo.Images {
				image := parsedImage
				keptReason := parsedImage.KeptData.Reason

				tableValues := make([]string, headersCount)
				tableColors := make([]tablewriter.Colors, headersCount)

				imageSizeBytes, err := strconv.ParseInt(image.ImageSizeBytes, 10, 64)
				if err != nil {
					imageSizeBytes = 0 // we will not crash the app for this reason
				}

				if keptReason == keepreasons.None {
					// needs to be deleted
					tableValues[0] = "✗ NO"
					tableColors[0] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor}
					deleteCount++
					deleteTotalSizeBytes = deleteTotalSizeBytes + imageSizeBytes
				} else {
					tableValues[0] = "✔ YES"
					tableColors[0] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor}
					keepCount++
					keepTotalSizeBytes = keepTotalSizeBytes + imageSizeBytes
				}

				tableValues[1] = stringutil.KeepAtMost(parsedRepo.Link, 80)
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

				tableValues[5] = "-"
				if keptReason == keepreasons.UsedInCluster {
					tableValues[5] = parsedImage.KeptData.Metadata
					tableColors[5] = tablewriter.Colors{tablewriter.FgGreenColor}
				} else {
					tableColors[5] = tablewriter.Colors{}
				}

				tableValues[6] = time.Unix(uploadedMs/1000, 0).Format(time.RFC822)
				if keptReason == keepreasons.Young {
					tableColors[6] = tablewriter.Colors{tablewriter.FgGreenColor}
				} else {
					tableColors[6] = tablewriter.Colors{}
				}

				table.Rich(tableValues, tableColors)
			}
		}
	} else {
		headers := []string{"repo", "deleted", "deleted size", "most recent to be deleted"}
		headersCount := len(headers)
		table.SetHeader(headers)
		for _, parsedRepo := range repos {
			tableValues := make([]string, headersCount)
			tableColors := make([]tablewriter.Colors, headersCount)

			tableValues[0] = parsedRepo.Link
			totalImagesCountInRepo := len(parsedRepo.Images)
			deletedImagesCountInRepo := 0
			var deleteTotalSizeInRepoBytes int64 = 0
			var keepTotalSizeBytesInRepo int64 = 0
			var latestUploadedTimeStampToBeDeleted int64 = 0

			for _, parsedImage := range parsedRepo.Images {
				image := parsedImage
				keptReason := parsedImage.KeptData.Reason

				imageSizeBytes, err := strconv.ParseInt(image.ImageSizeBytes, 10, 64)
				if err != nil {
					imageSizeBytes = 0 // we will not crash the app for this reason
				}

				if keptReason == keepreasons.None {
					// needs to be deleted
					deletedImagesCountInRepo++
					deleteTotalSizeInRepoBytes = deleteTotalSizeInRepoBytes + imageSizeBytes
					uploadedMs, err := strconv.ParseInt(image.TimeUploadedMs, 10, 64)
					if err == nil {
						if uploadedMs > latestUploadedTimeStampToBeDeleted {
							latestUploadedTimeStampToBeDeleted = uploadedMs
						}
					}

				} else {
					keepTotalSizeBytesInRepo = keepTotalSizeBytesInRepo + imageSizeBytes
				}

			}

			deleteCount = deleteCount + deletedImagesCountInRepo
			keepCount = keepCount + (totalImagesCountInRepo - deletedImagesCountInRepo)
			deleteTotalSizeBytes = deleteTotalSizeBytes + deleteTotalSizeInRepoBytes
			keepTotalSizeBytes = keepTotalSizeBytes + keepTotalSizeBytesInRepo

			tableValues[1] = fmt.Sprintf("%6.2f%% / %v/%v", float64(deletedImagesCountInRepo)/float64(totalImagesCountInRepo)*100, deletedImagesCountInRepo, totalImagesCountInRepo)

			colorToPaint := tablewriter.Colors{tablewriter.Normal}
			if deletedImagesCountInRepo == totalImagesCountInRepo {
				colorToPaint = tablewriter.Colors{tablewriter.FgRedColor}
			} else if deletedImagesCountInRepo > 0 {
				colorToPaint = tablewriter.Colors{tablewriter.FgYellowColor}
			} else if deletedImagesCountInRepo == 0 {
				colorToPaint = tablewriter.Colors{tablewriter.FgGreenColor}
			}

			for i, _ := range tableColors {
				tableColors[i] = colorToPaint
			}

			tableValues[2] = fmt.Sprintf("%v/%v", stringutil.HumanFriendlySize(deleteTotalSizeInRepoBytes), stringutil.HumanFriendlySize(deleteTotalSizeInRepoBytes+keepTotalSizeBytesInRepo))

			if latestUploadedTimeStampToBeDeleted == 0 {
				tableValues[3] = "-"
			} else {
				tableValues[3] = timeago.Of(time.Unix(latestUploadedTimeStampToBeDeleted/1000, 0).UTC())
			}

			table.Rich(tableValues, tableColors)
		}
	}

	table.Render()

	totalBytes := deleteTotalSizeBytes + keepTotalSizeBytes
	totalImages := deleteCount + keepCount

	fmt.Println(
		deleteCount,
		"image(s) will be deleted",
		color.Red(
			fmt.Sprintf(
				"/ %v / %.2f%% of total images / %.2f%% of total size",
				stringutil.HumanFriendlySize(deleteTotalSizeBytes),
				float64(deleteCount)/float64(totalImages)*100,
				float64(deleteTotalSizeBytes)/float64(totalBytes)*100,
			),
		),
	)

	fmt.Println(
		keepCount,
		"image(s) will be kept",
		color.Green(fmt.Sprintf("/ %v", stringutil.HumanFriendlySize(keepTotalSizeBytes))),
	)
}
