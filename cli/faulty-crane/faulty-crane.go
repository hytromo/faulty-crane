package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/hytromo/faulty-crane/internal/argsparser"
	"github.com/hytromo/faulty-crane/internal/configuration"
	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/imagefilters"
	"github.com/hytromo/faulty-crane/internal/orchestrator"
	"github.com/hytromo/faulty-crane/internal/reporter"
	color "github.com/logrusorgru/aurora"
)

func main() {
	initLogging()

	appOptions, err := argsparser.Parse(os.Args)

	if err != nil {
		log.Fatal(err.Error())
	}

	if appOptions.Configure.SubcommandEnabled {
		configuration.CreateNew(appOptions.Configure, nil)
		fmt.Printf("Configuration written in %v\n", color.Green(appOptions.Configure.Config))
		return
	}

	if appOptions.Show.SubcommandEnabled {
		parsedRepos := configuration.ReadPlan(appOptions.Show.Plan, true)
		reporter.ReportRepositoriesStatus(parsedRepos, appOptions.Show.Analytical)
	}

	if appOptions.Apply.SubcommandEnabled || appOptions.Plan.SubcommandEnabled {
		options := appOptions.ApplyPlanCommon

		var parsedRepos []containerregistry.Repository

		// parsed repos are read from a plan file only if it is specified during a normal apply run
		if appOptions.Apply.SubcommandEnabled && options.Plan != "" {
			// normal run, reading from an existent plan file the parsed repos
			log.Infof("Reading from plan file %v\n", options.Plan)
			parsedRepos = configuration.ReadPlan(options.Plan, true)
		} else {
			orchestrator := orchestrator.NewOrchestrator(&appOptions)
			orchestrator.Init()
			parsedRepos = imagefilters.Parse(
				orchestrator.GetAllRepos(),
				options.Keep,
			)
		}

		if appOptions.Plan.SubcommandEnabled {
			reporter.ReportRepositoriesStatus(parsedRepos, false)

			if options.Plan == "" {
				return
			}

			configuration.WritePlan(parsedRepos, options.Plan)

			// if the user used a config to produce the dry run, they can use the same config to execute the plan, so here we prepare fully the command for them
			configStrInfo := ""
			if options.Config != "" {
				configStrInfo = fmt.Sprintf(" -config %v", options.Config)
			}

			log.Infof("Plan saved to %v", options.Plan)

			fmt.Printf(
				"\n\nTo delete exactly what is planned:\n",
			)

			fmt.Printf(
				"    %v apply%v %v",
				filepath.Base(os.Args[0]), configStrInfo, options.Plan,
			)

			fmt.Printf(
				"\n\nTo show analytically what is going to be kept:\n",
			)

			fmt.Printf(
				"    %v show -analytical -plan %v\n",
				filepath.Base(os.Args[0]), options.Plan,
			)
		}

		if appOptions.Apply.SubcommandEnabled {
			orchestrator := orchestrator.NewOrchestrator(&appOptions)
			orchestrator.Init()
			results := orchestrator.DeleteImagesWithNoKeepReason(parsedRepos)

			if results.ShouldDeleteCount > 0 {
				log.Infof("Deleted %.2f%% (%v/%v) of the images", float64(results.ManagedToDeleteCount)/float64(results.ShouldDeleteCount)*100, results.ManagedToDeleteCount, results.ShouldDeleteCount)
			} else {
				log.Info("Nothing to do")
			}
		}
	}
}
