package lbase

import (
	"testing"
)

func TestStub(t *testing.T) {

	db, err := Open("test=1 global=1")
	if err != nil || db == nil {
		t.Fatal("Open test global failed")
	}
	defer Close()

	tFunc(t)
}
