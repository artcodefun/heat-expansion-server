package dtos

import "strings"

// storageListQuery contains query params for StorageListRequest.
type storageListQuery struct {
	Category string `form:"category" binding:"required,storage_category"`
}

// StorageListRequest captures the base ID needed for listing storage items.
type StorageListRequest = Request[BaseURI, storageListQuery, None]

// storageItemURI bundles the base and item ID for storage mutations.
type storageItemURI struct {
	BaseURI
	ItemID UuidStr `uri:"itemId" binding:"required,uuid"`
}

// StorageItemRequest bundles the base and item ID for storage mutations.
type StorageItemRequest = Request[storageItemURI, None, None]

// IsValidStorageCategory returns true if value matches one of the predefined StorageCategory constants.
func IsValidStorageCategory(value string) bool {
	upper := strings.ToUpper(value)
	switch StorageCategory(upper) {
	case StorageCategoryBuff, StorageCategoryIntel, StorageCategoryDamaged, StorageCategoryArtifact, StorageCategoryConsumable:
		return true
	default:
		return false
	}
}
