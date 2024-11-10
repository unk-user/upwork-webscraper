package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestScanQuery(t *testing.T) {
	t.Run("reads a query from the input", func(t *testing.T) {
		buffer := &bytes.Buffer{}

		got, err := ScanQuery(buffer, strings.NewReader("random string\n"))
		want := "random string"

		assertSuccess(t, err)
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})
	t.Run("writes a prompt to the output", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		_, err := ScanQuery(buffer, strings.NewReader("random string\n"))

		assertSuccess(t, err)
		if buffer.String() != QueryPrompt {
			t.Errorf("expected prompt %q, got %q", QueryPrompt, buffer.String())
		}
	})
}

func TestMakeParams(t *testing.T) {
	keywords := "html javascript css"

	got := MakeParams(keywords)
	want := `?q=%28html%20OR%20javascript%20OR%20css%29`

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func assertSuccess(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
