package configuration

// KubernetesCluster encapsulates all the information needed per kubernetes cluster
type KubernetesCluster struct {
	Context string
}

// ContainerRegistry keeps the needed data for the container registry
type ContainerRegistry struct {
	Host   string
	Access string
}

// UsedIn defines a list of resources that could use container images
type UsedIn struct {
	KubernetesClusters []KubernetesCluster
}

// Image defines various image-related fields
type Image struct {
	Tags         []string
	Digests      []string
	Repositories []string
}

// KeepImages specifies what conditions we should use in order to keep images from being deleted
type KeepImages struct {
	YoungerThan string
	UsedIn      UsedIn
	Image       Image
}

// Configuration struct shows the structure of the configuration file used by this app
type Configuration struct {
	ContainerRegistry ContainerRegistry
	Keep              KeepImages
}

// CleanSubcommandOptions defines the options of the clean subcommand
type CleanSubcommandOptions struct {
	SubcommandEnabled bool
	DryRun            bool
	// Config is where to read configuration from
	Config            string
	ContainerRegistry ContainerRegistry
	Keep              KeepImages
}

// ConfigureSubcommandOptions defines the options of the configure subcommand
type ConfigureSubcommandOptions struct {
	SubcommandEnabled bool
	// Config is where to save the configuration file
	Config string
}

// AppOptions groups all the possible application options in a single struct
type AppOptions struct {
	Clean     CleanSubcommandOptions
	Configure ConfigureSubcommandOptions
}
