package database

import (
	"github.com/l3uddz/wantarr/pvr"
	"github.com/pkg/errors"
)

func SetMediaItems(pvrName string, wantedType string, mediaItems []pvr.MediaItem) error {
	// begin transaction
	tx := db.Begin()

	// bulk insert/update items
	for _, item := range mediaItems {
		// set item to insert/update
		mediaItem := MediaItem{
			PvrName:    pvrName,
			WantedType: wantedType,
			AirDateUtc: item.AirDateUtc,
		}

		// update items last search if available
		if !item.LastSearch.IsZero() {
			mediaItem.LastSearchDateUtc = &item.LastSearch
		}

		// create or update media item
		err := tx.Where(MediaItem{
			Id:         item.ItemId,
			PvrName:    pvrName,
			WantedType: wantedType,
		}).Assign(mediaItem).FirstOrCreate(&mediaItem).Error

		if err != nil {
			log.WithError(err).Errorf("Failed inserting media item: %v", item.ItemId)
		}
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		log.WithError(err).Error("Failed commit bulk insert/update of media items...")
		return errors.Wrap(err, "failed committing bulk transaction")
	}

	return nil
}
