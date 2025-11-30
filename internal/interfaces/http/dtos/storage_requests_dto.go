package dtos

// StorageItemURI binds /bases/:baseId/storage/:itemId style routes.
type StorageItemURI struct {
	BaseID int    `uri:"baseId" binding:"required,min=1"`
	ItemID string `uri:"itemId" binding:"required"`
}
