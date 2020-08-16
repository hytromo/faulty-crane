package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hytromo/faulty-crane/internal/argsparser"
	"github.com/hytromo/faulty-crane/internal/configurationhelper"
)

func main() {
	appOptions, err := argsparser.Parse(os.Args)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("parsed cli options is: %+v\n", appOptions)

	if appOptions.Clean.SubcommandEnabled {

	} else if appOptions.Configure.SubcommandEnabled {
		configurationhelper.CreateNew(appOptions.Configure)
	}
}
