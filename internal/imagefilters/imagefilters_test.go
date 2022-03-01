package imagefilters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/hytromo/faulty-crane/internal/configuration"
	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
)

func TestParse(t *testing.T) {
	testReposPath := "../../test/repos.json"
	configBytes, err := ioutil.ReadFile(testReposPath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not read configuration file %v: %v\n", testReposPath, err))
	}

	repos := []containerregistry.Repository{}

	err = json.Unmarshal([]byte(configBytes), &repos)

	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid format of configuration file %v: %v\n", testReposPath, err))
	}

	// mutate a repo to whitelist it because it is young
	for i, repo := range repos {
		if repo.Link == "hytromo/whitelistedDueToTime" {
			for j := range repo.Images {
				nowInMs := strconv.FormatInt(time.Now().UnixMilli(), 10)
				repos[i].Images[j].TimeCreatedMs = nowInMs
				repos[i].Images[j].TimeUploadedMs = nowInMs
			}
		}
	}

	parsedRepos := Parse(repos, configuration.KeepImages{
		YoungerThan: "2d",
		AtLeast:     1,
		UsedIn: configuration.UsedIn{
			KubernetesClusters: []configuration.KubernetesCluster{},
		},
		Image: configuration.Image{
			Tags:         []string{"whitelistedTag"},
			Digests:      []string{"sha256:whitelistedDigestwhitelistedDigestwhitelistedDigestwhitelistedge"},
			Repositories: []string{"hytromo/whitelistedRepo"},
		},
	})

	keptCount := 0
	deletedCount := 0

	for _, repo := range parsedRepos {
		for _, image := range repo.Images {
			if repo.Link == "hytromo/whitelistedDueToTime" {
				if image.KeptData.Reason != keepreasons.Young {
					t.Error("Image should be kept because it is young")
				} else {
					keptCount++
				}
			} else if repo.Link == "hytromo/whitelistedRepo" {
				if image.KeptData.Reason != keepreasons.WhitelistedRepository {
					t.Error("Image should be kept due to its whitelisted repository")
				} else {
					keptCount++
				}
			} else if repo.Link == "hytromo/whitelistedDueToOnlyOne1" || repo.Link == "hytromo/whitelistedDueToOnlyOne2" {
				if image.KeptData.Reason != keepreasons.OneOfFew {
					t.Error("Image should be kept because it is the only one")
				} else {
					keptCount++
				}
			}

			if image.Digest[0] == "sha256:whitelistedDigestwhitelistedDigestwhitelistedDigestwhitelistedge" {
				if image.KeptData.Reason != keepreasons.WhitelistedDigest {
					t.Error("Image should be whitelisted due to its digest")
				} else {
					keptCount++
				}
			}

			for _, tag := range image.Tag {
				if tag == "whitelistedTag" {
					if image.KeptData.Reason != keepreasons.WhitelistedTag {
						t.Error("Image should be whitelisted due to its tag")
					} else {
						keptCount++
					}
				}
			}

			if image.KeptData.Reason == keepreasons.None {
				deletedCount++
			}

		}
	}

	if keptCount != 6 {
		t.Errorf("Exactly 6 images should be kept, not %v", keptCount)
	}

	if deletedCount != 2 {
		t.Errorf("Exactly 2 images should be deleted, not %v", deletedCount)
	}
}

func TestParse2(t *testing.T) {
	testReposPath := "../../test/repos.json"
	configBytes, err := ioutil.ReadFile(testReposPath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not read configuration file %v: %v\n", testReposPath, err))
	}

	repos := []containerregistry.Repository{}

	err = json.Unmarshal([]byte(configBytes), &repos)

	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid format of configuration file %v: %v\n", testReposPath, err))
	}

	// mutate a repo to whitelist it because it is young
	for i, repo := range repos {
		if repo.Link == "hytromo/whitelistedDueToTime" {
			for j := range repo.Images {
				nowInMs := strconv.FormatInt(time.Now().UnixMilli(), 10)
				repos[i].Images[j].TimeCreatedMs = nowInMs
				repos[i].Images[j].TimeUploadedMs = nowInMs
			}
		}
	}

	parsedRepos := Parse(repos, configuration.KeepImages{
		YoungerThan: "2d",
		AtLeast:     0,
		UsedIn: configuration.UsedIn{
			KubernetesClusters: []configuration.KubernetesCluster{},
		},
		Image: configuration.Image{
			Tags:         []string{"whitelistedTag"},
			Digests:      []string{"sha256:whitelistedDigestwhitelistedDigestwhitelistedDigestwhitelistedge"},
			Repositories: []string{"hytromo/whitelistedRepo"},
		},
	})

	keptCount := 0
	deletedCount := 0

	for _, repo := range parsedRepos {
		for _, image := range repo.Images {
			if repo.Link == "hytromo/whitelistedDueToTime" {
				if image.KeptData.Reason != keepreasons.Young {
					t.Error("Image should be kept because it is young")
				} else {
					keptCount++
				}
			} else if repo.Link == "hytromo/whitelistedRepo" {
				if image.KeptData.Reason != keepreasons.WhitelistedRepository {
					t.Error("Image should be kept due to its whitelisted repository")
				} else {
					keptCount++
				}
			} else if repo.Link == "hytromo/whitelistedDueToOnlyOne1" || repo.Link == "hytromo/whitelistedDueToOnlyOne2" {
				if image.KeptData.Reason != keepreasons.None {
					t.Error("Should not be whitelisted actually as at least is 0 in this test")
				}
			}

			if image.Digest[0] == "sha256:whitelistedDigestwhitelistedDigestwhitelistedDigestwhitelistedge" {
				if image.KeptData.Reason != keepreasons.WhitelistedDigest {
					t.Error("Image should be whitelisted due to its digest")
				} else {
					keptCount++
				}
			}

			for _, tag := range image.Tag {
				if tag == "whitelistedTag" {
					if image.KeptData.Reason != keepreasons.WhitelistedTag {
						t.Error("Image should be whitelisted due to its tag")
					} else {
						keptCount++
					}
				}
			}

			if image.KeptData.Reason == keepreasons.None {
				deletedCount++
			}

		}
	}

	if keptCount != 4 {
		t.Errorf("Exactly 4 images should be kept, not %v", keptCount)
	}

	if deletedCount != 4 {
		t.Errorf("Exactly 4 images should be deleted, not %v", deletedCount)
	}
}

func TestParse3(t *testing.T) {
	testReposPath := "../../test/repos.json"
	configBytes, err := ioutil.ReadFile(testReposPath)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not read configuration file %v: %v\n", testReposPath, err))
	}

	repos := []containerregistry.Repository{}

	err = json.Unmarshal([]byte(configBytes), &repos)

	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid format of configuration file %v: %v\n", testReposPath, err))
	}

	// mutate a repo to whitelist it because it is young
	for i, repo := range repos {
		if repo.Link == "hytromo/whitelistedDueToTime" {
			for j := range repo.Images {
				nowInMs := strconv.FormatInt(time.Now().UnixMilli(), 10)
				repos[i].Images[j].TimeCreatedMs = nowInMs
				repos[i].Images[j].TimeUploadedMs = nowInMs
			}
		}
	}

	parsedRepos := Parse(repos, configuration.KeepImages{
		YoungerThan: "2d",
		AtLeast:     10,
		UsedIn: configuration.UsedIn{
			KubernetesClusters: []configuration.KubernetesCluster{},
		},
		Image: configuration.Image{
			Tags:         []string{"whitelistedTag"},
			Digests:      []string{"sha256:whitelistedDigestwhitelistedDigestwhitelistedDigestwhitelistedge"},
			Repositories: []string{"hytromo/whitelistedRepo"},
		},
	})

	keptCount := 0
	deletedCount := 0

	for _, repo := range parsedRepos {
		for _, image := range repo.Images {
			if image.KeptData.Reason == keepreasons.None {
				deletedCount++
			} else {
				keptCount++
			}
		}
	}

	if keptCount != 8 {
		t.Errorf("Exactly 8 images should be kept, not %v", keptCount)
	}

	if deletedCount != 0 {
		t.Errorf("Exactly 0 images should be deleted, not %v", deletedCount)
	}
}
