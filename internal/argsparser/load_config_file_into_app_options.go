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
	if appOptions.ApplyPlanCommon.GoogleContainerRegistry.Host == "" {
		appOptions.ApplyPlanCommon.GoogleContainerRegistry.Host = configOptions.GCR.Host
	}

	if appOptions.ApplyPlanCommon.GoogleContainerRegistry.Token == "" {
		appOptions.ApplyPlanCommon.GoogleContainerRegistry.Token = configOptions.GCR.Token
	}

	if appOptions.ApplyPlanCommon.DockerhubContainerRegistry.Namespace == "" {
		appOptions.ApplyPlanCommon.DockerhubContainerRegistry.Namespace = configOptions.Dockerhub.Namespace
	}

	if appOptions.ApplyPlanCommon.DockerhubContainerRegistry.Username == "" {
		appOptions.ApplyPlanCommon.DockerhubContainerRegistry.Username = configOptions.Dockerhub.Username
	}

	if appOptions.ApplyPlanCommon.DockerhubContainerRegistry.Password == "" {
		appOptions.ApplyPlanCommon.DockerhubContainerRegistry.Password = configOptions.Dockerhub.Password
	}

	if len(appOptions.ApplyPlanCommon.Keep.UsedIn.KubernetesClusters) == 0 {
		appOptions.ApplyPlanCommon.Keep.UsedIn.KubernetesClusters = configOptions.Keep.UsedIn.KubernetesClusters
	}

	if appOptions.ApplyPlanCommon.Keep.YoungerThan == "" {
		appOptions.ApplyPlanCommon.Keep.YoungerThan = configOptions.Keep.YoungerThan
	}

	if appOptions.ApplyPlanCommon.Keep.AtLeast == 0 {
		appOptions.ApplyPlanCommon.Keep.AtLeast = configOptions.Keep.AtLeast
	}

	if len(appOptions.ApplyPlanCommon.Keep.Image.Digests) == 0 {
		appOptions.ApplyPlanCommon.Keep.Image.Digests = configOptions.Keep.Image.Digests
	}

	if len(appOptions.ApplyPlanCommon.Keep.Image.Tags) == 0 {
		appOptions.ApplyPlanCommon.Keep.Image.Tags = configOptions.Keep.Image.Tags
	}

	if len(appOptions.ApplyPlanCommon.Keep.Image.Repositories) == 0 {
		appOptions.ApplyPlanCommon.Keep.Image.Repositories = configOptions.Keep.Image.Repositories
	}
}
