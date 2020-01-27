package database

import (
	"github.com/pkg/errors"
	"io/ioutil"
)

func (d *Database) Close() error {
	// dont dump database when no changes made
	if !d.changed {
		return nil
	}

	// marshal vault
	jsonData, err := json.MarshalIndent(d.vault, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal database data")
	}

	// write json data to file
	if err := ioutil.WriteFile(d.filePath, jsonData, 0644); err != nil {
		return errors.Wrapf(err, "failed to write marshalled database data to: %q", d.filePath)
	}

	return nil
}
