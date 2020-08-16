package argsparser

import "testing"

func TestParse(t *testing.T) {
	cliOptions, err := Parse([]string{"app", "clean", "-dry-run"})

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

	cliOptions, err = Parse([]string{"app", "clean"})

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

	cliOptions, err = Parse([]string{"app", "non-existent-subcommand"})
	if err == nil {
		t.Error("Non-existent subcommand should return an error")
	}

	cliOptions, err = Parse([]string{"app", "configure", "-o", "file.json"})

	if err != nil {
		t.Error("Err should be nil")
	}

	if !cliOptions.Configure.SubcommandEnabled {
		t.Error("The configure command should be enabled")
	}

	if cliOptions.Configure.Config != "file.json" {
		t.Error("Configuration file has wrong value")
	}

	if cliOptions.Clean.SubcommandEnabled {
		t.Error("The clean subcommand should not be enabled")
	}

	cliOptions, err = Parse([]string{"app"})

	if err == nil {
		t.Error("An error should be returned if the app is called without cli arguments")
	}
}
