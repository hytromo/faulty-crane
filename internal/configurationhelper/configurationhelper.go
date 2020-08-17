package configurationhelper

import (
	"log"

	"github.com/hytromo/faulty-crane/internal/argsparser"
	"github.com/hytromo/faulty-crane/internal/utils/fileutil"
)

func constructConfigurationFromAnswers(answers UserInput) Configuration {
	config := Configuration{}
	config.ContainerRegistry = containerRegistry{
		Type: answers.containerRegistry,
	}

	config.Keep.YoungerThan = answers.youngerThan

	config.Keep.UsedIn.KubernetesClusters = make([]kubernetesCluster, len(answers.kubernetesClusters))

	for i, cluster := range answers.kubernetesClusters {
		config.Keep.UsedIn.KubernetesClusters[i].Context = cluster
	}

	config.Keep.Image.Tags = answers.imageTags
	config.Keep.Image.Digests = answers.imageDigests
	config.Keep.Image.IDs = answers.imageIDs
	config.Keep.Image.Repositories = answers.imageRepositories

	return config
}

func saveConfig(path string, config Configuration) {
	err := fileutil.SaveJSON(path, config)

	if err != nil {
		log.Fatalf("Error saving configuration file: %v", err)
	}
}

// CreateNew asks the user for configuration input and then creates a configuration file based on the answers
func CreateNew(params argsparser.ConfigureSubcommandOptions) {
	answers := AskUserInput()
	config := constructConfigurationFromAnswers(answers)
	saveConfig(params.Config, config)
}
