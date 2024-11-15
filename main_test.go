package main

import (
	"testing"
)

func TestMakeParams(t *testing.T) {
	keywords := "html javascript css"

	got := MakeParams(keywords)
	want := `?q=%28html%20OR%20javascript%20OR%20css%29`

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
