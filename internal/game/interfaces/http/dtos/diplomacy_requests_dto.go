package dtos

type diplomacyUserURI struct {
	UserID UuidStr `uri:"userId" binding:"required,uuid"`
}

type diplomacyRelationshipsQuery struct {
	Status *DiplomaticStatus `form:"status" binding:"omitempty,diplomatic_relationship_status"`
}

type diplomacyRequestURI struct {
	RequestID UuidStr `uri:"requestId" binding:"required,uuid"`
}

type diplomacyRequestActionBody struct {
	SenderBaseID int `json:"senderBaseId" binding:"required,min=1"`
}

type diplomacySendMessageBody struct {
	SenderBaseID   int                      `json:"senderBaseId" binding:"required,min=1"`
	ReceiverUserID UuidStr                  `json:"receiverUserId" binding:"required,uuid"`
	ReceiverBaseID *int                     `json:"receiverBaseId,omitempty"`
	Content        DiplomaticMessageContent `json:"content" binding:"required,diplomatic_message_content"`
}

type diplomacySendRequestBody struct {
	SenderBaseID   int                   `json:"senderBaseId" binding:"required,min=1"`
	ReceiverUserID UuidStr               `json:"receiverUserId" binding:"required,uuid"`
	ReceiverBaseID *int                  `json:"receiverBaseId,omitempty"`
	Kind           DiplomaticRequestKind `json:"kind" binding:"required,diplomatic_request_kind"`
}

type diplomacyRelationshipActionBody struct {
	SenderBaseID   int  `json:"senderBaseId" binding:"required,min=1"`
	ReceiverBaseID *int `json:"receiverBaseId,omitempty"`
}

type DiplomacyRelationshipGetRequest = Request[diplomacyUserURI, None, None]
type DiplomacyRelationshipsListRequest = Request[None, diplomacyRelationshipsQuery, None]
type DiplomacyChatRequest = Request[diplomacyUserURI, None, None]
type DiplomacyRequestActionRequest = Request[diplomacyRequestURI, None, diplomacyRequestActionBody]
type DiplomacySendInformationalMessageRequest = Request[None, None, diplomacySendMessageBody]
type DiplomacySendRequestRequest = Request[None, None, diplomacySendRequestBody]
type DiplomacyRelationshipActionRequest = Request[diplomacyUserURI, None, diplomacyRelationshipActionBody]
