package configurationhelper

import (
	"fmt"

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

func askContainerRegistry() string {
	return ask.Str(ask.Question{
		Description:     "Container registry for cleanup",
		PossibleAnswers: []string{"gcr", "azurer"},
		DefaultValue:    "gcr",
	})
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

func askImageRepositories() []string {
	return askListOfStrings("Keep images having repository")
}

// UserInput is a struct holding the user's answers
type UserInput struct {
	containerRegistry  string
	youngerThan        string
	kubernetesClusters []string
	imageTags          []string
	imageDigests       []string
	imageIDs           []string
	imageRepositories  []string
}

// AskUserInput asks for user input in order to create a new configuration
func AskUserInput() UserInput {
	return UserInput{
		containerRegistry:  askContainerRegistry(),
		youngerThan:        askYoungerThan(),
		kubernetesClusters: askKubernetesClusters(),
		imageTags:          askImageTags(),
		imageDigests:       askImageDigests(),
		imageIDs:           askImageIds(),
		imageRepositories:  askImageRepositories(),
	}
}
