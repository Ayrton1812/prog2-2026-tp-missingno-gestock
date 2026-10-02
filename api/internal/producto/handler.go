package producto

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
	productos, err := h.service.FindAll(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, productos)
}

func (h *Handler) FindByID(c *gin.Context) {
	producto, err := h.service.FindByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, producto)
}

func (h *Handler) Create(c *gin.Context) {
	var request CrearProductoRequest

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

	producto, err := h.service.Create(
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

	c.JSON(http.StatusCreated, producto)
}

func (h *Handler) Update(c *gin.Context) {
	var request ActualizarProductoRequest

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

	producto, err := h.service.Update(
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

	c.JSON(http.StatusOK, producto)
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
	router.GET("/api/productos", handler.FindAll)
	router.GET("/api/productos/:id", handler.FindByID)
	router.POST("/api/productos", handler.Create)
	router.PUT("/api/productos/:id", handler.Update)
	router.DELETE("/api/productos/:id", handler.Delete)
}
