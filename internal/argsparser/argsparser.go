package argsparser

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/hytromo/faulty-crane/internal/configuration"
	"github.com/hytromo/faulty-crane/internal/optionsvalidator"
)

// EnvPrefix is the common prefix of all the environment variables we respect
const EnvPrefix = "FAULTY_CRANE_"

func getWrongOptionsError(subCommandsMap map[string]func()) (err error) {
	allSubcommands := make([]string, len(subCommandsMap))

	i := 0
	for subcommand := range subCommandsMap {
		allSubcommands[i] = subcommand
		i++
	}

	return errors.New(
		fmt.Sprintln(
			"Please specify one of the valid subcommands:",
			strings.Join(allSubcommands, ", "),
			"\nYou can use the -h/--help switch on the subcommands for further assistance on their usage",
		),
	)
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func lookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatalf("LookupEnvOrBool[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

// it parses the arguments even when there are not enough of them
func safeParseArguments(flagset *flag.FlagSet, args []string) {
	if len(args) > 2 {
		flagset.Parse(args[2:])
	} else {
		flagset.Parse([]string{})
	}
}

/*
This function is required so that sensitive values passed from the environment are not shown in stdout when listing command usage through --help
*/
func registerStrParameter(cmd *flag.FlagSet, p *string, name string, envKey string, defaultValue string, usage string) {
	defaultValueFriendly := "empty"

	if defaultValue != "" {
		defaultValueFriendly = defaultValue
	}

	cmd.StringVar(p, name, "", fmt.Sprintf("Description: %v\nEnv var:     %v\nDefault:     %v", usage, envKey, defaultValueFriendly))

	if *p == "" {
		*p = lookupEnvOrString(envKey, defaultValue)
	}
}
func registerBoolParameter(cmd *flag.FlagSet, p *bool, name string, envKey string, defaultValue bool, usage string) {
	defaultValueFriendly := "false"

	if defaultValue {
		defaultValueFriendly = "true"
	}

	cmd.BoolVar(p, name, false, fmt.Sprintf("Description: %v\nEnv var:     %v\nDefault:     %v", usage, envKey, defaultValueFriendly))

	if !*p {
		*p = lookupEnvOrBool(envKey, defaultValue)
	}
}

func addApplyPlanCommonVars(cmd *flag.FlagSet, appOptions *configuration.AppOptions, args []string) {

	registerStrParameter(cmd, &appOptions.ApplyPlanCommon.Config, "config", EnvPrefix+"CONFIG", "", "path to the configuration file; can be created through 'faulty-crane configure'; other options can override the configuration")

	registerStrParameter(cmd, &appOptions.ApplyPlanCommon.GoogleContainerRegistry.Host, "registry", EnvPrefix+"GOOGLE_CONTAINER_REGISTRY_HOST", "", "the registry to clean, e.g. eu.gcr.io")

	password := ""
	registerStrParameter(cmd, &password, "password", EnvPrefix+"CONTAINER_REGISTRY_PASSWORD", "", "the registry password, access key etc. For GCR it's the output of 'gcloud auth print-access-token', we HIGHLY recommend you use an env variable for this")

	username := ""
	registerStrParameter(cmd, &username, "username", EnvPrefix+"CONTAINER_REGISTRY_USERNAME", "", "the registry username, not all registries require this, e.g. GCR does not")

	registerStrParameter(cmd, &appOptions.ApplyPlanCommon.Keep.YoungerThan, "keep-younger-than", EnvPrefix+"KEEP_YOUNGER_THAN", "", "images younger than this value will be kept; provide a duration value, e.g. '10d', '1w3d' or '1d3h'")

	atLeastStr := ""
	registerStrParameter(cmd, &atLeastStr, "keep-at-least", EnvPrefix+"KEEP_AT_LEAST", "", "at least that many images will be kept in this specific repo, prioritising the younger ones")

	k8sClustersStr := ""
	imageTags := ""
	imageDigests := ""
	imageIDs := ""

	registerStrParameter(cmd, &k8sClustersStr, "keep-used-in-k8s", EnvPrefix+"KEEP_USED_IN_K8S", "", "comma-separated list of k8s contexts; any image that is used by these clusters won't be deleted")

	registerStrParameter(cmd, &imageTags, "keep-image-tags", EnvPrefix+"KEEP_IMAGE_TAGS", "", "comma-separated list of tags; images with any of these tags will be kept")

	registerStrParameter(cmd, &imageDigests, "keep-image-digests", EnvPrefix+"KEEP_IMAGE_DIGESTS", "", "comma-separated list of digests; images with these digests will be kept")

	registerStrParameter(cmd, &imageIDs, "keep-image-repos", EnvPrefix+"KEEP_IMAGE_REPOS", "", "comma-separated list of repos; images with in these repos will be kept")

	safeParseArguments(cmd, args)

	appOptions.ApplyPlanCommon.Keep.AtLeast = 0
	if atLeastStr != "" {
		atLeast, err := strconv.Atoi(atLeastStr)

		if err != nil {
			log.Fatalf("Could not convert keep-at-least value '%s' to integer", atLeastStr)
		}

		appOptions.ApplyPlanCommon.Keep.AtLeast = atLeast
	}

	if len(k8sClustersStr) > 0 {
		k8sClustersArr := strings.Split(k8sClustersStr, ",")
		appOptions.ApplyPlanCommon.Keep.UsedIn.KubernetesClusters = make([]configuration.KubernetesCluster, len(k8sClustersArr))
		for i, context := range k8sClustersArr {
			appOptions.ApplyPlanCommon.Keep.UsedIn.KubernetesClusters[i] = configuration.KubernetesCluster{
				Context: context,
			}
		}
	}

	if len(imageTags) > 0 {
		imageTagsArr := strings.Split(imageTags, ",")
		appOptions.ApplyPlanCommon.Keep.Image.Tags = make([]string, len(imageTagsArr))
		for i, imageTag := range imageTagsArr {
			appOptions.ApplyPlanCommon.Keep.Image.Tags[i] = imageTag
		}
	}

	if len(imageDigests) > 0 {
		imageDigestsArr := strings.Split(imageDigests, ",")
		appOptions.ApplyPlanCommon.Keep.Image.Digests = make([]string, len(imageDigestsArr))
		for i, imageTag := range imageDigestsArr {
			appOptions.ApplyPlanCommon.Keep.Image.Digests[i] = imageTag
		}
	}

	if len(imageIDs) > 0 {
		imageIDsArr := strings.Split(imageIDs, ",")
		appOptions.ApplyPlanCommon.Keep.Image.Repositories = make([]string, len(imageIDsArr))
		for i, imageTag := range imageIDsArr {
			appOptions.ApplyPlanCommon.Keep.Image.Repositories[i] = imageTag
		}
	}

	if appOptions.ApplyPlanCommon.Config != "" {
		replaceMissingAppOptionsFromConfig(appOptions, appOptions.ApplyPlanCommon.Config)
	}

	if configuration.IsGCR(appOptions) {
		appOptions.ApplyPlanCommon.GoogleContainerRegistry.Token = password
	} else if configuration.IsDockerhub(appOptions) {
		appOptions.ApplyPlanCommon.DockerhubContainerRegistry.Password = password
		appOptions.ApplyPlanCommon.DockerhubContainerRegistry.Username = username
	}
}

// Parse parses a list of strings as cli options and returns the final configuration.
// Returns an error if the list of strings cannot be parsed.
func Parse(args []string) (configuration.AppOptions, error) {
	applySubCmd := "apply"
	planSubCmd := "plan"
	configureSubCmd := "configure"
	showSubCmd := "show"

	var appOptions configuration.AppOptions

	subCommandsMap := map[string]func(){
		planSubCmd: func() {
			appOptions.Plan.SubcommandEnabled = true

			planCmd := flag.NewFlagSet(planSubCmd, flag.ExitOnError)

			registerStrParameter(planCmd, &appOptions.ApplyPlanCommon.Plan, "out", EnvPrefix+"PLAN", "", "a plan file to write")

			addApplyPlanCommonVars(planCmd, &appOptions, args)
		},
		applySubCmd: func() {
			appOptions.Apply.SubcommandEnabled = true

			applyCmd := flag.NewFlagSet(applySubCmd, flag.ExitOnError)

			addApplyPlanCommonVars(applyCmd, &appOptions, args)

			if applyCmd.NArg() != 0 {
				// has plan nfile
				appOptions.ApplyPlanCommon.Plan = applyCmd.Args()[0]
				log.Info("Got plan file ", appOptions.ApplyPlanCommon.Plan)
			}

		},
		configureSubCmd: func() {
			appOptions.Configure.SubcommandEnabled = true

			configureCmd := flag.NewFlagSet(configureSubCmd, flag.ExitOnError)
			registerStrParameter(configureCmd, &appOptions.Configure.Config, "out", EnvPrefix+"CONFIG", filepath.Base(os.Args[0])+".json", "the file to save the configuration to")
			safeParseArguments(configureCmd, args)
		},
		showSubCmd: func() {
			appOptions.Show.SubcommandEnabled = true

			showCmd := flag.NewFlagSet(showSubCmd, flag.ExitOnError)
			registerStrParameter(showCmd, &appOptions.Show.Plan, "plan", EnvPrefix+"PLAN", "plan.out", "the plan file to show")
			registerBoolParameter(showCmd, &appOptions.Show.Analytical, "analytical", EnvPrefix+"ANALYTICAL", false, "print the whole plan, not an aggregation")
			safeParseArguments(showCmd, args)
		},
	}

	chosenCommand := "non-existent-subcommand"

	if len(args) >= 2 {
		chosenCommand = args[1]
	}

	parseCliOptionsOfSubcommand, subcommandExists := subCommandsMap[chosenCommand]

	if !subcommandExists {
		return appOptions, getWrongOptionsError(subCommandsMap)
	}

	parseCliOptionsOfSubcommand()

	return appOptions, optionsvalidator.Validate(appOptions)
}
