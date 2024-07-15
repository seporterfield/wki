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
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	output := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output <- buf.String()
	}()

	f()
	w.Close()
	os.Stdout = stdout
	return <-output
}

// Test for the help dialogue on invalid subcommand
func TestInvalidSubcommandShowsHelp(t *testing.T) {
	// Simulate running the command `wki garbage`
	output := captureOutput(func() {
		cmd := exec.Command("go", "run", "main.go", "garbage")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}
	})

	// Check that the output contains the expected help message
	expectedHelpMessage := "help blah blah"
	if !strings.Contains(output, expectedHelpMessage) {
		t.Errorf("Expected help message not found. Got: %s", output)
	}
}
