package deposito

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) FindAll(c *gin.Context) {
	depositos, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, depositos)
}

func (h *Handler) FindByID(c *gin.Context) {
	deposito, err := h.service.FindByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, deposito)
}

func (h *Handler) Create(c *gin.Context) {
	var request CrearDepositoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	usuarioID := c.GetString("user_id")
	if usuarioID == "" {
		usuarioID = "sistema"
	}

	deposito, err := h.service.Create(
		c.Request.Context(),
		request,
		usuarioID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, deposito)
}

func (h *Handler) Update(c *gin.Context) {
	var request ActualizarDepositoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	usuarioID := c.GetString("user_id")
	if usuarioID == "" {
		usuarioID = "sistema"
	}

	deposito, err := h.service.Update(
		c.Request.Context(),
		c.Param("id"),
		request,
		usuarioID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, deposito)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.service.Delete(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	router.GET("/api/depositos", handler.FindAll)
	router.GET("/api/depositos/:id", handler.FindByID)
	router.POST("/api/depositos", handler.Create)
	router.PUT("/api/depositos/:id", handler.Update)
	router.DELETE("/api/depositos/:id", handler.Delete)
}
