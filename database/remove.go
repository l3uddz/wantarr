package database

func (d *Database) Delete(mediaId int) {
	if _, exists := d.vault[mediaId]; exists {
		delete(d.vault, mediaId)
		d.changed = true
	}
}
