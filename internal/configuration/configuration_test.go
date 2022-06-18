package configuration

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func writeDockerhubAnswers(writer *io.PipeWriter) {
	defer writer.Close()
	answers := []string{"dockerhub", "hytromo", "namespace", "1234", "10d", "k1", "k2", "", "t1", "t2", "", "d1", "d2", "", "i1", "i2", ""}
	for _, answer := range answers {
		_, err := io.WriteString(writer, answer+"\r\n")
		if err != nil {
			log.Fatal("Could not write dockerhub answers")
		}
		time.Sleep(time.Duration(time.Microsecond * 100))
	}
}

func TestAskUserInput(t *testing.T) {
	reader, writer := io.Pipe()

	go writeDockerhubAnswers(writer)

	userInput := AskUserInput(reader)

	if userInput.ContainerRegistryUsername != "hytromo" {
		t.Error("Wrong username")
	}

	if userInput.ContainerRegistryNamespace != "namespace" {
		t.Error("Wrong namespace")
	}

	if userInput.ContainerRegistryPassword != "1234" {
		t.Error("Wrong password")
	}

	if userInput.YoungerThan != "10d" {
		t.Error("Wrong younger than")
	}

	if !reflect.DeepEqual(userInput.KubernetesClusters, []string{"k1", "k2"}) {
		t.Error("Wrong clusters")
	}

	if !reflect.DeepEqual(userInput.ImageTags, []string{"t1", "t2"}) {
		t.Error("Wrong image tags")
	}

	if !reflect.DeepEqual(userInput.ImageDigests, []string{"d1", "d2"}) {
		t.Error("Wrong image digests")
	}

	if !reflect.DeepEqual(userInput.ImageIDs, []string{"i1", "i2"}) {
		t.Error("Wrong image ids")
	}

}

func TestCreateNew(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "unit_test")

	if err != nil {
		t.Error("Could not create temp file")
	}

	defer os.Remove(tmpFile.Name())

	reader, writer := io.Pipe()

	go writeDockerhubAnswers(writer)

	CreateNew(ConfigureSubcommandOptions{
		SubcommandEnabled: true,
		Config:            tmpFile.Name(),
	}, reader)
}

func TestPlanRW(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "unit_test")

	if err != nil {
		t.Error("Could not create temp file")
	}

	defer os.Remove(tmpFile.Name())

	plan := ReadPlan("../../test/repos.json", false)

	if len(plan) <= 0 {
		// just a failsafe to ensure we have read something
		t.Error("Not enough repos in the plan to test")
	}

	WritePlan(plan, tmpFile.Name(), true)
	newPlan := ReadPlan(tmpFile.Name(), true)

	if !reflect.DeepEqual(plan, newPlan) {
		t.Error("Written plan does not equal initial plan")
	}
}

func TestRegistryDetection(t *testing.T) {
	if !IsGCR(&AppOptions{
		ApplyPlanCommon: ApplyPlanCommonSubcommandOptions{
			GoogleContainerRegistry: GoogleContainerRegistry{
				Host:  "asia.gcr.io",
				Token: "",
			},
		},
	}) {

		t.Error("Should be detected as GCR")
	}

	if !IsDockerhub(&AppOptions{
		ApplyPlanCommon: ApplyPlanCommonSubcommandOptions{
			DockerhubContainerRegistry: DockerhubContainerRegistry{
				Username:  "test",
				Namespace: "test",
				Password:  "test",
			},
		},
	}) {
		t.Error("Should be detected as dockerhub")
	}
}
