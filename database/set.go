package database

import "time"

func (d *Database) Set(mediaId int, expiry *time.Time) error {
	d.vault[mediaId] = expiry
	d.changed = true
	return nil
}
