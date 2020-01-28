package database

import (
	"github.com/l3uddz/wantarr/pvr"
)

func (d *Database) Set(mediaId int, mediaItem *pvr.MediaItem) error {
	d.vault[mediaId] = *mediaItem
	d.changed = true
	return nil
}
