package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hytromo/faulty-crane/internal/argsparser"
	"github.com/hytromo/faulty-crane/internal/configurationhelper"
	color "github.com/logrusorgru/aurora"
)

func main() {
	appOptions, err := argsparser.Parse(os.Args)

	if err != nil {
		log.Fatal(err.Error())
	}

	if appOptions.Clean.SubcommandEnabled {
		fmt.Printf("App options is %+v\n", appOptions)
	} else if appOptions.Configure.SubcommandEnabled {
		configurationhelper.CreateNew(appOptions.Configure)
		fmt.Printf("Configuration written in %v\n", color.Green(appOptions.Configure.Config))
	}
}
