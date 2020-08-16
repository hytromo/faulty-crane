package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hytromo/faulty-crane/internal/argsparser"
	"github.com/hytromo/faulty-crane/internal/configurationhelper"
)

func main() {
	cliOptions, err := argsparser.Parse(os.Args)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("parsed cli options is: %+v\n", cliOptions)

	if cliOptions.Clean.SubcommandEnabled {

	} else if cliOptions.Configure.SubcommandEnabled {
		configurationhelper.CreateNewConfiguration(cliOptions.Configure)
	}
}
