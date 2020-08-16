package argsparser

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type CleanCliOptions struct {
	SubcommandEnabled bool
	DryRun            bool
}

type ConfigureCliOptions struct {
	SubcommandEnabled bool
	Config            string
}

type CliOptions struct {
	Clean     CleanCliOptions
	Configure ConfigureCliOptions
}

func getWrongOptionsError(subCommandsMap map[string]func()) (err error) {
	allSubcommands := make([]string, len(subCommandsMap))

	i := 0
	for k := range subCommandsMap {
		allSubcommands[i] = k
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
func Parse(args []string) (CliOptions, error) {
	cleanSubCmd := "clean"
	configureSubCmd := "configure"

	var cliOptions CliOptions

	subCommandsMap := map[string]func(){
		cleanSubCmd: func() {
			cliOptions.Clean.SubcommandEnabled = true
			cleanCmd := flag.NewFlagSet(cleanSubCmd, flag.ExitOnError)
			cleanCmd.BoolVar(&cliOptions.Clean.DryRun, "dry-run", false, "just output what is expected to be deleted")
			cleanCmd.Parse(args[2:])
		},
		configureSubCmd: func() {
			cliOptions.Configure.SubcommandEnabled = true
			configureCmd := flag.NewFlagSet(configureSubCmd, flag.ExitOnError)
			configureCmd.StringVar(&cliOptions.Configure.Config, "o", "faulty-crane.json", "the file to save the configuration to")
			configureCmd.Parse(args[2:])
		},
	}

	if len(args) < 2 {
		return cliOptions, getWrongOptionsError(subCommandsMap)
	}

	populateCliOptionsOfSubcommand, subcommandExists := subCommandsMap[args[1]]

	if !subcommandExists {
		return cliOptions, getWrongOptionsError(subCommandsMap)
	}

	populateCliOptionsOfSubcommand()

	return cliOptions, nil
}
