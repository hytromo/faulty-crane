package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/hytromo/faulty-crane/internal/argsparser"
	"github.com/hytromo/faulty-crane/internal/configurationhelper"
	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/imagefilters"
	"github.com/hytromo/faulty-crane/internal/reporter"
	color "github.com/logrusorgru/aurora"
)

func main() {
	initLogging()

	appOptions, err := argsparser.Parse(os.Args)

	if err != nil {
		log.Fatal(err.Error())
	}

	if appOptions.Clean.SubcommandEnabled {
		options := appOptions.Clean

		client := containerregistry.MakeGCRClient(containerregistry.GCRClient{
			Host:      options.ContainerRegistry.Host,
			AccessKey: options.ContainerRegistry.Access,
		})

		parsedRepos := imagefilters.Parse(client.GetAllRepos(), options.Keep)

		if options.DryRun {
			reporter.ReportRepositoriesStatus(parsedRepos)

			if options.Plan != "" {
				configurationhelper.WritePlan(parsedRepos, options.Plan)
				log.Infof(
					"Plan saved to %v, run '%v clean -plan %v -config ...' to delete exactly the images that have been marked for deletion above",
					options.Plan, filepath.Base(os.Args[0]), options.Plan,
				)
			}
		} else {
			// TODO: really clean
		}
	} else if appOptions.Configure.SubcommandEnabled {
		configurationhelper.CreateNew(appOptions.Configure)
		fmt.Printf("Configuration written in %v\n", color.Green(appOptions.Configure.Config))
	}
}
