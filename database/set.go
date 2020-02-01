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
		mediaItem := MediaItem{}

		// create or update media item
		err := tx.Where(MediaItem{
			Id:         item.ItemId,
			PvrName:    pvrName,
			WantedType: wantedType,
		}).Assign(MediaItem{
			PvrName:    pvrName,
			WantedType: wantedType,
			AirDateUtc: item.AirDateUtc,
		}).FirstOrCreate(&mediaItem).Error

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
