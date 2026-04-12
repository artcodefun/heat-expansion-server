package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type DiplomacyHandler struct {
	queries    cqrs.DiplomacyQueries
	commands   cqrs.DiplomacyCommands
	translator ports.Translator
}

func NewDiplomacyHandler(queries cqrs.DiplomacyQueries, commands cqrs.DiplomacyCommands, translator ports.Translator) *DiplomacyHandler {
	return &DiplomacyHandler{queries: queries, commands: commands, translator: translator}
}

// ListRelationships handles GET /diplomacy/relationships.
func (h *DiplomacyHandler) ListRelationships(c *gin.Context) {
	var req dtos.DiplomacyRelationshipsListRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	items, err := h.queries.ListRelationships(c.Request.Context(), actor, dtos.DiplomaticStatusPtrFromDTO(req.Query.Status))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.DiplomaticRelationshipsFromReadModels(items))
}

// GetRelationship handles GET /diplomacy/relationships/:userId.
func (h *DiplomacyHandler) GetRelationship(c *gin.Context) {
	var req dtos.DiplomacyRelationshipGetRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	item, err := h.queries.GetRelationship(c.Request.Context(), actor, req.Uri.UserID.Uuid())
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.DiplomaticRelationshipFromReadModel(item))
}

// ListAvailableInformationalMessages handles GET /diplomacy/messages/available.
func (h *DiplomacyHandler) ListAvailableInformationalMessages(c *gin.Context) {
	locale := getLocale(c)
	items := domain.InformationalDiplomaticMessageContents()
	result := make(map[string]string, len(items))
	for _, key := range items {
		result[string(key)] = h.translator.T(locale, key, nil)
	}
	c.JSON(http.StatusOK, result)
}

// ListChats handles GET /diplomacy/messages/chats.
func (h *DiplomacyHandler) ListChats(c *gin.Context) {
	actor := actor(c)
	items, err := h.queries.ListChats(c.Request.Context(), actor)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.DiplomaticChatsFromReadModels(items, h.translator, getLocale(c)))
}

// GetUnreadCount handles GET /diplomacy/messages/unread-count.
func (h *DiplomacyHandler) GetUnreadCount(c *gin.Context) {
	actor := actor(c)
	count, err := h.queries.GetUnreadMessagesCount(c.Request.Context(), actor)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// ListChatMessages handles GET /diplomacy/messages/chats/:userId.
func (h *DiplomacyHandler) ListChatMessages(c *gin.Context) {
	var req dtos.DiplomacyChatRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	items, err := h.queries.ListChatMessages(c.Request.Context(), actor, req.Uri.UserID.Uuid())
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.DiplomaticMessagesFromReadModels(items, h.translator, getLocale(c)))
}

// MarkChatAsRead handles POST /diplomacy/messages/chats/:userId/read.
func (h *DiplomacyHandler) MarkChatAsRead(c *gin.Context) {
	var req dtos.DiplomacyChatRequest
	if !bindRequest(c, &req) {
		return
	}
	if err := h.commands.MarkChatAsRead(c.Request.Context(), actor(c), req.Uri.UserID.Uuid()); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}

// ListPendingRequests handles GET /diplomacy/requests/pending.
func (h *DiplomacyHandler) ListPendingRequests(c *gin.Context) {
	actor := actor(c)
	items, err := h.queries.ListPendingRequests(c.Request.Context(), actor)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.DiplomaticRequestsFromReadModels(items))
}

// SendInformationalMessage handles POST /diplomacy/messages.
func (h *DiplomacyHandler) SendInformationalMessage(c *gin.Context) {
	var req dtos.DiplomacySendInformationalMessageRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	id, err := h.commands.SendInformationalMessage(c.Request.Context(), actor, req.Body.SenderBaseID, req.Body.ReceiverUserID.Uuid(), req.Body.ReceiverBaseID, domain.TranslationKey(req.Body.Content))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// SendRequest handles POST /diplomacy/requests.
func (h *DiplomacyHandler) SendRequest(c *gin.Context) {
	var req dtos.DiplomacySendRequestRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	id, err := h.commands.SendRequest(c.Request.Context(), actor, req.Body.SenderBaseID, req.Body.ReceiverUserID.Uuid(), req.Body.ReceiverBaseID, domain.DiplomaticRequestKind(req.Body.Kind))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// DeclareWar handles POST /diplomacy/relationships/:userId/declare-war.
func (h *DiplomacyHandler) DeclareWar(c *gin.Context) {
	var req dtos.DiplomacyRelationshipActionRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	id, err := h.commands.DeclareWar(c.Request.Context(), actor, req.Body.SenderBaseID, req.Uri.UserID.Uuid(), req.Body.ReceiverBaseID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// BreakAlliance handles POST /diplomacy/relationships/:userId/break-alliance.
func (h *DiplomacyHandler) BreakAlliance(c *gin.Context) {
	var req dtos.DiplomacyRelationshipActionRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	id, err := h.commands.BreakAlliance(c.Request.Context(), actor, req.Body.SenderBaseID, req.Uri.UserID.Uuid(), req.Body.ReceiverBaseID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// AcceptRequest handles POST /diplomacy/requests/:requestId/accept.
func (h *DiplomacyHandler) AcceptRequest(c *gin.Context) {
	var req dtos.DiplomacyRequestActionRequest
	if !bindRequest(c, &req) {
		return
	}
	if err := h.commands.AcceptRequest(c.Request.Context(), actor(c), req.Body.SenderBaseID, req.Uri.RequestID.Uuid()); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}

// RejectRequest handles POST /diplomacy/requests/:requestId/reject.
func (h *DiplomacyHandler) RejectRequest(c *gin.Context) {
	var req dtos.DiplomacyRequestActionRequest
	if !bindRequest(c, &req) {
		return
	}
	if err := h.commands.RejectRequest(c.Request.Context(), actor(c), req.Body.SenderBaseID, req.Uri.RequestID.Uuid()); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}
