package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func buildRsyncCommand(platform string, dryRun bool) (*exec.Cmd, error) {
	upperPlatform := strings.ToUpper(platform)

	sshUser := os.Getenv(upperPlatform + "_SSH_USER")
	if sshUser == "" {
		return nil, errors.New(upperPlatform + "_SSH_USER environment variable must be defined.")
	}

	sshHost := os.Getenv(upperPlatform + "_SSH_HOST")
	if sshHost == "" {
		return nil, errors.New(upperPlatform + "_SSH_HOST environment variable must be defined.")
	}

	sshPath := os.Getenv(upperPlatform + "_SSH_PATH")
	if sshPath == "" {
		return nil, errors.New(upperPlatform + "_SSH_PATH environment variable must be defined.")
	}

	command := []string{"-amv"}
	if dryRun {
		command = append(command, "--dry-run")
	}
	command = append(command, fmt.Sprintf("--include='linux-%s'", platform))
	command = append(command, "--include='*.conda'")
	command = append(command, "--include='*.tar.bz2'")
	command = append(command, "--exclude='*'")
	command = append(command, fmt.Sprintf("%s@%s:%s", sshUser, sshHost, sshPath))
	command = append(command, "./")

	return exec.Command("rsync", command...), nil
}

func parsePlatform(platform string) (string, error) {
	switch platform {
	case "s390x", "aarch64":
		return platform, nil
	default:
		return "", errors.New("the provided platform must be either s390x or aarch64")
	}
}

func usage() {
	const usage = `Usage: copy-from [options] <target>

Options:
    -n, --dry-run           show what would have been transferred

Copies packages from the identified TARGET to the local directory. The following
environment variables must be set for the respective TARGET. The variable names
must be in uppercase.

These are:
    ${TARGET}_SSH_USER      the user to use on the remote machine
    ${TARGET}_SSH_HOST      the host / IP of the remote machine
    ${TARGET}_SSH_PATH      the artifact path to copy from the remote machine
                            typically be the conda-bld / croot directory.

Valid TARGETs:
    aarch64                 linux-aarch64
    s390x                   linux-s390x

Example:
    copy-from -n s390x
`
	fmt.Fprintf(os.Stderr, "%s\n", usage)
}

func main() {
	flag.Usage = usage

	var dryRun bool
	flag.BoolVar(&dryRun, "n", false, "")
	flag.BoolVar(&dryRun, "dry-run", false, "")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	platform, err := parsePlatform(flag.Arg(0))
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	rsyncCommand, err := buildRsyncCommand(platform, dryRun)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var outBytes, errBytes bytes.Buffer
	rsyncCommand.Stdout = &outBytes
	rsyncCommand.Stderr = &errBytes

	fmt.Println("Running: " + rsyncCommand.String())

	err = rsyncCommand.Run()

	if err != nil {
		fmt.Println(errBytes.String())
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(outBytes.String())

	os.Exit(0)
}
