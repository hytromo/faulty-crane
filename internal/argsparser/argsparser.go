package argsparser

import (
	"errors"
	"flag"
	"fmt"
	"strings"

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

// Parse parses a list of strings as cli options and returns the final configuration.
// Returns an error if the list of strings cannot be parsed.
func Parse(args []string) (configuration.AppOptions, error) {
	cleanSubCmd := "clean"
	configureSubCmd := "configure"
	showSubCmd := "show"

	var appOptions configuration.AppOptions

	subCommandsMap := map[string]func(){
		cleanSubCmd: func() {
			appOptions.Clean.SubcommandEnabled = true

			cleanCmd := flag.NewFlagSet(cleanSubCmd, flag.ExitOnError)

			cleanCmd.BoolVar(&appOptions.Clean.DryRun, "dry-run", false, "just output what is expected to be deleted without actually deleting anything")

			cleanCmd.StringVar(&appOptions.Clean.Plan, "plan", "", "a plan file: use with -dry-run to create a new plan file containing the images marked for deletion; use without -dry-run to read from a plan file which images to delete (if a plan file is specified all the other filters are skipped/ignored)")

			cleanCmd.StringVar(&appOptions.Clean.Config, "config", "", "path to the configuration file; can be created through 'faulty-crane configure'; other options can override the configuration")

			cleanCmd.StringVar(&appOptions.Clean.ContainerRegistry.Host, "registry", "", "the registry to clean, e.g. eu.gcr.io")
			cleanCmd.StringVar(&appOptions.Clean.ContainerRegistry.Access, "key", "", "the registry access key, e.g. 'gcloud auth print-access-token'")

			cleanCmd.StringVar(&appOptions.Clean.Keep.YoungerThan, "younger-than", "", "images younger than this value will be kept; provide a duration value, e.g. '10d', '1w3d' or '1d3h'")

			k8sClustersStr := cleanCmd.String("keep-used-in-k8s", "", "comma-separated list of k8s contexts; any image that is used by these clusters won't be deleted")

			imageTags := cleanCmd.String("keep-image-tags", "", "comma-separated list of tags; images with any of these tags will be kept")

			imageDigests := cleanCmd.String("keep-image-digests", "", "comma-separated list of digests; images with these digests will be kept")

			imageIDs := cleanCmd.String("keep-image-ids", "", "comma-separated list of IDs; images with these IDs will be kept")

			cleanCmd.Parse(args[2:])

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
			configureCmd.StringVar(&appOptions.Configure.Config, "o", "faulty-crane.json", "the file to save the configuration to")
			configureCmd.Parse(args[2:])
		},
		showSubCmd: func() {
			appOptions.Show.SubcommandEnabled = true

			showCmd := flag.NewFlagSet(showSubCmd, flag.ExitOnError)
			showCmd.StringVar(&appOptions.Show.Plan, "plan", "", "the plan file to show")
			showCmd.Parse(args[2:])
		},
	}

	if len(args) < 2 {
		return appOptions, getWrongOptionsError(subCommandsMap)
	}

	parseCliOptionsOfSubcommand, subcommandExists := subCommandsMap[args[1]]

	if !subcommandExists {
		return appOptions, getWrongOptionsError(subCommandsMap)
	}

	parseCliOptionsOfSubcommand()

	return appOptions, optionsvalidator.Validate(appOptions)
}
