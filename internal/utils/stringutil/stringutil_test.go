package stringutil

import "testing"

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
