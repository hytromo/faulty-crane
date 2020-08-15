package main

import (
	"fmt"
	"os"

	"github.com/hytromo/faulty-crane/internal/argsParser"
	"github.com/hytromo/faulty-crane/internal/configurationHelper"
)

func main() {
	fmt.Println("New version v5")
	cliOptions := argsParser.Parse(os.Args)

	fmt.Printf("parsed cli options is: %+v\n", cliOptions)

	if cliOptions.Clean.SubcommandEnabled {

	} else if cliOptions.Configure.SubcommandEnabled {
		configurationHelper.CreateNewConfiguration(cliOptions.Configure)
	}
}
