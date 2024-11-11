package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBasicCommand(t *testing.T) {
	cmd := exec.Command("echo", "hello")
	err := cmd.Run()
	if err != nil {
		t.Errorf("Command failed with %s", err)
	}

	// Check exit code
	if cmd.ProcessState.ExitCode() != 0 {
		t.Errorf("Expected exit code 0, got %d", cmd.ProcessState.ExitCode())
	}
}

func TestLoginCommand(t *testing.T) {
	// Setup - might need to create a temporary config file

	t.Run("login with no username", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "login")
		err := cmd.Run()
		if err == nil {
			t.Error("Expected error for missing username, got nil")
		}
		if cmd.ProcessState.ExitCode() != 1 {
			t.Errorf("Expected exit code 1 for missing username, got %d", cmd.ProcessState.ExitCode())
		}
	})

	tmpConfigPath := filepath.Join(os.TempDir(), ".gatorconfig.json")
	initialConfig := []byte(`{"db_url": "test_url"}`)
	err := os.WriteFile(tmpConfigPath, initialConfig, 0644)
	if err != nil {
		t.Fatalf("Failed to create initial config: %v", err)
	}

	t.Run("login with username", func(t *testing.T) {
		defer os.Remove(tmpConfigPath)

		originalConfig := os.Getenv("RSS_GATOR_CONFIG")
		os.Setenv("RSS_GATOR_CONFIG", tmpConfigPath)
		defer os.Setenv("RSS_GATOR_CONFIG", originalConfig)

		cmd := exec.Command("go", "run", ".", "login", "testlogin")
		cmd.Env = append(os.Environ(), "RSS_GATOR_CONFIG="+tmpConfigPath)
		err := cmd.Run()
		if err != nil {
			t.Error("not expecting error with provided username")
		}
		if cmd.ProcessState.ExitCode() != 0 {
			t.Errorf("Expected exit code 0, got %d", cmd.ProcessState.ExitCode())
		}

		configData, err := os.ReadFile(tmpConfigPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		var config struct {
			CurrentUserName string `json:"CurrentUserName"`
		}
		if err := json.Unmarshal(configData, &config); err != nil {
			t.Fatalf("Failed to parse config file: %v", err)
		}

		// Verify the username was set correctly
		if config.CurrentUserName != "testlogin" {
			t.Errorf("Expected user to be 'testlogin', got '%s'", config.CurrentUserName)
		}
	})
}
