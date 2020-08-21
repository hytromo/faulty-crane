package k8s

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/hytromo/faulty-crane/internal/configuration"
	log "github.com/sirupsen/logrus"
	apiAppsV1 "k8s.io/api/apps/v1"
	apiCoreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ClusterWithAPI is a struct that contains information about the cluster as well as api clients
type ClusterWithAPI struct {
	Context       string
	Namespace     string
	RunningInside bool // RunningInside means that this app is running inside this cluster (and thus different configuration options need to be specified)
	CoreV1        coreV1.CoreV1Interface
	AppsV1        appsV1.AppsV1Interface
}

// K8s struct provides an object that fetches resources from multiple k8s clusters
type K8s struct {
	Clusters []ClusterWithAPI
}

type podsContainer struct {
	ClusterWithAPI *ClusterWithAPI
	PodList        *apiCoreV1.PodList
}

type deploymentsContainer struct {
	ClusterWithAPI *ClusterWithAPI
	DeploymentList *apiAppsV1.DeploymentList
}

type replicaSetsContainer struct {
	ClusterWithAPI *ClusterWithAPI
	ReplicaSetList *apiAppsV1.ReplicaSetList
}

type statefulSetsContainer struct {
	ClusterWithAPI  *ClusterWithAPI
	StatefulSetList *apiAppsV1.StatefulSetList
}

func buildConfigFromFlags(kubectlContext, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: kubectlContext,
		}).ClientConfig()
}

// NewK8s connects to the k8s cluster
func NewK8s(clusters []configuration.KubernetesCluster) K8s {
	// TODO: dynamic kubeconfig connection
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	clustersWithAPI := make([]ClusterWithAPI, len(clusters))
	for clusterIndex, cluster := range clusters {
		var config *rest.Config

		if cluster.RunningInside {
			// this app is running inside this cluster
			internalConfig, err := clientcmd.BuildConfigFromFlags("", "")

			if err != nil {
				log.Fatal("Could not parse kubeconfig: %v", err.Error())
			}

			config = internalConfig
		} else {
			// external cluster
			externalConfig, err := buildConfigFromFlags(cluster.Context, kubeconfig)

			if err != nil {
				log.Fatal("Could not parse kubeconfig: %v", err.Error())
			}

			config = externalConfig
		}

		clientset, err := kubernetes.NewForConfig(config)

		if err != nil {
			log.Fatal("Cannot initialize kubernetes client with this config: %v", err.Error())
		}

		clustersWithAPI[clusterIndex] = ClusterWithAPI{
			Namespace: cluster.Namespace,
			Context:   cluster.Context,
			CoreV1:    clientset.CoreV1(),
			AppsV1:    clientset.AppsV1(),
		}
	}

	return K8s{Clusters: clustersWithAPI}
}

func (k8s *K8s) getPods(waitGroup *sync.WaitGroup, podsChan chan<- podsContainer) {
	defer waitGroup.Done()

	for index, cluster := range k8s.Clusters {
		pods, err := cluster.CoreV1.Pods(cluster.Namespace).List(metav1.ListOptions{})

		if err != nil {
			log.Fatal("Could not read pods: %v", err.Error())
		}

		podsChan <- podsContainer{
			ClusterWithAPI: &k8s.Clusters[index],
			PodList:        pods,
		}
	}

	close(podsChan)
}

func (k8s *K8s) getDeployments(waitGroup *sync.WaitGroup, deploymentsChan chan<- deploymentsContainer) {
	defer waitGroup.Done()

	for index, cluster := range k8s.Clusters {
		deployments, err := cluster.AppsV1.Deployments(cluster.Namespace).List(metav1.ListOptions{})

		if err != nil {
			log.Fatal("Could not read deployments: %v", err.Error())
		}

		deploymentsChan <- deploymentsContainer{
			ClusterWithAPI: &k8s.Clusters[index],
			DeploymentList: deployments,
		}
	}

	close(deploymentsChan)
}

func (k8s *K8s) getReplicaSets(waitGroup *sync.WaitGroup, replicaSetsChan chan<- replicaSetsContainer) {
	defer waitGroup.Done()

	for index, cluster := range k8s.Clusters {
		replicaSets, err := cluster.AppsV1.ReplicaSets(cluster.Namespace).List(metav1.ListOptions{})

		if err != nil {
			log.Fatal("Could not read replicaSets: %v", err.Error())
		}

		replicaSetsChan <- replicaSetsContainer{
			ClusterWithAPI: &k8s.Clusters[index],
			ReplicaSetList: replicaSets,
		}
	}

	close(replicaSetsChan)
}

func (k8s *K8s) getStatefulSets(waitGroup *sync.WaitGroup, statefulSetsChan chan<- statefulSetsContainer) {
	defer waitGroup.Done()

	for index, cluster := range k8s.Clusters {
		statefulSets, err := cluster.AppsV1.StatefulSets(cluster.Namespace).List(metav1.ListOptions{})

		if err != nil {
			log.Fatal("Could not read statefulSets: %v", err.Error())
		}

		statefulSetsChan <- statefulSetsContainer{
			ClusterWithAPI:  &k8s.Clusters[index],
			StatefulSetList: statefulSets,
		}
	}

	close(statefulSetsChan)
}

// GetUsedImages gets all the images used inside a kubernetes cluster by fetching the corresponding resources concurrently
func (k8s K8s) GetUsedImages() map[string]*ClusterWithAPI {
	waitGroup := sync.WaitGroup{}

	// TODO: can we make a channel that never blocks on range read? E.g. it reads whatever is there and then unblocks
	podsChan := make(chan podsContainer, len(k8s.Clusters))
	deploymentsChan := make(chan deploymentsContainer, len(k8s.Clusters))
	replicaSetsChan := make(chan replicaSetsContainer, len(k8s.Clusters))
	statefulSetsChan := make(chan statefulSetsContainer, len(k8s.Clusters))

	log.Info("Reading kubernetes cluster(s)...")

	waitGroup.Add(4)

	go k8s.getPods(&waitGroup, podsChan)
	go k8s.getDeployments(&waitGroup, deploymentsChan)
	go k8s.getReplicaSets(&waitGroup, replicaSetsChan)
	go k8s.getStatefulSets(&waitGroup, statefulSetsChan)

	waitGroup.Wait()

	images := map[string]*ClusterWithAPI{}

	for pods := range podsChan {
		extractImagesFromPods(pods, &images)
	}

	for deployments := range deploymentsChan {
		extractImagesFromDeployments(deployments, &images)
	}

	for replicaSets := range replicaSetsChan {
		extractImagesFromReplicaSets(replicaSets, &images)
	}

	for statefulSets := range statefulSetsChan {
		extractImagesFromStatefulSets(statefulSets, &images)
	}

	return images
}
