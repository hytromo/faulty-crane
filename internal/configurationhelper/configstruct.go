package configurationhelper

type kubernetesCluster struct {
	Context string
}

type containerRegistry struct {
	Type string
}

type usedIn struct {
	KubernetesClusters []kubernetesCluster
}

type image struct {
	Tags         []string
	Digests      []string
	IDs          []string
	Repositories []string
}

type keep struct {
	YoungerThan string
	UsedIn      usedIn
	Image       image
}

// Configuration struct shows the structure of the configuration file used by this app
type Configuration struct {
	ContainerRegistry containerRegistry
	Keep              keep
}
