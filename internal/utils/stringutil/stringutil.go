package stringutil

import (
	"fmt"
	"strings"
)

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

// TrimRightChars trims the from the right part of a string as many chars as specified
func TrimRightChars(strToTrim string, numberOfCharsToRemove int) string {
	return strToTrim[:len(strToTrim)-numberOfCharsToRemove]
}

// KeepAtMost shortens a string to at most numberOfCharsToKeep characters, including the dots at the end
func KeepAtMost(strToShorten string, numberOfCharsToKeep int) string {
	if len(strToShorten) <= numberOfCharsToKeep {
		return strToShorten
	}

	numOfDots := 2

	return TrimRightChars(strToShorten, len(strToShorten)-numberOfCharsToKeep+numOfDots) + strings.Repeat(".", numOfDots)
}

// HumanFriendlySize returns a human friendly representation of bytes count (IEC)
func HumanFriendlySize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "kMGTPE"[exp])
}
