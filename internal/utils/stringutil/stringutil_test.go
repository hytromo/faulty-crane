package stringutil

import (
	"testing"
)

func TestStrInSlice(t *testing.T) {
	if !StrInSlice("exists", []string{"this", "indeed", "exists"}) {
		t.Error("Wrong assertion")
	}

	if StrInSlice("does-not-exist", []string{"this", "does", "not", "exist"}) {
		t.Error("Wrong assertion")
	}

	if StrInSlice("does-not-exist", []string{}) {
		t.Error("Wrong assertion")
	}
}

func TestTrimLeftChars(t *testing.T) {
	if TrimLeftChars("abccname:", 4) != "name:" {
		t.Error("Wrong trim")
	}

	if TrimLeftChars("abccname:", 8) != ":" {
		t.Error("Wrong trim")
	}
}

func TestTrimRightChars(t *testing.T) {
	if TrimRightChars("abccname:", 1) != "abccname" {
		t.Error("Wrong trim")
	}

	if TrimRightChars("abccname:", 8) != "a" {
		t.Error("Wrong trim")
	}
}

func TestKeepAtMost(t *testing.T) {
	if KeepAtMost("abcde", 4) != "ab.." {
		t.Error("Wrong keep at most")
	}

	if KeepAtMost("abcde", 3) != "a.." {
		t.Error("Wrong keep at most")
	}

	if KeepAtMost("a", 3) != "a" {
		t.Error("Wrong keep at most")
	}

	if KeepAtMost("My name is alex!", 15) != "My name is al.." {
		t.Error("Wrong keep at most")
	}
}

func TestHumanFriendlySize(t *testing.T) {
	friendlyToSize := map[string]int64{
		"2.0 kiB":  2048,
		"3.0 MiB":  1024 * 1024 * 3,
		"11.0 GiB": 1024 * 1024 * 1024 * 11,
		"11.5 GiB": 1024 * 1024 * 1024 * 11.5,
		"12.5 TiB": 1024 * 1024 * 1024 * 1024 * 12.5,
		"123 B":    123,
	}

	for friendlySize, bytes := range friendlyToSize {
		if HumanFriendlySize(bytes) != friendlySize {
			t.Errorf("Wrong human friendly size %v", HumanFriendlySize(bytes))
		}
	}
}
