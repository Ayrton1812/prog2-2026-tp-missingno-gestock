package stock

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

func (h *Handler) CrearMovimiento(c *gin.Context) {

	var request CrearMovimientoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	usuarioID := c.GetString("user_id")

	if usuarioID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "usuario no autenticado",
		})
		return
	}

	err := h.service.CrearMovimiento(
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

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "movimiento registrado correctamente",
	})
}

func (h *Handler) FindMovimientos(c *gin.Context) {

	movimientos, err := h.service.FindMovimientos(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movimientos)
}

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {

	router.POST(
		"/movimientos",
		handler.CrearMovimiento,
	)

	router.GET(
		"/movimientos",
		handler.FindMovimientos,
	)
}
