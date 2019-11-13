package docker

import (
	"testing"

	"github.com/sbugalski/terratest/modules/logger"
	"github.com/sbugalski/terratest/modules/shell"
	"github.com/stretchr/testify/require"
)

// BuildOptions defines options that can be passed to the 'docker build' command.
type BuildOptions struct {
	// Tags for the Docker image
	Tags []string

	// Build args to pass the 'docker build' command
	BuildArgs []string

	// Custom CLI options that will be passed as-is to the 'docker build' command. This is an "escape hatch" that allows
	// Terratest to not have to support every single command-line option offered by the 'docker build' command, and
	// solely focus on the most important ones.
	OtherOptions []string
}

// Build runs the 'docker build' command at the given path with the given options and fails the test if there are any
// errors.
func Build(t *testing.T, path string, options *BuildOptions) {
	require.NoError(t, BuildE(t, path, options))
}

// BuildE runs the 'docker build' command at the given path with the given options and returns any errors.
func BuildE(t *testing.T, path string, options *BuildOptions) error {
	logger.Logf(t, "Running 'docker build' in %s", path)

	args, err := formatDockerBuildArgs(path, options)
	if err != nil {
		return err
	}

	cmd := shell.Command{
		Command: "docker",
		Args:    args,
	}

	_, buildErr := shell.RunCommandAndGetOutputE(t, cmd)
	return buildErr
}

// formatDockerBuildArgs formats the arguments for the 'docker build' command.
func formatDockerBuildArgs(path string, options *BuildOptions) ([]string, error) {
	args := []string{"build"}

	for _, tag := range options.Tags {
		args = append(args, "--tag", tag)
	}

	for _, arg := range options.BuildArgs {
		args = append(args, "--build-arg", arg)
	}

	for _, opt := range options.OtherOptions {
		args = append(args, opt)
	}

	args = append(args, path)

	return args, nil
}
