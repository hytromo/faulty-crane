package configuration

import (
	"fmt"

	"github.com/hytromo/faulty-crane/internal/ask"
	"github.com/hytromo/faulty-crane/internal/utils/stringutil"
	color "github.com/logrusorgru/aurora"
	"maze.io/x/duration"
)

func askListOfStrings(question string) []string {
	emptyToFinishSuffix := "(empty=done)"

	var answers = []string{}

	for {
		answer := ask.Str(ask.Question{
			Description: fmt.Sprintf(question+" #%v %v", len(answers)+1, emptyToFinishSuffix),
		})

		if answer == "" {
			break
		}

		if !stringutil.StrInSlice(answer, answers) {
			answers = append(answers, answer)
		}
	}

	return answers
}

var gcrRegistryHosts = []string{
	"gcr.io", "eu.gcr.io", "us.gcr.io", "asia.gcr.io",
}

func isGCR(registryLink string) bool {
	return registryLink == "gcr.io" || registryLink == "eu.gcr.io" || registryLink == "us.gcr.io" || registryLink == "asia.gcr.io"
}

func askContainerRegistryHost() string {
	for {
		registryLink := ask.Str(ask.Question{
			Description:     fmt.Sprintf("Container %v for cleanup", color.Green("registry link")),
			PossibleAnswers: gcrRegistryHosts,
		})

		if !isGCR(registryLink) {
			fmt.Println("Only Google Container Registry (GCR) is supported for now, please try again")
		} else {
			return registryLink
		}
	}
}

func askContainerRegistryKey(containerRegistryLink string) string {
	if isGCR(containerRegistryLink) {
		return ask.Str(ask.Question{
			Description: fmt.Sprintf("%v (gcloud auth print-access-token)", color.Green("Access token")),
		})
	}

	return ""
}

func askYoungerThan() string {
	for {
		youngerThan := ask.Str(ask.Question{
			Description: fmt.Sprintf("Keep images %v (e.g. 10d3h, empty=ignore age)", color.Green("younger than")),
		})

		if youngerThan == "" {
			break
		}

		youngerDuration, err := duration.ParseDuration(youngerThan)

		if err == nil && youngerDuration.Seconds() > 0 {
			return youngerThan
		}

		fmt.Println("Please give a valid duration")
	}

	return ""
}

func askKubernetesClusters() []string {
	return askListOfStrings(fmt.Sprintf("Keep images used in %v", color.Green("k8s cluster context")))
}

func askImageTags() []string {
	return askListOfStrings(fmt.Sprintf("Keep images having %v", color.Green("tag")))
}

func askImageDigests() []string {
	return askListOfStrings(fmt.Sprintf("Keep images having %v", color.Green("digest")))
}

func askImageIds() []string {
	return askListOfStrings(fmt.Sprintf("Keep images having %v", color.Green("id")))
}

// UserInput is a struct holding the user's answers
type UserInput struct {
	containerRegistryLink   string
	containerRegistryAccess string
	youngerThan             string
	kubernetesClusters      []string
	imageTags               []string
	imageDigests            []string
	imageIDs                []string
}

// AskUserInput asks for user input in order to create a new configuration
func AskUserInput() UserInput {
	containerRegistryLink := askContainerRegistryHost()

	return UserInput{
		containerRegistryLink:   containerRegistryLink,
		containerRegistryAccess: askContainerRegistryKey(containerRegistryLink),
		youngerThan:             askYoungerThan(),
		kubernetesClusters:      askKubernetesClusters(),
		imageTags:               askImageTags(),
		imageDigests:            askImageDigests(),
		imageIDs:                askImageIds(),
	}
}
