package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/hytromo/faulty-crane/internal/argsparser"
	"github.com/hytromo/faulty-crane/internal/configurationhelper"
	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/imagefilters"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
	"github.com/hytromo/faulty-crane/internal/reporter"
	color "github.com/logrusorgru/aurora"
)

func main() {
	initLogging()

	appOptions, err := argsparser.Parse(os.Args)

	if err != nil {
		log.Fatal(err.Error())
	}
	switch {
	case appOptions.Clean.SubcommandEnabled:
		{
			options := appOptions.Clean

			var parsedRepos []imagefilters.ParsedRepo

			// parsed repos are read from a plan file only if it is specified during a normal clean run
			if !options.DryRun && options.Plan != "" {
				// normal run, reading from an existent plan file the parsed repos
				log.Infof("Reading from plan file %v\n", options.Plan)
				parsedRepos = configurationhelper.ReadPlan(options.Plan)
			} else {
				parsedRepos = imagefilters.Parse(
					containerregistry.MakeGCRClient(containerregistry.GCRClient{
						Host:      options.ContainerRegistry.Host,
						AccessKey: options.ContainerRegistry.Access,
					}).GetAllRepos(),
					options.Keep,
				)
			}

			if options.DryRun {
				reporter.ReportRepositoriesStatus(parsedRepos, appOptions.Clean.AnalyticalPlan)

				if options.Plan != "" {
					configurationhelper.WritePlan(parsedRepos, options.Plan)
					log.Infof(
						"Plan saved to %v, run '%v clean -plan %v -config ...' to delete exactly the images that have been marked for deletion above",
						options.Plan, filepath.Base(os.Args[0]), options.Plan,
					)
				}
			} else {
				// TODO: really clean the parsed repos
				gcrClient := containerregistry.MakeGCRClient(containerregistry.GCRClient{
					Host:      options.ContainerRegistry.Host,
					AccessKey: options.ContainerRegistry.Access,
				})
				for _, repo := range parsedRepos {
					for _, image := range repo.Images {
						if image.KeptData.Reason == keepreasons.None {
							log.Infof("Need to delete %v with tags %v and digest %v", image.Image.Repo, image.Image.Tag, image.Image.Digest)
							b, err := json.MarshalIndent(image, "", "  ")
							if err != nil {
								fmt.Println(err)
							}
							fmt.Print(string(b))
							gcrClient.DeleteImage(
								strings.Replace(image.Image.Repo, "eu.gcr.io/", "", 1), image.Image)
							break
						}
					}
					if true {
						break
					}
				}
			}
		}
	case appOptions.Configure.SubcommandEnabled:
		{
			configurationhelper.CreateNew(appOptions.Configure)
			fmt.Printf("Configuration written in %v\n", color.Green(appOptions.Configure.Config))
		}
	case appOptions.Show.SubcommandEnabled:
		{
			parsedRepos := configurationhelper.ReadPlan(appOptions.Show.Plan)
			reporter.ReportRepositoriesStatus(parsedRepos, appOptions.Show.AnalyticalPlan)
		}
	}
}
