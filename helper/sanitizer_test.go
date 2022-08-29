package helper

import "testing"

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
