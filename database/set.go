package database

import (
	"github.com/l3uddz/wantarr/pvr"
)

func (d *Database) Set(mediaId int, mediaItem *pvr.MediaItem, keepLastSearch bool) error {
	definitiveMediaItem := mediaItem

	// does this media item already exist?
	if keepLastSearch && mediaItem.LastSearch.IsZero() {
		if existingMediaItem, ok := d.vault[mediaId]; ok {
			definitiveMediaItem = &pvr.MediaItem{
				AirDateUtc: mediaItem.AirDateUtc,
				LastSearch: existingMediaItem.LastSearch,
				Name:       mediaItem.Name,
			}
		}
	}

	// set item in vault
	d.vault[mediaId] = *definitiveMediaItem
	d.changed = true
	return nil
}
