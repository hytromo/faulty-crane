package imagefilters

import (
	"github.com/hytromo/faulty-crane/internal/configuration"
	"github.com/hytromo/faulty-crane/internal/containerregistry"
	"github.com/hytromo/faulty-crane/internal/k8s"
	"github.com/hytromo/faulty-crane/internal/keepreasons"
)

func k8sFilter(repos []containerregistry.Repository, clusters []configuration.KubernetesCluster) {
	if len(clusters) == 0 {
		return
	}

	usedImages := k8s.NewK8s(clusters).GetUsedImages()

	for repoIndex := range repos {
	imageLoop:
		for imageIndex, parsedImage := range repos[repoIndex].Images {
			if parsedImage.KeptData.Reason != keepreasons.None {
				// image already kept for some other reason
				continue
			}

			for _, tag := range parsedImage.Tag {
				fullNameWithTag := parsedImage.Repo + ":" + tag
				cluster, exists := usedImages[fullNameWithTag]

				if exists {
					// image used in a k8s cluster
					repos[repoIndex].Images[imageIndex].KeptData.Reason = keepreasons.UsedInCluster
					repos[repoIndex].Images[imageIndex].KeptData.Metadata = cluster.Context
					continue imageLoop
				}
			}

			for _, digest := range parsedImage.Digest {
				fullNameWithDigest := parsedImage.Repo + "@" + digest
				cluster, exists := usedImages[fullNameWithDigest]

				if exists {
					// image used in a k8s cluster
					repos[repoIndex].Images[imageIndex].KeptData.Reason = keepreasons.UsedInCluster
					repos[repoIndex].Images[imageIndex].KeptData.Metadata = cluster.Context
					continue imageLoop
				}
			}
		}
	}

}
