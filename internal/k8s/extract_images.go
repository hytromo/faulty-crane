package k8s

import (
	apiCoreV1 "k8s.io/api/core/v1"
)

func extractImagesFromContainers(containers []apiCoreV1.Container) []string {
	images := []string{}

	for _, container := range containers {
		images = append(images, container.Image)
	}

	return images
}

func extractImagesFromSpec(spec apiCoreV1.PodSpec) []string {
	images := []string{}

	images = append(images, extractImagesFromContainers(spec.InitContainers)...)
	images = append(images, extractImagesFromContainers(spec.Containers)...)

	return images
}

func extractImagesFromPods(pods podsContainer, images *map[string]*ClusterWithAPI) {
	for _, pod := range pods.PodList.Items {
		for _, image := range extractImagesFromSpec(pod.Spec) {
			(*images)[image] = pods.ClusterWithAPI
		}
	}
}

func extractImagesFromJobs(jobs jobsContainer, images *map[string]*ClusterWithAPI) {
	for _, job := range jobs.JobList.Items {
		for _, image := range extractImagesFromSpec(job.Spec.Template.Spec) {
			(*images)[image] = jobs.ClusterWithAPI
		}
	}
}

func extractImagesFromCronJobs(cronJobs cronJobsContainer, images *map[string]*ClusterWithAPI) {
	for _, job := range cronJobs.CronJobList.Items {
		for _, image := range extractImagesFromSpec(job.Spec.JobTemplate.Spec.Template.Spec) {
			(*images)[image] = cronJobs.ClusterWithAPI
		}
	}
}

func extractImagesFromDeployments(deployments deploymentsContainer, images *map[string]*ClusterWithAPI) {
	for _, deployment := range deployments.DeploymentList.Items {
		for _, image := range extractImagesFromSpec(deployment.Spec.Template.Spec) {
			(*images)[image] = deployments.ClusterWithAPI
		}
	}
}

func extractImagesFromReplicaSets(replicaSets replicaSetsContainer, images *map[string]*ClusterWithAPI) {
	for _, replicaSet := range replicaSets.ReplicaSetList.Items {
		for _, image := range extractImagesFromSpec(replicaSet.Spec.Template.Spec) {
			(*images)[image] = replicaSets.ClusterWithAPI
		}
	}
}

func extractImagesFromStatefulSets(statefulSets statefulSetsContainer, images *map[string]*ClusterWithAPI) {
	for _, statefulSet := range statefulSets.StatefulSetList.Items {
		for _, image := range extractImagesFromSpec(statefulSet.Spec.Template.Spec) {
			(*images)[image] = statefulSets.ClusterWithAPI
		}
	}
}
