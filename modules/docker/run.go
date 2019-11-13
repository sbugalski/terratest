package docker

import (
	"testing"

	"github.com/sbugalski/terratest/modules/logger"
	"github.com/sbugalski/terratest/modules/shell"
	"github.com/stretchr/testify/require"
)

// RunOptions defines options that can be passed to the 'docker run' command.
type RunOptions struct {
	// Override the default COMMAND of the Docker image
	Command []string

	// If set to true, pass the --detach flag to 'docker run' to run the container in the background
	Detach bool

	// Override the default ENTRYPOINT of the Docker image
	Entrypoint string

	// Set environment variables
	EnvironmentVariables []string

	// If set to true, pass the --init flag to 'docker run' to run an init inside the container that forwards signals
	// and reaps processes
	Init bool

	// Assign a name to the container
	Name string

	// If set to true, pass the --privileged flag to 'docker run' to give extended privileges to the container
	Privileged bool

	// If set to true, pass the --rm flag to 'docker run' to automatically remove the container when it exits
	Remove bool

	// If set to true, pass the -tty flag to 'docker run' to allocate a pseudo-TTY
	Tty bool

	// Username or UID
	User string

	// Bind mount these volume(s) when running the container
	Volumes []string

	// Custom CLI options that will be passed as-is to the 'docker run' command. This is an "escape hatch" that allows
	// Terratest to not have to support every single command-line option offered by the 'docker run' command, and
	// solely focus on the most important ones.
	OtherOptions []string
}

// Run runs the 'docker run' command on the given image with the given options and return stdout/stderr. This method
// fails the test if there are any errors.
func Run(t *testing.T, image string, options *RunOptions) string {
	out, err := RunE(t, image, options)
	require.NoError(t, err)
	return out
}

// Run runs the 'docker run' command on the given image with the given options and return stdout/stderr, or any error.
func RunE(t *testing.T, image string, options *RunOptions) (string, error) {
	logger.Logf(t, "Running 'docker run' on image '%s'", image)

	args, err := formatDockerRunArgs(image, options)
	if err != nil {
		return "", err
	}

	cmd := shell.Command{
		Command: "docker",
		Args:    args,
	}

	return shell.RunCommandAndGetOutputE(t, cmd)
}

// formatDockerRunArgs formats the arguments for the 'docker run' command.
func formatDockerRunArgs(image string, options *RunOptions) ([]string, error) {
	args := []string{"run"}

	if options.Detach {
		args = append(args, "--detach")
	}

	if options.Entrypoint != "" {
		args = append(args, "--entrypoint", options.Entrypoint)
	}

	for _, envVar := range options.EnvironmentVariables {
		args = append(args, "--env", envVar)
	}

	if options.Init {
		args = append(args, "--init")
	}

	if options.Name != "" {
		args = append(args, "--name", options.Name)
	}

	if options.Privileged {
		args = append(args, "--privileged")
	}

	if options.Remove {
		args = append(args, "--rm")
	}

	if options.Tty {
		args = append(args, "--tty")
	}

	if options.User != "" {
		args = append(args, "--user", options.User)
	}

	for _, volume := range options.Volumes {
		args = append(args, "--volume", volume)
	}

	for _, opt := range options.OtherOptions {
		args = append(args, opt)
	}

	args = append(args, image)

	for _, arg := range options.Command {
		args = append(args, arg)
	}

	return args, nil
}
