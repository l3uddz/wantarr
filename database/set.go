package database

import (
	"github.com/l3uddz/wantarr/pvr"
)

func (d *Database) Set(mediaId int, mediaItem *pvr.MediaItem, keepLastSearch bool) error {
	var definitiveMediaItem *pvr.MediaItem = mediaItem

	if keepLastSearch {
		// does this media item already exist?
		if existingMediaItem, ok := d.vault[mediaId]; ok {
			definitiveMediaItem = &pvr.MediaItem{
				AirDateUtc: mediaItem.AirDateUtc,
				LastSearch: existingMediaItem.LastSearch,
				Name:       mediaItem.Name,
			}
		}
	}

	d.vault[mediaId] = *definitiveMediaItem
	d.changed = true
	return nil
}
