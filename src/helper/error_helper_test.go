package helper

import (
	"errors"
	"testing"
)

func TestNewError(t *testing.T) {

	got := NewError("test error")
	want := errors.New("test error")

	if got.Error() != want.Error() {
		t.Errorf("Got %q, wanted %q", got, want)
	}
}

func TestNewErrorWithFormattedMsg(t *testing.T) {

	got := NewError("test error: [%s]", "abc")
	want := errors.New("test error: [abc]")

	if got.Error() != want.Error() {
		t.Errorf("Got %q, wanted %q", got, want)
	}
}
