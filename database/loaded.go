package database

func (d *Database) FromDisk() bool {
	return d.loaded
}
