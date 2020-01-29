package cmd

import (
	"fmt"
	"github.com/l3uddz/wantarr/config"
	"github.com/l3uddz/wantarr/database"
	pvrObj "github.com/l3uddz/wantarr/pvr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tommysolsen/capitalise"
	"strings"
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

		// init pvr object
		if err := pvr.Init(); err != nil {
			log.WithError(err).Fatalf("Failed initializing pvr object for: %s", pvrName)
		}

		// load database
		db, err := database.New(strings.ToLower(pvrName), "missing", flagConfigFolder)
		if err != nil {
			log.WithError(err).Fatal("Failed initializing database file...")
		}

		// close database on finish
		defer func() {
			if err := db.Close(); err != nil {
				log.WithError(err).Error("Failed closing database gracefully...")
			}
		}()

		// retrieve missing records from pvr and stash in database
		if flagRefreshCache || !db.FromDisk() {
			log.Infof("Retrieving missing media from %s: %q", capitalise.First(pvrConfig.Type), pvrName)

			missingRecords, err := pvr.GetWantedMissing()
			if err != nil {
				log.WithError(err).Fatal("Failed retrieving wanted missing pvr items...")
			}

			// stash missing media in database
			log.Debug("Stashing media items in database...")

			for itemId, record := range missingRecords {
				if err := db.Set(itemId, &record, true); err != nil {
					log.WithError(err).Errorf("Failed stashing mediaItem %q in database", record)
				}
			}

			log.Info("Stashed media items in database")

			// remove media no longer missing
			if db.FromDisk() {
				log.Debug("Removing media items from database that are no longer missing...")

				removedItems, err := db.RemoveMissingMediaItems(missingRecords)
				if err != nil {
					log.WithError(err).Error("Failed removing media items from database that are no longer missing")
				} else {
					log.WithField("media_items", removedItems).
						Info("Removed media items from database that are no longer missing")
				}
			}
		}

		// start queue monitor

		// start searching
	},
}

func init() {
	rootCmd.AddCommand(missingCmd)

	missingCmd.Flags().IntVarP(&maxQueueSize, "queue-size", "q", 5, "Exit when queue size reached.")
	missingCmd.Flags().BoolVarP(&flagRefreshCache, "refresh-cache", "r", false, "Refresh the locally stored cache.")
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
