package proveedor

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

	proveedores, err := h.service.FindAll(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, proveedores)
}

func (h *Handler) FindByID(c *gin.Context) {

	proveedor, err := h.service.FindByID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, proveedor)
}

func (h *Handler) Create(c *gin.Context) {

	var request CrearProveedorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	usuarioID := c.GetString("user_id")

	proveedor, err := h.service.Create(
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

	c.JSON(http.StatusCreated, proveedor)
}

func (h *Handler) Update(c *gin.Context) {

	var request ActualizarProveedorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	usuarioID := c.GetString("user_id")

	proveedor, err := h.service.Update(
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

	c.JSON(http.StatusOK, proveedor)
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

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {

	router.GET("/proveedores", handler.FindAll)
	router.GET("/proveedores/:id", handler.FindByID)
	router.POST("/proveedores", handler.Create)
	router.PUT("/proveedores/:id", handler.Update)
	router.DELETE("/proveedores/:id", handler.Delete)
}
