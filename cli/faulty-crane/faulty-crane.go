package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/hytromo/faulty-crane/internal/argsparser"
	"github.com/hytromo/faulty-crane/internal/configurationhelper"
	"github.com/hytromo/faulty-crane/internal/containerregistry"
	color "github.com/logrusorgru/aurora"
)

func main() {
	initLogging()

	appOptions, err := argsparser.Parse(os.Args)

	if err != nil {
		log.Fatal(err.Error())
	}

	if appOptions.Clean.SubcommandEnabled {
		fmt.Printf("App options is %+v\n", appOptions)
		client := containerregistry.MakeGCRClient(containerregistry.GCRClient{
			Host:      appOptions.Clean.ContainerRegistry.Host,
			AccessKey: appOptions.Clean.ContainerRegistry.Access,
		})

		client.GetAllImages()
	} else if appOptions.Configure.SubcommandEnabled {
		configurationhelper.CreateNew(appOptions.Configure)
		fmt.Printf("Configuration written in %v\n", color.Green(appOptions.Configure.Config))
	}
}
