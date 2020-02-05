package database

func GetItemsCount(pvrName string, wantedType string) int {
	itemCount := 0
	db.Model(&MediaItem{}).Where("pvr_name = ? AND wanted_type = ?", pvrName, wantedType).Count(&itemCount)
	return itemCount
}
