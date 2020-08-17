package configurationhelper

import (
	"log"

	"github.com/hytromo/faulty-crane/internal/configuration"
	"github.com/hytromo/faulty-crane/internal/utils/fileutil"
)

func constructConfigurationFromAnswers(answers UserInput) configuration.Configuration {
	config := configuration.Configuration{}
	config.ContainerRegistry = configuration.ContainerRegistry{
		Access: answers.containerRegistryAccess,
		Link:   answers.containerRegistryLink,
	}

	config.Keep.YoungerThan = answers.youngerThan

	config.Keep.UsedIn.KubernetesClusters = make([]configuration.KubernetesCluster, len(answers.kubernetesClusters))

	for i, cluster := range answers.kubernetesClusters {
		config.Keep.UsedIn.KubernetesClusters[i].Context = cluster
	}

	config.Keep.Image.Tags = answers.imageTags
	config.Keep.Image.Digests = answers.imageDigests
	config.Keep.Image.IDs = answers.imageIDs

	return config
}

func saveConfig(path string, config configuration.Configuration) {
	err := fileutil.SaveJSON(path, config)

	if err != nil {
		log.Fatalf("Error saving configuration file: %v", err)
	}
}

// CreateNew asks the user for configuration input and then creates a configuration file based on the answers
func CreateNew(params configuration.ConfigureSubcommandOptions) {
	answers := AskUserInput()
	config := constructConfigurationFromAnswers(answers)
	saveConfig(params.Config, config)
}
