package configurationhelper

import (
	"fmt"

	"github.com/hytromo/faulty-crane/internal/argsparser"
	"github.com/hytromo/faulty-crane/internal/ask"
	"github.com/hytromo/faulty-crane/internal/utils/stringutil"
)

// CreateNew asks the user for configuration input and then creates a configuration file based on the answers
func CreateNew(args argsparser.ConfigureSubcommandOptions) {
	containerRegistry := ask.Str(ask.Question{
		Description:     "Container registry for cleanup",
		PossibleAnswers: []string{"gcr", "azurer"},
		DefaultValue:    "gcr",
	})

	var kubernetesClusters []string
	for {

		finishText := " (leave empty if you wish to finish)"
		if len(kubernetesClusters) == 0 {
			finishText = ""
		}

		kubernetesCluster := ask.Str(ask.Question{
			Description: fmt.Sprintf("Kubernetes cluster #%v context%v", len(kubernetesClusters)+1, finishText),
		})

		if kubernetesCluster == "" {
			if len(kubernetesClusters) == 0 {
				fmt.Println("You need to give me at least one valid kubernetes cluster")
				continue
			}

			break
		}

		if !stringutil.StrInSlice(kubernetesCluster, kubernetesClusters) {
			kubernetesClusters = append(kubernetesClusters, kubernetesCluster)
		}
	}

	fmt.Printf("Container registry: %v", containerRegistry)
	fmt.Printf("kubernetes cluster contexts: %v", kubernetesClusters)
}
