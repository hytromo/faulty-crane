package argsparser

import "testing"

func TestParse(t *testing.T) {
	cliOptions, err := Parse([]string{"app", "plan", "-out", "plan.out"})

	if err != nil {
		t.Error("Err should be nil")
	}

	if !cliOptions.Apply.SubcommandEnabled {
		t.Error("The apply command should be enabled")
	}

	if cliOptions.ApplyPlanCommon.Plan == "" {
		t.Error("Plan should not be empty")
	}

	if cliOptions.Configure.SubcommandEnabled {
		t.Error("The configure subcommand should not be enabled")
	}

	cliOptions, err = Parse([]string{"app", "apply"})

	if err != nil {
		t.Error("Err should be nil")
	}

	if !cliOptions.Apply.SubcommandEnabled {
		t.Error("The apply command should be enabled")
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

	if cliOptions.Apply.SubcommandEnabled {
		t.Error("The apply subcommand should not be enabled")
	}

	cliOptions, err = Parse([]string{"app"})

	if err == nil {
		t.Error("An error should be returned if the app is called without cli arguments")
	}
}
