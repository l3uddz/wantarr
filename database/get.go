package database

import "github.com/pkg/errors"

func GetMediaItems(pvrName string, wantedType string) ([]MediaItem, error) {
	var mediaItems []MediaItem

	if err := db.Where(MediaItem{
		PvrName:    pvrName,
		WantedType: wantedType,
	}).Order("air_date_utc desc").Find(&mediaItems).Error; err != nil {
		return nil, errors.Wrap(err, "failed querying for media items")
	}

	return mediaItems, nil
}
