package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/htfs"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var holotreeHashCmd = &cobra.Command{
	Use:   "hash <conda.yaml+>",
	Short: "Calculates a blueprint hash for managed holotree virtual environment from conda.yaml files.",
	Long:  "Calculates a blueprint hash for managed holotree virtual environment from conda.yaml files.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Conda YAML hash calculation lasted").Report()
		}
		_, holotreeBlueprint, err := htfs.ComposeFinalBlueprint(args, "")
		pretty.Guard(err == nil, 1, "Blueprint calculation failed: %v", err)
		hash := htfs.BlueprintHash(holotreeBlueprint)
		common.Log("Blueprint hash for %v is %v.", args, hash)
		if common.Silent {
			common.Stdout("%s\n", hash)
		}
	},
}

func init() {
	holotreeCmd.AddCommand(holotreeHashCmd)
}
