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

const ENV_PREFIX = "FAULTY_CRANE_"

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

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func LookupEnvOrBool(key string, defaultVal bool) bool {
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
		*p = LookupEnvOrString(envKey, defaultValue)
	}
}

func registerBoolParameter(cmd *flag.FlagSet, p *bool, name string, envKey string, defaultValue bool, usage string) {
	defaultValueFriendly := "false"

	if defaultValue {
		defaultValueFriendly = "true"
	}

	cmd.BoolVar(p, name, false, fmt.Sprintf("Description: %v\nEnv var:     %v\nDefault:     %v", usage, envKey, defaultValueFriendly))

	if !*p {
		*p = LookupEnvOrBool(envKey, defaultValue)
	}
}

func addApplyPlanCommonVars(cmd *flag.FlagSet, appOptions *configuration.AppOptions, args []string) {

	registerStrParameter(cmd, &appOptions.ApplyPlanCommon.Config, "config", ENV_PREFIX+"CONFIG", "", "path to the configuration file; can be created through 'faulty-crane configure'; other options can override the configuration")

	registerStrParameter(cmd, &appOptions.ApplyPlanCommon.ContainerRegistry.Host, "registry", ENV_PREFIX+"CONTAINER_REGISTRY_HOST", "", "the registry to clean, e.g. eu.gcr.io")
	// cmd.StringVar(&appOptions.ApplyPlanCommon.ContainerRegistry.Access, "key", LookupEnvOrString(ENV_PREFIX+"CONTAINER_REGISTRY_ACCESS", ""), "the path to the registry access key file, e.g. a file containing the output of 'gcloud auth print-access-token'")
	registerStrParameter(cmd, &appOptions.ApplyPlanCommon.ContainerRegistry.Access, "key", ENV_PREFIX+"CONTAINER_REGISTRY_ACCESS", "", "the registry access key, e.g. the output of 'gcloud auth print-access-token', we highly recommend you use an env variable for this")

	registerStrParameter(cmd, &appOptions.ApplyPlanCommon.Keep.YoungerThan, "keep-younger-than", ENV_PREFIX+"KEEP_YOUNGER_THAN", "", "images younger than this value will be kept; provide a duration value, e.g. '10d', '1w3d' or '1d3h'")

	k8sClustersStr := ""
	imageTags := ""
	imageDigests := ""
	imageIDs := ""

	registerStrParameter(cmd, &k8sClustersStr, "keep-used-in-k8s", ENV_PREFIX+"KEEP_USED_IN_K8S", "", "comma-separated list of k8s contexts; any image that is used by these clusters won't be deleted")

	registerStrParameter(cmd, &imageTags, "keep-image-tags", ENV_PREFIX+"KEEP_IMAGE_TAGS", "", "comma-separated list of tags; images with any of these tags will be kept")

	registerStrParameter(cmd, &imageDigests, "keep-image-digests", ENV_PREFIX+"KEEP_IMAGE_DIGESTS", "", "comma-separated list of digests; images with these digests will be kept")

	registerStrParameter(cmd, &imageIDs, "keep-image-repos", ENV_PREFIX+"KEEP_IMAGE_REPOS", "", "comma-separated list of repos; images with in these repos will be kept")

	safeParseArguments(cmd, args)

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
}

// Parse parses a list of strings as cli options and returns the final configuration.
// Returns an error if the list of strings cannot be parsed.
func Parse(args []string) (configuration.AppOptions, error) {
	applySubCmd := "apply"
	planSubCmd := "plan"
	configureSubCmd := "configure"
	showSubCmd := "show"
	defaultSubCmd := planSubCmd

	var appOptions configuration.AppOptions

	subCommandsMap := map[string]func(){
		planSubCmd: func() {
			appOptions.Plan.SubcommandEnabled = true

			planCmd := flag.NewFlagSet(planSubCmd, flag.ExitOnError)

			registerStrParameter(planCmd, &appOptions.ApplyPlanCommon.Plan, "out", ENV_PREFIX+"PLAN", "", "a plan file to write")

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
			registerStrParameter(configureCmd, &appOptions.Configure.Config, "out", ENV_PREFIX+"CONFIG", filepath.Base(os.Args[0])+".json", "the file to save the configuration to")
			safeParseArguments(configureCmd, args)
		},
		showSubCmd: func() {
			appOptions.Show.SubcommandEnabled = true

			showCmd := flag.NewFlagSet(showSubCmd, flag.ExitOnError)
			registerStrParameter(showCmd, &appOptions.Show.Plan, "plan", ENV_PREFIX+"PLAN", "plan.out", "the plan file to show")
			registerBoolParameter(showCmd, &appOptions.Show.Analytical, "analytical", ENV_PREFIX+"ANALYTICAL", false, "print the whole plan, not an aggregation")
			safeParseArguments(showCmd, args)
		},
	}

	chosenCommand := defaultSubCmd

	if len(args) < 2 {
		log.Infof("No subcommand specified in arguments, assuming %v", LookupEnvOrString(ENV_PREFIX+"SUBCMD", defaultSubCmd))
	} else {
		chosenCommand = args[1]
	}

	parseCliOptionsOfSubcommand, subcommandExists := subCommandsMap[chosenCommand]

	if !subcommandExists {
		return appOptions, getWrongOptionsError(subCommandsMap)
	}

	parseCliOptionsOfSubcommand()

	return appOptions, optionsvalidator.Validate(appOptions)
}
