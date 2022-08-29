package helper

import (
	"reflect"
	"strings"
	"testing"
)

func TestSanitize(t *testing.T) {
	ours := "abc"
	theirs := Sanitize("abc\n\n\n")

	if ours != theirs {
		t.Errorf("Got %s but wanted %s", theirs, ours)
	}
}

func TestSanitizeWithNoNewLines(t *testing.T) {
	ours := "abc"
	theirs := Sanitize("abc")

	if ours != theirs {
		t.Errorf("Got %s but wanted %s", theirs, ours)
	}
}

func TestSliceSanitize(t *testing.T) {
	ours := strings.Split("abc aaa def", " ")

	theirs := strings.Split("abc aaa\n def\n\n", " ")
	SanitizeStringSlice(theirs)

	if !reflect.DeepEqual(ours, theirs) {
		t.Errorf("Got %v but wanted %v", theirs, ours)
	}
}

func TestSliceWithNoNewLinesSanitize(t *testing.T) {
	ours := strings.Split("abc aaa def", " ")

	theirs := strings.Split("abc aaa def", " ")
	SanitizeStringSlice(theirs)

	if !reflect.DeepEqual(ours, theirs) {
		t.Errorf("Got %v but wanted %v", theirs, ours)
	}
}
