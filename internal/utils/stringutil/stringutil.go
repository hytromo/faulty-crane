package stringutil

// StrInSlice returns whether a string exists inside a slice
func StrInSlice(strToSearch string, slice []string) bool {
	for _, strOfSlice := range slice {
		if strOfSlice == strToSearch {
			return true
		}
	}

	return false
}
