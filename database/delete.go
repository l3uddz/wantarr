package database

import (
	"github.com/l3uddz/wantarr/pvr"
	"github.com/pkg/errors"
)

func DeleteMissingItems(pvrName string, wantedType string, newMediaItems []pvr.MediaItem) (int, error) {
	// build slice of new item ids
	newItemIds := make(map[int]*string)

	for _, item := range newMediaItems {
		newItemIds[item.ItemId] = nil
	}

	// begin transaction
	tx := db.Begin()

	// retrieve existing items
	var dbItems []MediaItem

	if err := tx.Where("pvr_name = ? AND wanted_type = ?", pvrName, wantedType).Find(&dbItems).Error; err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "failed finding existing media items")
	}

	// iterate database items finding those that no longer exist
	removedItems := 0

	for _, item := range dbItems {
		if _, ok := newItemIds[item.Id]; !ok {
			item := item

			// item no longer exists
			if err := tx.Unscoped().Delete(&item).Error; err != nil {
				tx.Rollback()
				return 0, errors.WithMessagef(err, "failed removing media item: %v", item.Id)
			} else {
				removedItems += 1
			}
		}
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		return 0, errors.New("failed committing bulk delete transaction")
	}

	return removedItems, nil
}
