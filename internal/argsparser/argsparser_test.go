package argsparser

import (
	"os"
	"reflect"
	"testing"
)

func TestPlan(t *testing.T) {
	cliOptions, err := Parse([]string{"app", "plan",
		"-config", "../../test/config.json",
		"-out", "plan.out",
		"-keep-at-least", "10",
		"-keep-image-digests", "d1,d2",
		"-keep-image-repos", "r1,r2",
		"-keep-image-tags", "t1,t2",
		"-keep-used-in-k8s", "k1,k2",
		"-keep-younger-than", "1w3d",
		"-username", "user",
		"-password", "pass",
	})

	if err != nil {
		t.Error("Err should be nil")
	}

	if !cliOptions.Plan.SubcommandEnabled {
		t.Error("The plan command should be enabled")
	}

	if cliOptions.Configure.SubcommandEnabled {
		t.Error("The configure subcommand should not be enabled")
	}

	if cliOptions.Show.SubcommandEnabled {
		t.Error("The show subcommand should not be enabled")
	}

	if cliOptions.ApplyPlanCommon.Plan != "plan.out" {
		t.Error("Plan should have a specific value")
	}

	if cliOptions.ApplyPlanCommon.Config != "../../test/config.json" {
		t.Error("Config should have a specific value")
	}

	if cliOptions.ApplyPlanCommon.DockerhubContainerRegistry.Username != "user" {
		t.Error("Dockerhub username should be user")
	}

	if cliOptions.ApplyPlanCommon.DockerhubContainerRegistry.Password != "pass" {
		t.Error("Dockerhub passname should be pass")
	}

	if cliOptions.ApplyPlanCommon.DockerhubContainerRegistry.Namespace != "namespace" {
		t.Error("Dockerhub namespace should be namespace")
	}

	if !reflect.DeepEqual(cliOptions.ApplyPlanCommon.Keep.Image.Digests, []string{"d1", "d2"}) {
		t.Error("Wrong keep image digests")
	}

	if !reflect.DeepEqual(cliOptions.ApplyPlanCommon.Keep.Image.Repositories, []string{"r1", "r2"}) {
		t.Error("Wrong keep image repos")
	}

	if !reflect.DeepEqual(cliOptions.ApplyPlanCommon.Keep.Image.Tags, []string{"t1", "t2"}) {
		t.Error("Wrong keep image tags")
	}

	if cliOptions.ApplyPlanCommon.Keep.AtLeast != 10 {
		t.Error("Wrong keep at least")
	}

	k8sContexts := []string{}

	for _, cluster := range cliOptions.ApplyPlanCommon.Keep.UsedIn.KubernetesClusters {
		k8sContexts = append(k8sContexts, cluster.Context)
	}

	if !reflect.DeepEqual(k8sContexts, []string{"k1", "k2"}) {
		t.Error("Wrong keep used in k8s contexts")
	}

	if cliOptions.ApplyPlanCommon.Keep.YoungerThan != "1w3d" {
		t.Error("Wrong keep younger than")
	}

}

// TestOptionsOrder tests whether the order that the options are evaluated is config < env < cli options
func TestOptionsOrder(t *testing.T) {
	os.Setenv(EnvPrefix+"CONTAINER_REGISTRY_USERNAME", "envOverridenUsername")
	os.Setenv(EnvPrefix+"KEEP_USED_IN_K8S", "k3,k4")

	cliOptions, err := Parse([]string{"app", "plan",
		"-config", "../../test/config.json",
		"-out", "plan.out",
		"-keep-at-least", "10",
		"-keep-image-digests", "d1,d2",
		"-keep-image-repos", "r1,r2",
		"-keep-image-tags", "t1,t2",
		"-keep-used-in-k8s", "k1,k2",
		"-keep-younger-than", "1w3d",
		"-password", "pass",
	})

	if err != nil {
		t.Error("Err should be nil")
	}

	if !cliOptions.Plan.SubcommandEnabled {
		t.Error("The plan command should be enabled")
	}

	if cliOptions.Configure.SubcommandEnabled {
		t.Error("The configure subcommand should not be enabled")
	}

	if cliOptions.Show.SubcommandEnabled {
		t.Error("The show subcommand should not be enabled")
	}

	if cliOptions.ApplyPlanCommon.Plan != "plan.out" {
		t.Error("Plan should have a specific value")
	}

	if cliOptions.ApplyPlanCommon.Config != "../../test/config.json" {
		t.Error("Config should have a specific value")
	}

	if cliOptions.ApplyPlanCommon.DockerhubContainerRegistry.Username != "envOverridenUsername" {
		t.Error("Dockerhub username should be envOverridenUsername")
	}

	if cliOptions.ApplyPlanCommon.DockerhubContainerRegistry.Password != "pass" {
		t.Error("Dockerhub passname should be pass")
	}

	if cliOptions.ApplyPlanCommon.DockerhubContainerRegistry.Namespace != "namespace" {
		t.Error("Dockerhub namespace should be namespace")
	}

	if !reflect.DeepEqual(cliOptions.ApplyPlanCommon.Keep.Image.Digests, []string{"d1", "d2"}) {
		t.Error("Wrong keep image digests")
	}

	if !reflect.DeepEqual(cliOptions.ApplyPlanCommon.Keep.Image.Repositories, []string{"r1", "r2"}) {
		t.Error("Wrong keep image repos")
	}

	if !reflect.DeepEqual(cliOptions.ApplyPlanCommon.Keep.Image.Tags, []string{"t1", "t2"}) {
		t.Error("Wrong keep image tags")
	}

	if cliOptions.ApplyPlanCommon.Keep.AtLeast != 10 {
		t.Error("Wrong keep at least")
	}

	k8sContexts := []string{}

	for _, cluster := range cliOptions.ApplyPlanCommon.Keep.UsedIn.KubernetesClusters {
		k8sContexts = append(k8sContexts, cluster.Context)
	}

	if !reflect.DeepEqual(k8sContexts, []string{"k1", "k2"}) {
		t.Error("Wrong keep used in k8s contexts")
	}

	if cliOptions.ApplyPlanCommon.Keep.YoungerThan != "1w3d" {
		t.Error("Wrong keep younger than")
	}

}
