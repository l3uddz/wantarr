package database

import (
	"github.com/pkg/errors"
	"time"
)

func GetMediaItems(pvrName string, wantedType string, excludeFuture bool) ([]MediaItem, error) {
	var mediaItems []MediaItem

	// generate query
	sqlQuery := "pvr_name = ? AND wanted_type = ?"
	sqlParams := []interface{}{
		pvrName,
		wantedType,
	}

	if excludeFuture {
		sqlQuery += " AND air_date_utc <= ?"
		sqlParams = append(sqlParams, time.Now().UTC())
	}

	// exec query
	if err := db.Where(sqlQuery, sqlParams...).Order("air_date_utc desc").Find(&mediaItems).Error; err != nil {
		return nil, errors.Wrap(err, "failed querying for media items")
	}

	return mediaItems, nil
}
