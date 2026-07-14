package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

// PrototypeHandler handles CRUD endpoints for all four prototype types.
type PrototypeHandler struct {
	commands   cqrs.PrototypeCommands
	queries    cqrs.PrototypeQueries
	translator ports.Translator
}

func NewPrototypeHandler(commands cqrs.PrototypeCommands, queries cqrs.PrototypeQueries, translator ports.Translator) *PrototypeHandler {
	return &PrototypeHandler{commands: commands, queries: queries, translator: translator}
}

// ── Army ──────────────────────────────────────────────────────────────────────

func (h *PrototypeHandler) ListArmy(c *gin.Context) {
	list, err := h.queries.ListArmyPrototypes(c.Request.Context(), actor(c))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	out := make([]dtos.ArmyPrototypeDTO, len(list))
	for i, p := range list {
		out[i] = dtos.ArmyPrototypeDTOFromModel(p)
	}
	c.JSON(http.StatusOK, out)
}

func (h *PrototypeHandler) GetArmy(c *gin.Context) {
	var uri dtos.GetPrototypeURI
	if !bindURI(c, &uri) {
		return
	}
	p, err := h.queries.GetArmyPrototype(c.Request.Context(), actor(c), uri.ID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ArmyPrototypeDTOFromModel(p))
}

func (h *PrototypeHandler) CreateArmy(c *gin.Context) {
	var req dtos.ArmyPrototypeDTO
	if !bindRequest(c, &req) {
		return
	}
	p, err := h.commands.CreateArmyPrototype(c.Request.Context(), actor(c), dtos.ArmyPrototypeDTOToModel(req))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusCreated, dtos.ArmyPrototypeDTOFromModel(p))
}

func (h *PrototypeHandler) UpdateArmy(c *gin.Context) {
	var uri dtos.GetPrototypeURI
	if !bindURI(c, &uri) {
		return
	}
	var req dtos.ArmyPrototypeDTO
	if !bindRequest(c, &req) {
		return
	}
	req.ID = uri.ID
	p, err := h.commands.UpdateArmyPrototype(c.Request.Context(), actor(c), dtos.ArmyPrototypeDTOToModel(req))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ArmyPrototypeDTOFromModel(p))
}

// ── Build ──────────────────────────────────────────────────────────────────────

func (h *PrototypeHandler) ListBuild(c *gin.Context) {
	list, err := h.queries.ListBuildPrototypes(c.Request.Context(), actor(c))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	out := make([]dtos.BuildPrototypeDTO, len(list))
	for i, p := range list {
		out[i] = dtos.BuildPrototypeDTOFromModel(p)
	}
	c.JSON(http.StatusOK, out)
}

func (h *PrototypeHandler) GetBuild(c *gin.Context) {
	var uri dtos.GetPrototypeURI
	if !bindURI(c, &uri) {
		return
	}
	p, err := h.queries.GetBuildPrototype(c.Request.Context(), actor(c), uri.ID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.BuildPrototypeDTOFromModel(p))
}

func (h *PrototypeHandler) CreateBuild(c *gin.Context) {
	var req dtos.BuildPrototypeDTO
	if !bindRequest(c, &req) {
		return
	}
	p, err := h.commands.CreateBuildPrototype(c.Request.Context(), actor(c), dtos.BuildPrototypeDTOToModel(req))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusCreated, dtos.BuildPrototypeDTOFromModel(p))
}

func (h *PrototypeHandler) UpdateBuild(c *gin.Context) {
	var uri dtos.GetPrototypeURI
	if !bindURI(c, &uri) {
		return
	}
	var req dtos.BuildPrototypeDTO
	if !bindRequest(c, &req) {
		return
	}
	req.ID = uri.ID
	p, err := h.commands.UpdateBuildPrototype(c.Request.Context(), actor(c), dtos.BuildPrototypeDTOToModel(req))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.BuildPrototypeDTOFromModel(p))
}

// ── Storage ───────────────────────────────────────────────────────────────────

func (h *PrototypeHandler) ListStorage(c *gin.Context) {
	list, err := h.queries.ListStoragePrototypes(c.Request.Context(), actor(c))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	out := make([]dtos.StoragePrototypeDTO, len(list))
	for i, p := range list {
		out[i] = dtos.StoragePrototypeDTOFromModel(p)
	}
	c.JSON(http.StatusOK, out)
}

func (h *PrototypeHandler) GetStorage(c *gin.Context) {
	var uri dtos.GetPrototypeURI
	if !bindURI(c, &uri) {
		return
	}
	p, err := h.queries.GetStoragePrototype(c.Request.Context(), actor(c), uri.ID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.StoragePrototypeDTOFromModel(p))
}

func (h *PrototypeHandler) CreateStorage(c *gin.Context) {
	var req dtos.StoragePrototypeDTO
	if !bindRequest(c, &req) {
		return
	}
	p, err := h.commands.CreateStoragePrototype(c.Request.Context(), actor(c), dtos.StoragePrototypeDTOToModel(req))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusCreated, dtos.StoragePrototypeDTOFromModel(p))
}

func (h *PrototypeHandler) UpdateStorage(c *gin.Context) {
	var uri dtos.GetPrototypeURI
	if !bindURI(c, &uri) {
		return
	}
	var req dtos.StoragePrototypeDTO
	if !bindRequest(c, &req) {
		return
	}
	req.ID = uri.ID
	p, err := h.commands.UpdateStoragePrototype(c.Request.Context(), actor(c), dtos.StoragePrototypeDTOToModel(req))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.StoragePrototypeDTOFromModel(p))
}

// ── Tech ──────────────────────────────────────────────────────────────────────

func (h *PrototypeHandler) ListTech(c *gin.Context) {
	list, err := h.queries.ListTechPrototypes(c.Request.Context(), actor(c))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	out := make([]dtos.TechPrototypeDTO, len(list))
	for i, p := range list {
		out[i] = dtos.TechPrototypeDTOFromModel(p)
	}
	c.JSON(http.StatusOK, out)
}

func (h *PrototypeHandler) GetTech(c *gin.Context) {
	var uri dtos.GetPrototypeURI
	if !bindURI(c, &uri) {
		return
	}
	p, err := h.queries.GetTechPrototype(c.Request.Context(), actor(c), uri.ID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechPrototypeDTOFromModel(p))
}

func (h *PrototypeHandler) CreateTech(c *gin.Context) {
	var req dtos.TechPrototypeDTO
	if !bindRequest(c, &req) {
		return
	}
	p, err := h.commands.CreateTechPrototype(c.Request.Context(), actor(c), dtos.TechPrototypeDTOToModel(req))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusCreated, dtos.TechPrototypeDTOFromModel(p))
}

func (h *PrototypeHandler) UpdateTech(c *gin.Context) {
	var uri dtos.GetPrototypeURI
	if !bindURI(c, &uri) {
		return
	}
	var req dtos.TechPrototypeDTO
	if !bindRequest(c, &req) {
		return
	}
	req.ID = uri.ID
	p, err := h.commands.UpdateTechPrototype(c.Request.Context(), actor(c), dtos.TechPrototypeDTOToModel(req))
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechPrototypeDTOFromModel(p))
}
