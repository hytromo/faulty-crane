package argsParser

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type CleanCliOptions struct {
	subcommandEnabled bool
	dryRun            bool
}

type ConfigureCliOptions struct {
	subcommandEnabled bool
	config            string
}

type CliOptions struct {
	clean     CleanCliOptions
	configure ConfigureCliOptions
}

func exitWithWrongOptionsMessage(subCommandsMap map[string]func()) {
	allSubcommands := make([]string, len(subCommandsMap))

	i := 0
	for k := range subCommandsMap {
		allSubcommands[i] = k
		i++
	}

	fmt.Println("Please specify one of the valid subcommands:", strings.Join(allSubcommands, ", "))
	fmt.Println("You can use the -h/--help switch on the subcommands for further assistance on their usage")
	os.Exit(1)
}

func Parse(args []string) CliOptions {
	cleanSubCmd := "clean"
	configureSubCmd := "configure"

	var cliOptions CliOptions

	subCommandsMap := map[string]func(){
		cleanSubCmd: func() {
			cliOptions.clean.subcommandEnabled = true
			cleanCmd := flag.NewFlagSet(cleanSubCmd, flag.ExitOnError)
			cleanCmd.BoolVar(&cliOptions.clean.dryRun, "dry-run", false, "just output what is expected to be deleted")
			cleanCmd.Parse(args[2:])
		},
		configureSubCmd: func() {
			cliOptions.configure.subcommandEnabled = true
			configureCmd := flag.NewFlagSet(configureSubCmd, flag.ExitOnError)
			configureCmd.StringVar(&cliOptions.configure.config, "o", "faulty-crane.json", "the file to save the configuration to")
			configureCmd.Parse(args[2:])
		},
	}

	if len(args) < 2 {
		exitWithWrongOptionsMessage(subCommandsMap)
	}

	populateCliOptionsOfSubcommand, subcommandExists := subCommandsMap[args[1]]

	if !subcommandExists {
		exitWithWrongOptionsMessage(subCommandsMap)
	}

	populateCliOptionsOfSubcommand()

	return cliOptions
}
