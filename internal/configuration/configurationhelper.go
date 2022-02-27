package configuration

import (
	"encoding/json"
	"io"

	log "github.com/sirupsen/logrus"

	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/utils/fileutil"
)

func constructConfigurationFromAnswers(answers UserInput) Configuration {
	config := Configuration{}
	config.GCR = GoogleContainerRegistry{
		Token: answers.ContainerRegistryPassword,
		Host:  answers.ContainerRegistryLink,
	}

	config.Dockerhub = DockerhubContainerRegistry{
		Username:  answers.ContainerRegistryUsername,
		Password:  answers.ContainerRegistryPassword,
		Namespace: answers.ContainerRegistryNamespace,
	}

	config.Keep.YoungerThan = answers.YoungerThan

	config.Keep.UsedIn.KubernetesClusters = make([]KubernetesCluster, len(answers.KubernetesClusters))

	for i, cluster := range answers.KubernetesClusters {
		config.Keep.UsedIn.KubernetesClusters[i].Context = cluster
	}

	config.Keep.Image.Tags = answers.ImageTags
	config.Keep.Image.Digests = answers.ImageDigests
	config.Keep.Image.Repositories = answers.ImageIDs

	return config
}

func saveConfig(path string, config Configuration) {
	err := fileutil.SaveJSON(path, config, false)

	if err != nil {
		log.Fatalf("Error saving configuration file: %v", err)
	}
}

// CreateNew asks the user for configuration input and then creates a configuration file based on the answers
func CreateNew(params ConfigureSubcommandOptions, reader io.Reader) {
	answers := AskUserInput(reader)
	config := constructConfigurationFromAnswers(answers)
	saveConfig(params.Config, config)
}

// WritePlan writes the parsed repos in a plan file; the plan file can then be used to remove specific images
func WritePlan(parsedRepos []containerregistry.Repository, planPath string) {
	fileutil.SaveJSON(planPath, parsedRepos, true)
}

// ReadPlan reads a plan file and returns the parsed repositories
func ReadPlan(planPath string) []containerregistry.Repository {
	planBytes, err := fileutil.ReadFile(planPath, true)
	if err != nil {
		log.Fatalf("Could not read plan file '%v': %v\n", planPath, err.Error())
	}

	parsedRepos := []containerregistry.Repository{}

	err = json.Unmarshal([]byte(planBytes), &parsedRepos)

	if err != nil {
		log.Fatalf("Cannot parse json of plan file %v: %v\n", planPath, err.Error())
	}

	return parsedRepos
}

// IsGCR returns if the configuration options point to GCR
func IsGCR(options *AppOptions) bool {
	config := options.ApplyPlanCommon

	return config.GoogleContainerRegistry != (GoogleContainerRegistry{})
}

// IsDockerhub returns if the configuration options point to Dockerhub
func IsDockerhub(options *AppOptions) bool {
	config := options.ApplyPlanCommon

	return config.DockerhubContainerRegistry != (DockerhubContainerRegistry{})
}
