package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
	"github.com/spinframework/otel-plugin/internal/stack"
)

var (
	removeVolumes = false
	cleanUpCmd    = &cobra.Command{
		Use:   "cleanup",
		Short: "Clean up OpenTelemetry dependencies",
		Long:  "Tears down the OpenTelemetry dependency containers started by \"spin otel setup\" (the inverse of setup). Works regardless of which stack was set up.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cleanUp(); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	cleanUpCmd.Flags().BoolVarP(&removeVolumes, "remove-volumes", "v", false, "Also remove named volumes declared in the compose files")
}

func cleanUp() error {
	runtime, err := detectContainerRuntime()
	if err != nil {
		return err
	}

	fmt.Println("Stopping and removing Spin OpenTelemetry containers...")

	for _, s := range stack.AllStacks() {
		composeFileName := s.GetComposeFileName()
		composeFilePath := path.Join(otelConfigPath, composeFileName)
		if _, err := os.Stat(composeFilePath); os.IsNotExist(err) {
			// A stack's compose file may not be present; nothing to tear down.
			continue
		}

		downArgs := []string{"compose", "-f", composeFilePath, "down", "--remove-orphans"}
		if removeVolumes {
			downArgs = append(downArgs, "--volumes")
		}

		cmd := exec.Command(runtime, downArgs...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			return err
		}
	}

	fmt.Println("All Spin OpenTelemetry resources have been cleaned up.")
	return nil
}
