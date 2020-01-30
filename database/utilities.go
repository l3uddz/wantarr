package database

import (
	"github.com/l3uddz/wantarr/pvr"
	"github.com/pkg/errors"
)

/* Public */

func (d *Database) RemoveMissingMediaItems(latestMediaItems map[int]pvr.MediaItem) (int, error) {
	removedItems := 0

	for itemId, _ := range *d.GetVault() {
		if _, itemStillMissing := latestMediaItems[itemId]; !itemStillMissing {
			// this item is no longer missing
			d.Delete(itemId)
			removedItems += 1
		}
	}

	return removedItems, nil
}

func (d *Database) SetMediaItems(latestMediaItems map[int]pvr.MediaItem) (int, error) {
	newItems := 0
	var err error = nil

	for itemId, record := range latestMediaItems {
		// does item already exist
		if _, itemExists := d.vault[itemId]; !itemExists {
			newItems += 1
		}

		if err := d.Set(itemId, &record, true); err != nil {
			err = errors.WithMessagef(err, "failed stashing mediaItem %q in database", record)
			log.WithError(err).Errorf("Failed stashing mediaItem %q in database", record)
		}
	}

	return newItems, err
}
