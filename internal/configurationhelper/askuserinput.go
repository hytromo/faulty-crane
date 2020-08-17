package configurationhelper

import (
	"fmt"
	"strings"

	"github.com/hytromo/faulty-crane/internal/ask"
	"github.com/hytromo/faulty-crane/internal/utils/stringutil"
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

func isGCR(registryLink string) bool {
	return strings.Contains(registryLink, "gcr.io/") && strings.Contains(registryLink, "/v2/")
}

func askContainerRegistryLink() string {
	for {
		registryLink := ask.Str(ask.Question{
			Description: "Container registry link for cleanup (e.g. https://eu.gcr.io/v2/project-name)",
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
			Description: "Access token (gcloud auth print-access-token)",
		})
	}

	return ""
}

func askYoungerThan() string {
	for {
		youngerThan := ask.Str(ask.Question{
			Description: "Keep images younger than (e.g. 10d3h, empty=ignore age)",
		})

		if youngerThan == "" {
			break
		}

		duration, err := duration.ParseDuration(youngerThan)

		if err == nil && duration.Seconds() > 0 {
			return youngerThan
		}

		fmt.Println("Please give a valid duration")
	}

	return ""
}

func askKubernetesClusters() []string {
	return askListOfStrings("Keep images used in k8s cluster context")
}

func askImageTags() []string {
	return askListOfStrings("Keep images having tag")
}

func askImageDigests() []string {
	return askListOfStrings("Keep images having digest")
}

func askImageIds() []string {
	return askListOfStrings("Keep images having id")
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
	containerRegistryLink := askContainerRegistryLink()

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
