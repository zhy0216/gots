package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var (
	buildOnce sync.Once
	buildErr  error
	binPath   string
)

func ensureBinary(t *testing.T) string {
	buildOnce.Do(func() {
		binPath = filepath.Join(os.TempDir(), "gots_test_bin")
		cmd := exec.Command("go", "build", "-o", binPath, "./cmd/gots")
		if output, err := cmd.CombinedOutput(); err != nil {
			buildErr = err
			t.Logf("Build output: %s", output)
		}
	})
	if buildErr != nil {
		t.Fatalf("Failed to build gots binary: %v", buildErr)
	}
	return binPath
}

// TestIntegration runs all .gts files in the test directory
func TestIntegration(t *testing.T) {
	bin := ensureBinary(t)

	files, err := filepath.Glob("test/*.gts")
	if err != nil {
		t.Fatalf("Failed to glob test files: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No .gts test files found in test/")
	}

	for _, file := range files {
		file := file
		testName := strings.TrimSuffix(filepath.Base(file), ".gts")

		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			cmd := exec.Command(bin, "run", file)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Failed to run %s: %v\nOutput:\n%s", file, err, output)
			} else {
				t.Logf("Output:\n%s", output)
			}
		})
	}
}

// TestIntegrationEmitGo tests that files compile to Go successfully
func TestIntegrationEmitGo(t *testing.T) {
	bin := ensureBinary(t)

	files, err := filepath.Glob("test/*.gts")
	if err != nil {
		t.Fatalf("Failed to glob test files: %v", err)
	}

	for _, file := range files {
		file := file
		testName := strings.TrimSuffix(filepath.Base(file), ".gts")

		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			cmd := exec.Command(bin, "emit-go", file)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Failed to emit Go for %s: %v\nOutput:\n%s", file, err, output)
			}
		})
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	// Cleanup
	if binPath != "" {
		os.Remove(binPath)
	}
	os.Exit(code)
}
