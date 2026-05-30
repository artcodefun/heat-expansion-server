package dtos

import "github.com/google/uuid"

type GetOrderURI struct {
	ID uuid.UUID `uri:"id"`
}

type CreateOrderRequest struct {
	PackageID uuid.UUID `json:"package_id" binding:"required"`
	ReturnURL string    `json:"return_url" binding:"required"`
}
