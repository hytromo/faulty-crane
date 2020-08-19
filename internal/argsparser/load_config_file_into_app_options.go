package argsparser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"github.com/hytromo/faulty-crane/internal/configuration"
)

func replaceMissingAppOptionsFromConfig(appOptions *configuration.AppOptions, configPath string) {
	configBytes, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not read configuration file %v: %v\n", configPath, err))
	}

	configOptions := configuration.Configuration{}

	err = json.Unmarshal([]byte(configBytes), &configOptions)

	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid format of configuration file %v: %v\n", configPath, err))
	}

	// cli options override config options, so config options should fill in the blanks only
	if appOptions.Clean.ContainerRegistry.Host == "" {
		appOptions.Clean.ContainerRegistry.Host = configOptions.ContainerRegistry.Host
	}

	if appOptions.Clean.ContainerRegistry.Access == "" {
		appOptions.Clean.ContainerRegistry.Access = configOptions.ContainerRegistry.Access
	}

	if len(appOptions.Clean.Keep.UsedIn.KubernetesClusters) == 0 {
		appOptions.Clean.Keep.UsedIn.KubernetesClusters = configOptions.Keep.UsedIn.KubernetesClusters
	}

	if len(appOptions.Clean.Keep.Image.Digests) == 0 {
		appOptions.Clean.Keep.Image.Digests = configOptions.Keep.Image.Digests
	}

	if len(appOptions.Clean.Keep.Image.Tags) == 0 {
		appOptions.Clean.Keep.Image.Tags = configOptions.Keep.Image.Tags
	}

	if len(appOptions.Clean.Keep.Image.IDs) == 0 {
		appOptions.Clean.Keep.Image.IDs = configOptions.Keep.Image.IDs
	}
}
