package transferencia

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

func (h *Handler) Crear(c *gin.Context) {

	var request CrearTransferenciaRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	usuarioID := c.GetString("user_id")
	rol := c.GetString("rol")
	depositoID := c.GetString("deposito_id")

	transferencia, err := h.service.Crear(
		c.Request.Context(),
		request,
		usuarioID,
		rol,
		depositoID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, transferencia)
}

func (h *Handler) Aprobar(c *gin.Context) {

	usuarioID := c.GetString("user_id")
	rol := c.GetString("rol")
	depositoID := c.GetString("deposito_id")

	transferencia, err := h.service.Aprobar(
		c.Request.Context(),
		c.Param("id"),
		usuarioID,
		rol,
		depositoID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transferencia)
}

func (h *Handler) Rechazar(c *gin.Context) {

	var request RechazarTransferenciaRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	usuarioID := c.GetString("user_id")
	rol := c.GetString("rol")
	depositoID := c.GetString("deposito_id")

	transferencia, err := h.service.Rechazar(
		c.Request.Context(),
		c.Param("id"),
		request.Motivo,
		usuarioID,
		rol,
		depositoID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transferencia)
}

func (h *Handler) Completar(c *gin.Context) {

	usuarioID := c.GetString("user_id")
	rol := c.GetString("rol")
	depositoID := c.GetString("deposito_id")

	transferencia, err := h.service.Completar(
		c.Request.Context(),
		c.Param("id"),
		usuarioID,
		rol,
		depositoID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transferencia)
}

func (h *Handler) FindAll(c *gin.Context) {

	transferencias, err := h.service.FindAll(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transferencias)
}

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {

	router.POST(
		"/transferencias",
		handler.Crear,
	)

	router.GET(
		"/transferencias",
		handler.FindAll,
	)

	router.PUT(
		"/transferencias/:id/aprobar",
		handler.Aprobar,
	)

	router.PUT(
		"/transferencias/:id/rechazar",
		handler.Rechazar,
	)

	router.PUT(
		"/transferencias/:id/completar",
		handler.Completar,
	)
}
