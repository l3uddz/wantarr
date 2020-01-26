package database

import (
	"github.com/asdine/storm/v3"
	"github.com/l3uddz/wantarr/logger"
	stringutils "github.com/l3uddz/wantarr/utils/strings"
	"github.com/pkg/errors"
)

var (
	log           = logger.GetLogger("db")
	Db  *storm.DB = nil
)

func Init(databasePath string) error {
	if db, err := storm.Open(databasePath); err != nil {
		return errors.Wrapf(err, "failed opening database: %q", databasePath)
	} else {
		Db = db
	}

	log.Infof("Using %s = %q", stringutils.StringLeftJust("DATABASE", " ", 10), databasePath)
	return nil
}

func Close() {
	if Db == nil {
		return
	}

	if err := Db.Close(); err != nil {
		log.WithError(err).Error("failed closing database gracefully...")
	}
}
