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

type CopyFromParameters struct {
	DryRun   bool
	Platform string
}

type SshLocation struct {
	User string
	Host string
	Path string
}

func RunCopyFrom(parameters CopyFromParameters, source SshLocation) error {
	var commandBuilder strings.Builder
	commandBuilder.WriteString("rsync -amv")
	if parameters.DryRun {
		commandBuilder.WriteString(" --dry-run")
	}
	commandBuilder.WriteString(fmt.Sprintf(" --include='linux-%s/'", parameters.Platform))
	commandBuilder.WriteString(" --include='*.conda'")
	commandBuilder.WriteString(" --include='*.tar.bz2'")
	commandBuilder.WriteString(" --exclude='*'")
	commandBuilder.WriteString(fmt.Sprintf(" %s@%s:%s", source.User, source.Host, source.Path))
	commandBuilder.WriteString(" ./")

	command := exec.Command("bash", "-c", commandBuilder.String())

	var outBytes, errBytes bytes.Buffer
	command.Stdout = &outBytes
	command.Stderr = &errBytes

	fmt.Println("Running: " + command.String())

	err := command.Run()

	if err != nil {
		fmt.Println(errBytes.String())
		return err
	}

	fmt.Println(outBytes.String())

	return nil
}

func ValidateParameters(parameters CopyFromParameters) error {
	err := validatePlatform(parameters.Platform)
	if err != nil {
		return err
	}
	return nil
}

func getSshLocation(platform string) (SshLocation, error) {
	upperPlatform := strings.ToUpper(platform)
	sshUser := os.Getenv(upperPlatform + "_SSH_USER")
	if sshUser == "" {
		return SshLocation{}, errors.New(upperPlatform + "_SSH_USER environment variable must be defined.")
	}

	sshHost := os.Getenv(upperPlatform + "_SSH_HOST")
	if sshHost == "" {
		return SshLocation{}, errors.New(upperPlatform + "_SSH_HOST environment variable must be defined.")
	}

	sshPath := os.Getenv(upperPlatform + "_SSH_PATH")
	if sshPath == "" {
		return SshLocation{}, errors.New(upperPlatform + "_SSH_PATH environment variable must be defined.")
	}

	return SshLocation{User: sshUser, Host: sshHost, Path: sshPath}, nil
}

func normalizePlatform(platform string) string {
	return strings.ToLower(platform)
}

func validatePlatform(platform string) error {
	switch platform {
	case "s390x", "aarch64":
		return nil
	default:
		return errors.New("the provided platform is not one of the expected values")
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

	parameters := CopyFromParameters{DryRun: dryRun, Platform: normalizePlatform(flag.Arg(0))}

	err := ValidateParameters(parameters)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	source, err := getSshLocation(parameters.Platform)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = RunCopyFrom(parameters, source); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
