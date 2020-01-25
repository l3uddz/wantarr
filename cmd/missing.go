package cmd

import "github.com/spf13/cobra"

var (
	maxQueueSize int
)

var missingCmd = &cobra.Command{
	Use:   "missing [PVR]",
	Short: "Search for missing media files",
	Long:  `This command can be used to search for missing media files from the respective arr wanted list.`,

	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pvrName := args[0]

		log.Infof("Searching for missing media in: %q with max queue size of %d", pvrName, maxQueueSize)
	},
}

func init() {
	rootCmd.AddCommand(missingCmd)

	missingCmd.Flags().IntVarP(&maxQueueSize, "queue-size", "q", 0, "Exit when queue size reached.")
}
