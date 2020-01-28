package cmd

import (
	"fmt"
	"github.com/l3uddz/wantarr/config"
	"github.com/l3uddz/wantarr/database"
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

		// load database
		db, err := database.New(pvrName, flagConfigFolder)
		if err != nil {
			log.WithError(err).Fatal("Failed initializing database file...")
		}

		// close database on finish
		defer func() {
			if err := db.Close(); err != nil {
				log.WithError(err).Error("Failed closing database gracefully...")
			}
		}()

		// retrieve missing records from pvr and store in database
		if flagRefreshCache || !db.FromDisk() {
			log.Infof("Retrieving missing media in %s named: %q", capitalise.First(pvrConfig.Type), pvrName)

			missingRecords, err := pvr.GetWantedMissing()
			if err != nil {
				log.WithError(err).Error("Failed retrieving wanted missing pvr items...")
			}

			// store missing media in database
			log.Info("Storing records in database...")

			for itemId, record := range missingRecords {
				if err := db.Set(itemId, &record, true); err != nil {
					log.WithError(err).Errorf("Failed storing mediaItem %q in database", record)
				}
			}

			log.Info("Stored records to database")

			// remove media no longer missing
			if db.FromDisk() {
				log.Info("Removing records from database that are no longer missing...")

				removedRecords := 0
				for itemId, _ := range *db.GetVault() {
					if _, itemStillMissing := missingRecords[itemId]; !itemStillMissing {
						// this item is no longer missing
						db.Delete(itemId)
						removedRecords += 1
					}
				}

				log.Infof("Removed %d records from database that are no longer missing!", removedRecords)
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
