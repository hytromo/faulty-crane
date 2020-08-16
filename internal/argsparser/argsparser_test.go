package argsparser

import "testing"

func TestParse(t *testing.T) {
	cliOptions, err := Parse([]string{"program", "clean", "-dry-run"})

	if err != nil {
		t.Error("Err should be nil")
	}

	if !cliOptions.Clean.SubcommandEnabled {
		t.Error("The clean command should be enabled")
	}

	if !cliOptions.Clean.DryRun {
		t.Error("Dry run option should be enabled")
	}

	if cliOptions.Configure.SubcommandEnabled {
		t.Error("The configure subcommand should not be enabled")
	}

	cliOptions, err = Parse([]string{"program", "clean"})

	if err != nil {
		t.Error("Err should be nil")
	}

	if !cliOptions.Clean.SubcommandEnabled {
		t.Error("The clean command should be enabled")
	}

	if cliOptions.Clean.DryRun {
		t.Error("Dry run option should not be enabled")
	}

	if cliOptions.Configure.SubcommandEnabled {
		t.Error("The configure subcommand should not be enabled")
	}

	cliOptions, err = Parse([]string{"program", "non-existent-subcommand"})
	if err == nil {
		t.Error("Non-existent subcommand should return an error")
	}
}
