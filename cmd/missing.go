package cmd

import (
	"fmt"
	"github.com/l3uddz/wantarr/config"
	pvrObj "github.com/l3uddz/wantarr/pvr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tommysolsen/capitalise"
)

var (
	maxQueueSize int
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

		// get missing
		log.Infof("Searching missing media in %s named: %q", capitalise.First(pvrConfig.Type), pvrName)

		if qSize, err := pvr.GetQueueSize(); err != nil {
			log.WithError(err).Error("Failed retrieving queued pvr items...")
		} else {
			log.WithField("size", qSize).Info("Refreshed download queue")
		}

		if _, err := pvr.GetWantedMissing(); err != nil {
			log.WithError(err).Error("Failed retrieving wanted missing pvr items...")
		}
	},
}

func init() {
	rootCmd.AddCommand(missingCmd)

	missingCmd.Flags().IntVarP(&maxQueueSize, "queue-size", "q", 5, "Exit when queue size reached.")
}

func parseValidateInputs(args []string) error {
	var ok bool = false
	var err error = nil

	// validate pvr exists in config
	pvrName = args[0]
	pvrConfig, ok = config.Config.Pvr[pvrName]
	if !ok {
		return fmt.Errorf("no pvr configuration found for: %q", pvrName)
	}

	// init pvrObj
	pvr, err = pvrObj.Get(pvrName, pvrConfig.Type, pvrConfig)
	if err != nil {
		return errors.WithMessage(err, "failed loading pvr object")
	}

	return nil
}
