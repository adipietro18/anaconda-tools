package main

import (
	"fmt"
	"os"
	"testing"
)

func resetSshEnv() {
	os.Unsetenv("TEST_SSH_HOST")
	os.Unsetenv("TEST_SSH_PATH")
	os.Unsetenv("TEST_SSH_USER")
}

func TestBuildRsyncCommand(t *testing.T) {
	// expectedString could be of type "func(string) bool" instead. then the
	// function could simply check if what is need is present or not?:
	// func(result string) bool { return strings.Contains(result, " --dry-run") }
	cases := []struct {
		parameters     CopyFromParameters
		source         SshLocation
		expectedString string
		errorMessage   string
	}{
		{
			CopyFromParameters{DryRun: true, Platform: "s390x"},
			SshLocation{Host: "localhost", Path: "./", User: "user"},
			"rsync -amv --dry-run --include='linux-s390x/' --include='*.conda' --include='*.tar.bz2' --exclude='*' user@localhost:./ ./",
			"Dry run flag was not present in the rsync command.",
		},
		{
			CopyFromParameters{DryRun: false, Platform: ""},
			SshLocation{Host: "", Path: "", User: ""},
			"rsync -amv --include='linux-/' --include='*.conda' --include='*.tar.bz2' --exclude='*' @: ./",
			"Dry run flag was unexpectedly present in the rsync command.",
		},
		{
			CopyFromParameters{DryRun: false, Platform: "s390x"},
			SshLocation{Host: "", Path: "", User: ""},
			"rsync -amv --include='linux-s390x/' --include='*.conda' --include='*.tar.bz2' --exclude='*' @: ./",
			"Linux platform was not included properly.",
		},
		{
			CopyFromParameters{DryRun: false, Platform: ""},
			SshLocation{Host: "host", Path: "/path/", User: "user"},
			"rsync -amv --include='linux-/' --include='*.conda' --include='*.tar.bz2' --exclude='*' user@host:/path/ ./",
			"SSH location was not included properly.",
		},
	}
	for _, testCase := range cases {
		command := BuildRsyncCommand(testCase.parameters, testCase.source)
		if command != testCase.expectedString {
			t.Error(fmt.Sprintf("%s\n\texpected: %s\n\treceived: %s", testCase.errorMessage, testCase.expectedString, command))
		}
	}
}

func TestValidateParameters(t *testing.T) {
	// todo: fuzzing!
	cases := []struct {
		parameters    CopyFromParameters
		errorExpected bool
		errorMessage  string
	}{
		{
			CopyFromParameters{DryRun: false, Platform: "s390x"},
			false,
			"Expected s390x to be valid.",
		},
		{
			CopyFromParameters{DryRun: false, Platform: "aarch64"},
			false,
			"Expected aarch64 to be valid.",
		},
		{
			CopyFromParameters{DryRun: false, Platform: "S390x"},
			true,
			"Expected S390x to be invalid (uppercase S).",
		},
		{
			CopyFromParameters{DryRun: false, Platform: "aarch65"},
			true,
			"Expected aarch65 to be invalid (unknown architecture).",
		},
		{
			CopyFromParameters{DryRun: false, Platform: "osx64"},
			true,
			"Expected osx64 to be invalid (invalid platform).",
		},
	}

	for _, testCase := range cases {
		if err := ValidateParameters(testCase.parameters); (err != nil) != testCase.errorExpected {
			t.Error(testCase.errorMessage)
		}
	}
}

func TestGetSshLocation(t *testing.T) {
	defer resetSshEnv()
	cases := []struct {
		platform         string
		setupEnv         func()
		expectedLocation SshLocation
		errorExpected    bool
		errorMessage     string
	}{
		{
			"TEST",
			func() {
				os.Setenv("TEST_SSH_USER", "user")
				os.Setenv("TEST_SSH_HOST", "host")
				os.Setenv("TEST_SSH_PATH", "/path")
			},
			SshLocation{User: "user", Host: "host", Path: "/path"},
			false,
			"Expected successful setup of SshLocation from environment variables.",
		},
		{
			"test",
			func() {
				os.Setenv("TEST_SSH_USER", "user")
				os.Setenv("TEST_SSH_HOST", "host")
				os.Setenv("TEST_SSH_PATH", "/path")
			},
			SshLocation{User: "user", Host: "host", Path: "/path"},
			false,
			"Expected the platform to be case insensitive.",
		},
		{
			"test",
			func() {
				os.Setenv("TEST_SSH_HOST", "host")
				os.Setenv("TEST_SSH_PATH", "/path")
			},
			SshLocation{},
			true,
			"Expected error due to missing TEST_SSH_USER.",
		},
		{
			"test",
			func() {
				os.Setenv("TEST_SSH_USER", "user")
				os.Setenv("TEST_SSH_PATH", "/path")
			},
			SshLocation{},
			true,
			"Expected error due to missing TEST_SSH_HOST.",
		},
		{
			"test",
			func() {
				os.Setenv("TEST_SSH_USER", "user")
				os.Setenv("TEST_SSH_HOST", "host")
			},
			SshLocation{},
			true,
			"Expected error due to missing TEST_SSH_PATH.",
		},
	}
	for _, testCase := range cases {
		resetSshEnv()
		testCase.setupEnv()
		location, err := getSshLocation(testCase.platform)
		if (err != nil) != testCase.errorExpected {
			t.Error(testCase.errorMessage)
		} else if location != testCase.expectedLocation {
			t.Error(testCase.errorMessage)
		}

	}
}
