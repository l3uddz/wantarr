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
		if err := database.Init(flagDatabaseFile); err != nil {
			log.WithError(err).Fatal("Failed opening database file")
		}
		defer database.Close()

		// retrieve missing records from pvr and stash in database
		if flagRefreshCache {
			log.Infof("Retrieving missing media from %s: %q", capitalise.First(pvrConfig.Type), pvrName)

			missingRecords, err := pvr.GetWantedMissing()
			if err != nil {
				log.WithError(err).Fatal("Failed retrieving wanted missing pvr items...")
			}

			// stash missing media in database
			log.Debug("Stashing media items in database...")

			if err := database.SetMediaItems(strings.ToLower(pvrName), "missing", missingRecords); err != nil {
				log.WithError(err).Fatal("Failed stashing media items in database")
			}

			log.Info("Stashed media items")

			// remove media no longer missing
			log.Debug("Removing media items from database that are no longer missing...")

			removedItems, err := database.DeleteMissingItems(strings.ToLower(pvrName), "missing", missingRecords)
			if err != nil {
				log.WithError(err).Fatal("Failed removing media items from database that are no longer missing...")
			}

			log.WithField("removed_items", removedItems).
				Info("Removed media items from database that are no longer missing")

			//if err := database.SetMediaItems(pvrName, "missing", missingRecords); err != nil {
			//	log.WithError(err).Fatal("Failed stashing media items in database...")
			//}

			//newItems, err := db.SetMediaItems(missingRecords)
			//if err != nil {
			//	log.WithError(err).Errorf("Failed stashing media items in database...")
			//	return
			//}
			//
			//log.WithField("new_media_items", newItems).Info("Stashed media items in database")
			//
			//// remove media no longer missing
			//if db.FromDisk() {
			//	log.Debug("Removing media items from database that are no longer missing...")
			//
			//	removedItems, err := db.RemoveMissingMediaItems(missingRecords)
			//	if err != nil {
			//		log.WithError(err).Error("Failed removing media items from database that are no longer missing")
			//	} else {
			//		log.WithField("media_items", removedItems).
			//			Info("Removed media items from database that are no longer missing")
			//	}
			//}
		}
		//
		//// start queue monitor
		//
		//// start searching
		//pos := 0
		//var searchItemIds []int
		//searchItems := make(map[int]pvrObj.MediaItem, 0)
		//
		//for itemId, item := range *db.GetVault() {
		//	if !item.LastSearch.IsZero() {
		//		continue
		//	}
		//
		//	if pos > 9 {
		//		break
		//	} else {
		//		pos++
		//	}
		//
		//	itm := item
		//	searchItems[itemId] = itm
		//	searchItemIds = append(searchItemIds, itemId)
		//}
		//
		//log.Info("Searching")
		//log.Info(searchItemIds)
		//
		//ok, err := pvr.SearchMediaItems(searchItemIds)
		//if err != nil {
		//	log.WithError(err).Fatal("Failed searching for items")
		//} else if !ok {
		//	log.Error("Failed searching for items!")
		//} else {
		//	log.Info("Searched for items!")
		//
		//	for itemId, item := range searchItems {
		//		item.LastSearch = time.Now().UTC()
		//		_ = db.Set(itemId, &item, false)
		//	}
		//}

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
