package database

import (
	"github.com/l3uddz/wantarr/pvr"
	"github.com/pkg/errors"
	"strconv"
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

		// create item if not exists
		err := tx.Where(MediaItem{
			Id:         strconv.Itoa(item.ItemId),
			PvrName:    pvrName,
			WantedType: wantedType,
		}).Assign(mediaItem).FirstOrCreate(&mediaItem).Error

		if err != nil {
			log.WithError(err).Errorf("Failed inserting media item: %v", item.ItemId)
		}

		// update item
		if !item.LastSearch.IsZero() {
			mediaItem.AirDateUtc = item.AirDateUtc
			mediaItem.LastSearchDateUtc = &item.LastSearch

			if err := tx.Save(&mediaItem).Error; err != nil {
				log.WithError(err).Errorf("Failed updating media item: %v", item.ItemId)
			}
		}
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		log.WithError(err).Error("Failed commit bulk insert/update of media items...")
		return errors.Wrap(err, "failed committing bulk transaction")
	}

	return nil
}
