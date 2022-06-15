package configuration

import (
	"fmt"
	"io"

	"github.com/hytromo/faulty-crane/internal/ask"
	"github.com/hytromo/faulty-crane/internal/utils/stringutil"
	color "github.com/logrusorgru/aurora"
	"maze.io/x/duration"
)

func askListOfStrings(readDevice io.Reader, question string) []string {
	emptyToFinishSuffix := "(empty=done)"

	var answers = []string{}

	for {
		answer := ask.Str(ask.Question{
			Description: fmt.Sprintf(question+" #%v %v", len(answers)+1, emptyToFinishSuffix),
			ReadDevice:  readDevice,
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

func askContainerRegistryNamespace(readDevice io.Reader) string {
	for {
		namespace := ask.Str(ask.Question{
			Description: "Namespace (e.g. organization name or your own username)",
			ReadDevice:  readDevice,
		})

		if namespace == "" {
			fmt.Println("A namespace is required")
		} else {
			return namespace
		}
	}
}

func askContainerRegistryUsername(readDevice io.Reader) string {
	for {
		username := ask.Str(ask.Question{
			Description: "Username",
			ReadDevice:  readDevice,
		})

		if username == "" {
			fmt.Println("A username is required")
		} else {
			return username
		}
	}
}

func askContainerType(readDevice io.Reader) string {
	return ask.Str(ask.Question{
		Description:     fmt.Sprintf("Container %v", color.Green("registry type")),
		PossibleAnswers: []string{"gcr", "dockerhub"},
		ReadDevice:      readDevice,
	})
}

func askContainerRegistryLink(readDevice io.Reader) string {
	return ask.Str(ask.Question{
		Description: fmt.Sprintf("Container %v for cleanup", color.Green("registry link")),
		PossibleAnswers: []string{
			"gcr.io", "eu.gcr.io", "us.gcr.io", "asia.gcr.io",
		},
		ReadDevice: readDevice,
	})
}

func askContainerRegistryPassword(readDevice io.Reader, containerRegistryLink string) string {
	return ask.Str(ask.Question{
		Description: fmt.Sprintf("%v (e.g. `gcloud auth print-access-token` for gcr, password for dockerhub)", color.Green("Access token")),
		ReadDevice:  readDevice,
	})
}

func askYoungerThan(readDevice io.Reader) string {
	for {
		youngerThan := ask.Str(ask.Question{
			Description: fmt.Sprintf("Keep images %v (e.g. 10d3h, empty=ignore age)", color.Green("younger than")),
			ReadDevice:  readDevice,
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

func askKubernetesClusters(readDevice io.Reader) []string {
	return askListOfStrings(readDevice, fmt.Sprintf("Keep images used in %v", color.Green("k8s cluster context")))
}

func askImageTags(readDevice io.Reader) []string {
	return askListOfStrings(readDevice, fmt.Sprintf("Keep images having %v", color.Green("tag")))
}

func askImageDigests(readDevice io.Reader) []string {
	return askListOfStrings(readDevice, fmt.Sprintf("Keep images having %v", color.Green("digest")))
}

func askImageIds(readDevice io.Reader) []string {
	return askListOfStrings(readDevice, fmt.Sprintf("Keep images having %v", color.Green("id")))
}

// UserInput is a struct holding the user's answers
type UserInput struct {
	ContainerRegistryLink      string
	ContainerRegistryUsername  string
	ContainerRegistryPassword  string
	ContainerRegistryNamespace string
	YoungerThan                string
	KubernetesClusters         []string
	ImageTags                  []string
	ImageDigests               []string
	ImageIDs                   []string
}

// AskUserInput asks for user input in order to create a new configuration
func AskUserInput(readDevice io.Reader) UserInput {
	containerType := askContainerType(readDevice)
	containerRegistryLink := ""
	containerRegistryUsername := ""
	containerRegistryNamespace := ""

	if containerType == "gcr" {
		containerRegistryLink = askContainerRegistryLink(readDevice)
	} else if containerType == "dockerhub" {
		containerRegistryUsername = askContainerRegistryUsername(readDevice)
		containerRegistryNamespace = askContainerRegistryNamespace(readDevice)

	}

	return UserInput{
		ContainerRegistryLink:      containerRegistryLink,
		ContainerRegistryPassword:  askContainerRegistryPassword(readDevice, containerRegistryLink),
		ContainerRegistryUsername:  containerRegistryUsername,
		ContainerRegistryNamespace: containerRegistryNamespace,
		YoungerThan:                askYoungerThan(readDevice),
		KubernetesClusters:         askKubernetesClusters(readDevice),
		ImageTags:                  askImageTags(readDevice),
		ImageDigests:               askImageDigests(readDevice),
		ImageIDs:                   askImageIds(readDevice),
	}
}
