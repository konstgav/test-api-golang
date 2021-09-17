package main

import (
	"net/http"
	"testing"
)

type TestWriter struct {
	http.ResponseWriter
	Text string
}

func (tw *TestWriter) Write(data []byte) (int, error) {
	tw.Text = string(data)
	return 0, nil
}

func TestHomePage(t *testing.T) {
	w := &TestWriter{}
	HomePage(w, nil)
	got := w.Text
	want := "Welcome to the test CRUD API!"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
