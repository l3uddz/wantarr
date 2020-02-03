package cmd

import (
	"github.com/l3uddz/wantarr/database"
	pvrObj "github.com/l3uddz/wantarr/pvr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tommysolsen/capitalise"
	"time"
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
		if maxQueueSize > 0 {
			go func() {
				log.Info("Started queue monitor")
				for {
					// retrieve queue size
					qs, err := pvr.GetQueueSize()
					if err != nil {
						log.WithError(err).Error("Failed retrieving queue size, aborting...")
						continueRunning.Store(false)
						break
					}

					// check queue size
					if qs >= maxQueueSize {
						log.Warnf("Queue size has been reached, aborting....")
						continueRunning.Store(false)
						break
					}

					// sleep before check
					time.Sleep(10 * time.Second)
				}
				log.Info("Finished queue monitor")
			}()
		}

		// get media items from database
		mediaItems, err := database.GetMediaItems(lowerPvrName, "missing", true)
		if err != nil {
			log.WithError(err).Fatal("Failed retrieving media items from database...")
		}
		log.WithField("media_items", len(mediaItems)).Debug("Retrieved media items from database")

		// start searching
		var searchItems []pvrObj.MediaItem
		searchedItemsCount := 0

		for _, item := range mediaItems {
			// abort if required (queue monitor will set this)
			if !continueRunning.Load() {
				break
			}

			// dont search this item if we already searched it within N days
			if item.LastSearchDateUtc != nil && !item.LastSearchDateUtc.IsZero() {
				retryAfterDate := item.LastSearchDateUtc.Add((24 * pvrConfig.RetryDaysAge.Missing) * time.Hour)
				if time.Now().UTC().Before(retryAfterDate) {
					log.WithField("retry_min_date", retryAfterDate).
						Tracef("Skipping media item %v until allowed retry date", item.Id)
					continue
				}
			}

			// add item to batch
			searchItems = append(searchItems, pvrObj.MediaItem{
				ItemId:     item.Id,
				AirDateUtc: item.AirDateUtc,
			})

			// not enough items batched yet
			batchedItemsCount := len(searchItems)
			if batchedItemsCount < searchBatchSize {
				continue
			}

			// do search
			log.WithFields(logrus.Fields{
				"search_items":   batchedItemsCount,
			}).Info("Searching...")

			searchedItemsCount += batchedItemsCount

			if _, err := searchForItems(searchItems, "missing"); err != nil {
				log.WithError(err).Error("Failed searching for items...")
			} else {
				log.WithFields(logrus.Fields{
					"searched_items": searchedItemsCount,
				}).Info("Search complete")
			}

			// reset batch
			searchItems = []pvrObj.MediaItem{}

			// max search items reached?
			if maxSearchItems > 0 && searchedItemsCount >= maxSearchItems {
				log.WithField("searched_items", searchedItemsCount).
					Info("Max search items reached, aborting...")
				break
			}

			// sleep before next batch
			time.Sleep(5 * time.Second)
		}

		// search for any leftover items from batching
		if continueRunning.Load() && len(searchItems) > 0 {
			// search items
			log.WithFields(logrus.Fields{
				"search_items":   len(searchItems),
			}).Info("Searching...")

			searchedItemsCount += len(searchItems)

			if _, err := searchForItems(searchItems, "missing"); err != nil {
				log.WithError(err).Error("Failed searching for items...")
			} else {
				log.WithFields(logrus.Fields{
					"searched_items": searchedItemsCount,
				}).Info("Search complete")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(missingCmd)

	missingCmd.Flags().IntVarP(&maxQueueSize, "queue-size", "q", 0, "Exit when queue size reached.")
	missingCmd.Flags().IntVarP(&maxSearchItems, "max-search", "m", 0, "Exit when this many items have been searched.")
	missingCmd.Flags().IntVarP(&searchBatchSize, "search-size", "s", 10, "How many items to search at once.")
	missingCmd.Flags().BoolVarP(&flagRefreshCache, "refresh-cache", "r", false, "Refresh the locally stored cache.")
}
