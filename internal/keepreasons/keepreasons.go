package keepreasons

// KeptReason enum represents the reason why an image was kept e.g. not cleaned
type KeptReason int

const (
	// None kept reason means that the image does not have a reason to be kept and thus it WILL be deleted, it is important that this has value 0, as non-parsed images have automatically kept reason 0
	None KeptReason = iota
	// Young kept reason means that the image was uploaded recently and thus is not filtered
	Young
	// UsedInCluster kept reason means that the image is being used in a k8s cluster and thus will not be deleted
	UsedInCluster
	// WhitelistedTag kept reason means that the image has a tag which is whitelisted and thus will not be deleted
	WhitelistedTag
	// WhitelistedDigest kept reason means that the image has a digest which is whitelisted and thus will not be deleted
	WhitelistedDigest
	// WhitelistedRepository kept reason means that the image has a repository which is whitelisted and thus will not be deleted
	WhitelistedRepository
	// OneOfFew kept reason means that the repository needs to keep a minimum number of images, that's why the image was kept
	OneOfFew
)

// KeptData contains all the data needed to figure out why an image was kept from being deleted
type KeptData struct {
	Reason KeptReason
	// Metadata contains extra data about the reason, e.g. if the image is kept because it is used in a k8s cluster, this may contain the cluster context
	Metadata string
}
