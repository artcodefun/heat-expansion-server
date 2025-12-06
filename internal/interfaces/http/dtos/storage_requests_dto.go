package dtos

// StorageListRequest captures the base ID needed for listing storage items.
type StorageListRequest = Request[BaseURI, None, None]

// storageItemURI bundles the base and item ID for storage mutations.
type storageItemURI struct {
	BaseURI
	ItemID UuidStr `uri:"itemId" binding:"required,uuid"`
}

// StorageItemRequest bundles the base and item ID for storage mutations.
type StorageItemRequest = Request[storageItemURI, None, None]
