package cmd

import (
	"fmt"
	"github.com/l3uddz/wantarr/config"
	"github.com/spf13/cobra"
	"github.com/tommysolsen/capitalise"
)

var (
	maxQueueSize int
	pvrName      string
	pvrConfig    *config.Pvr
)

var missingCmd = &cobra.Command{
	Use:   "missing [PVR]",
	Short: "Search for missing media files",
	Long:  `This command can be used to search for missing media files from the respective arr wanted list.`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// validate inputs
		if err := parseValidateInputs(args); err != nil {
			log.WithError(err).Fatal("Failed validating inputs")
		}

		// do search
		log.Infof("Searching missing media in %s named: %q", capitalise.First(pvrConfig.Type), pvrName)
	},
}

func init() {
	rootCmd.AddCommand(missingCmd)

	missingCmd.Flags().IntVarP(&maxQueueSize, "queue-size", "q", 5, "Exit when queue size reached.")
}

func parseValidateInputs(args []string) error {
	ok := false

	// validate pvr exists in config
	pvrName = args[0]
	pvrConfig, ok = config.Config.Pvr[pvrName]
	if !ok {
		return fmt.Errorf("no pvr configuration found for: %q", pvrName)
	}

	return nil
}
