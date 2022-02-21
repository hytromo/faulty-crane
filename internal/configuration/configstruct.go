package configuration

// KubernetesCluster encapsulates all the information needed per kubernetes cluster
type KubernetesCluster struct {
	Context       string
	Namespace     string
	RunningInside bool // RunningInside means that faulty-crane is running inside this cluster and thus the k8s client needs specific options to communicate with this cluster
}

// GoogleContainerRegistry keeps the needed data for the google container registry
type GoogleContainerRegistry struct {
	Host  string
	Token string
}

// DockerhubContainerRegistry keeps the needed data for the google container registry
type DockerhubContainerRegistry struct {
	Username string
	Password string
	// Namespace is where do you want us to search for images, could be same as the username, could be an org name etc
	Namespace string
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
	// Keep images younger than e.g. 5d
	YoungerThan string
	// Keep at least N images
	AtLeast int
	// Keep the images used in the below contexts
	UsedIn UsedIn
	// Keep images with the below image-related characteristics
	Image Image
}

// Configuration struct shows the structure of the configuration file used by this app
type Configuration struct {
	GCR       GoogleContainerRegistry    `json:",omitempty"`
	Dockerhub DockerhubContainerRegistry `json:",omitempty"`
	Keep      KeepImages
}

// ApplySubcommandOptions defines the options of the apply subcommand
type ApplySubcommandOptions struct {
	SubcommandEnabled bool
}

// PlanSubcommandOptions defines the options of the plan subcommand
type PlanSubcommandOptions struct {
	SubcommandEnabled bool
}

// ApplyPlanCommonSubcommandOptions defines the common options between plan and apply subcommands
type ApplyPlanCommonSubcommandOptions struct {
	// Plan file to write, or to read from for deleting images
	Plan string
	// Config is the path of the configuration file
	Config                     string
	GoogleContainerRegistry    GoogleContainerRegistry
	DockerhubContainerRegistry DockerhubContainerRegistry
	Keep                       KeepImages
}

// ConfigureSubcommandOptions defines the options of the configure subcommand
type ConfigureSubcommandOptions struct {
	SubcommandEnabled bool
	// Config is where to save the configuration file
	Config string
}

// ShowSubcommandOptions defines the options of the show subcommand
type ShowSubcommandOptions struct {
	SubcommandEnabled bool
	// Plan is the path to the plan file to show
	Plan       string
	Analytical bool
}

// AppOptions groups all the possible application options in a single struct
type AppOptions struct {
	Configure       ConfigureSubcommandOptions
	Plan            PlanSubcommandOptions
	Show            ShowSubcommandOptions
	Apply           ApplySubcommandOptions
	ApplyPlanCommon ApplyPlanCommonSubcommandOptions
}
