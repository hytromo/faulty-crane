package argsparser

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/hytromo/faulty-crane/internal/configuration"
	"github.com/hytromo/faulty-crane/internal/optionsvalidator"
)

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

// Parse parses a list of strings as cli options and returns the final configuration.
// Returns an error if the list of strings cannot be parsed.
func Parse(args []string) (configuration.AppOptions, error) {
	cleanSubCmd := "clean"
	configureSubCmd := "configure"
	showSubCmd := "show"
	defaultSubCmd := cleanSubCmd
	const ENV_PREFIX = "FAULTY_CRANE_"

	var appOptions configuration.AppOptions

	subCommandsMap := map[string]func(){
		cleanSubCmd: func() {
			appOptions.Clean.SubcommandEnabled = true

			cleanCmd := flag.NewFlagSet(cleanSubCmd, flag.ExitOnError)

			cleanCmd.BoolVar(&appOptions.Clean.DryRun, "dry-run", LookupEnvOrBool(ENV_PREFIX+"DRY_RUN", false), "just output what is expected to be deleted without actually deleting anything")
			cleanCmd.BoolVar(&appOptions.Clean.AnalyticalPlan, "analytically", LookupEnvOrBool(ENV_PREFIX+"ANALYTICALLY", false), "print the whole plan, not an aggregation")

			cleanCmd.StringVar(&appOptions.Clean.Plan, "plan", LookupEnvOrString(ENV_PREFIX+"PLAN", ""), "a plan file: use with -dry-run to create a new plan file containing the images marked for deletion; use without -dry-run to read from a plan file which images to delete (if a plan file is specified all the other filters are skipped/ignored)")

			cleanCmd.StringVar(&appOptions.Clean.Config, "config", LookupEnvOrString(ENV_PREFIX+"CONFIG", ""), "path to the configuration file; can be created through 'faulty-crane configure'; other options can override the configuration")

			cleanCmd.StringVar(&appOptions.Clean.ContainerRegistry.Host, "registry", LookupEnvOrString(ENV_PREFIX+"CONTAINER_REGISTRY_HOST", ""), "the registry to clean, e.g. eu.gcr.io")
			cleanCmd.StringVar(&appOptions.Clean.ContainerRegistry.Access, "key", LookupEnvOrString(ENV_PREFIX+"CONTAINER_REGISTRY_ACCESS", ""), "the path to the registry access key file, e.g. a file containing the output of 'gcloud auth print-access-token'")

			cleanCmd.StringVar(&appOptions.Clean.Keep.YoungerThan, "keep-younger-than", LookupEnvOrString(ENV_PREFIX+"KEEP_YOUNGER_THAN", ""), "images younger than this value will be kept; provide a duration value, e.g. '10d', '1w3d' or '1d3h'")

			k8sClustersStr := cleanCmd.String("keep-used-in-k8s", LookupEnvOrString(ENV_PREFIX+"KEEP_USED_IN_K8S", ""), "comma-separated list of k8s contexts; any image that is used by these clusters won't be deleted")

			imageTags := cleanCmd.String("keep-image-tags", LookupEnvOrString(ENV_PREFIX+"KEEP_IMAGE_TAGS", ""), "comma-separated list of tags; images with any of these tags will be kept")

			imageDigests := cleanCmd.String("keep-image-digests", LookupEnvOrString(ENV_PREFIX+"KEEP_IMAGE_DIGESTS", ""), "comma-separated list of digests; images with these digests will be kept")

			imageIDs := cleanCmd.String("keep-image-repos", LookupEnvOrString(ENV_PREFIX+"KEEP_IMAGE_REPOS", ""), "comma-separated list of repos; images with in these repos will be kept")

			safeParseArguments(cleanCmd, args)

			if len(*k8sClustersStr) > 0 {
				k8sClustersArr := strings.Split(*k8sClustersStr, ",")
				appOptions.Clean.Keep.UsedIn.KubernetesClusters = make([]configuration.KubernetesCluster, len(k8sClustersArr))
				for i, context := range k8sClustersArr {
					appOptions.Clean.Keep.UsedIn.KubernetesClusters[i] = configuration.KubernetesCluster{
						Context: context,
					}
				}
			}

			if len(*imageTags) > 0 {
				imageTagsArr := strings.Split(*imageTags, ",")
				appOptions.Clean.Keep.Image.Tags = make([]string, len(imageTagsArr))
				for i, imageTag := range imageTagsArr {
					appOptions.Clean.Keep.Image.Tags[i] = imageTag
				}
			}

			if len(*imageDigests) > 0 {
				imageDigestsArr := strings.Split(*imageDigests, ",")
				appOptions.Clean.Keep.Image.Digests = make([]string, len(imageDigestsArr))
				for i, imageTag := range imageDigestsArr {
					appOptions.Clean.Keep.Image.Digests[i] = imageTag
				}
			}

			if len(*imageIDs) > 0 {
				imageIDsArr := strings.Split(*imageIDs, ",")
				appOptions.Clean.Keep.Image.Repositories = make([]string, len(imageIDsArr))
				for i, imageTag := range imageIDsArr {
					appOptions.Clean.Keep.Image.Repositories[i] = imageTag
				}
			}

			if appOptions.Clean.Config != "" {
				replaceMissingAppOptionsFromConfig(&appOptions, appOptions.Clean.Config)
			}
		},
		configureSubCmd: func() {
			appOptions.Configure.SubcommandEnabled = true

			configureCmd := flag.NewFlagSet(configureSubCmd, flag.ExitOnError)
			configureCmd.StringVar(&appOptions.Configure.Config, "o", LookupEnvOrString(ENV_PREFIX+"CONFIG", "faulty-crane.json"), "the file to save the configuration to")
			safeParseArguments(configureCmd, args)
		},
		showSubCmd: func() {
			appOptions.Show.SubcommandEnabled = true

			showCmd := flag.NewFlagSet(showSubCmd, flag.ExitOnError)
			showCmd.StringVar(&appOptions.Show.Plan, "plan", LookupEnvOrString(ENV_PREFIX+"PLAN", "plan.out"), "the plan file to show")
			showCmd.BoolVar(&appOptions.Show.AnalyticalPlan, "analytically", LookupEnvOrBool(ENV_PREFIX+"ANALYTICALLY", false), "print the whole plan, not an aggregation")
			// TODO: parallelize image deletion with prorgessbar etc
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
