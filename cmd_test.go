package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Helper function to capture output
func captureOutput(f func()) string {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	os.Stdout = stdout
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

const ExpectedHelpMessage = `wki - Wikipedia CLI
usage: wki [-h] [-s SEARCH]

options:
  -h, --help            show this help message
  -s SEARCH, --search SEARCH
			search for topic <SEARCH>
`

var cmdtests = []struct {
	in  string
	out string
}{
	{"garbage", ExpectedHelpMessage},
	{"garbage lol 3", ""},
	{"", ""},
}

// Test for the help dialogue on invalid subcommand
func TestSubcommand(t *testing.T) {
	for _, tc := range cmdtests {
		t.Run(tc.in, func(t *testing.T) {
			output := captureOutput(func() {
				argsIn := append([]string{"run", "main.go"}, strings.Fields(tc.in)...)
				cmd := exec.Command("go", argsIn...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					t.Fatalf("Command failed: %v", err)
				}
			})

			if !strings.EqualFold(output, tc.out) {
				t.Errorf("Expected message not found. Got: %s", output)
			}
		})
	}
}
