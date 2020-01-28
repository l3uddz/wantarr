package database

import "github.com/l3uddz/wantarr/pvr"

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
