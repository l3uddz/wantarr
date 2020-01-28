package database

import (
	"fmt"
	"github.com/l3uddz/wantarr/pvr"
)

func (d *Database) Get(mediaId int) (*pvr.MediaItem, error) {
	mediaItem, ok := d.vault[mediaId]
	if !ok {
		return nil, fmt.Errorf("mediaItem not found with id: %q", mediaId)
	}
	return &mediaItem, nil
}

func (d *Database) GetVault() *map[int]pvr.MediaItem {
	return &d.vault
}
