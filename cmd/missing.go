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
	"time"
)

var (
	maxQueueSize    int
	searchBatchSize int
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
		if err := database.Init(flagDatabaseFile); err != nil {
			log.WithError(err).Fatal("Failed opening database file")
		}
		defer database.Close()

		// retrieve missing records from pvr and stash in database
		existingItemsCount := database.GetItemsCount(lowerPvrName, "missing")
		if flagRefreshCache || existingItemsCount < 1 {
			log.Infof("Retrieving missing media from %s: %q", capitalise.First(pvrConfig.Type), pvrName)

			missingRecords, err := pvr.GetWantedMissing()
			if err != nil {
				log.WithError(err).Fatal("Failed retrieving wanted missing pvr items...")
			}

			// stash missing media in database
			log.Debug("Stashing media items in database...")

			if err := database.SetMediaItems(lowerPvrName, "missing", missingRecords); err != nil {
				log.WithError(err).Fatal("Failed stashing media items in database")
			}

			log.Info("Stashed media items")

			// remove media no longer missing
			if existingItemsCount >= 1 {
				log.Debug("Removing media items from database that are no longer missing...")

				removedItems, err := database.DeleteMissingItems(lowerPvrName, "missing", missingRecords)
				if err != nil {
					log.WithError(err).Fatal("Failed removing media items from database that are no longer missing...")
				}

				log.WithField("removed_items", removedItems).
					Info("Removed media items from database that are no longer missing")
			}
		}

		// start queue monitor

		// get media items from database
		mediaItems, err := database.GetMediaItems(lowerPvrName, "missing")
		if err != nil {
			log.WithError(err).Fatal("Failed retrieving media items from database...")
		}
		log.WithField("media_items", len(mediaItems)).Debug("Retrieved media items from database")

		// start searching
		var searchItems []pvrObj.MediaItem
		for _, item := range mediaItems {
			// abort if required (queue monitor will set this)
			if !continueRunning.Load() {
				break
			}

			// dont search this item if we already searched it within N days
			if item.LastSearchDateUtc != nil && !item.LastSearchDateUtc.IsZero() {
				retryAfterDate := item.LastSearchDateUtc.Add((24 * pvrConfig.RetryDaysAge) * time.Hour)
				if time.Now().UTC().Before(retryAfterDate) {
					log.WithField("retry_min_date", retryAfterDate).Tracef("Skipping media item %v until allowed retry date", item.Id)
					continue
				}
			}

			// add item to batch
			searchItems = append(searchItems, pvrObj.MediaItem{
				ItemId:     item.Id,
				AirDateUtc: item.AirDateUtc,
			})

			// not enough items batched yet
			if len(searchItems) < searchBatchSize {
				continue
			}

			// search items
			if _, err := searchForItems(searchItems); err != nil {
				log.WithError(err).Error("Failed searching for items...")
			}

			// reset batch
			searchItems = []pvrObj.MediaItem{}
		}

		// search for any leftover items from batching
		if continueRunning.Load() && len(searchItems) > 0 {
			// search items
			if _, err := searchForItems(searchItems); err != nil {
				log.WithError(err).Error("Failed searching for items...")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(missingCmd)

	missingCmd.Flags().IntVarP(&maxQueueSize, "queue-size", "q", 10, "Exit when queue size reached.")
	missingCmd.Flags().IntVarP(&searchBatchSize, "search-batch-size", "s", 10, "How many items to search at once.")
	missingCmd.Flags().BoolVarP(&flagRefreshCache, "refresh-cache", "r", false, "Refresh the locally stored cache.")
}

func parseValidateInputs(args []string) error {
	var ok bool = false
	var err error = nil

	// validate pvr exists in config
	pvrName = args[0]
	lowerPvrName = strings.ToLower(pvrName)
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
