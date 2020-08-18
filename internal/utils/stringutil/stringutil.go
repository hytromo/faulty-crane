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

// TrimLeftChars trims the from the left part of a string as many chars as specified
func TrimLeftChars(strToTrim string, numberOfCharsToRemove int) string {
	removedTillNow := 0
	for i := range strToTrim {
		if removedTillNow >= numberOfCharsToRemove {
			return strToTrim[i:]
		}
		removedTillNow++
	}
	return strToTrim[:0]
}
