package ordencompra

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
	var req CrearOrdenCompraRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	rol := c.GetString("rol")

	orden, err := h.service.Crear(
		c.Request.Context(),
		req,
		c.GetString("user_id"),
		rol,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, orden)
}

func (h *Handler) Listar(c *gin.Context) {
	ordenes, err := h.service.Listar(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ordenes)
}

func (h *Handler) BuscarPorID(c *gin.Context) {
	orden, err := h.service.BuscarPorID(
		c.Request.Context(),
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orden)
}

func (h *Handler) Actualizar(c *gin.Context) {
	var req ActualizarOrdenCompraRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.Actualizar(
		c.Request.Context(),
		c.Param("id"),
		req,
		c.GetString("user_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "orden actualizada correctamente",
	})
}

func (h *Handler) Confirmar(c *gin.Context) {
	err := h.service.Confirmar(
		c.Request.Context(),
		c.Param("id"),
		c.GetString("user_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "orden confirmada correctamente",
	})
}

func (h *Handler) Recibir(c *gin.Context) {
	var req RecibirOrdenCompraRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.Recibir(
		c.Request.Context(),
		c.Param("id"),
		req,
		c.GetString("user_id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "mercadería recibida correctamente",
	})
}

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {
	router.GET("/ordenes-compra", handler.Listar)
	router.GET("/ordenes-compra/:id", handler.BuscarPorID)

	router.POST("/ordenes-compra", handler.Crear)

	router.PUT(
		"/ordenes-compra/:id",
		handler.Actualizar,
	)

	router.PUT(
		"/ordenes-compra/:id/confirmar",
		handler.Confirmar,
	)

	router.POST(
		"/ordenes-compra/:id/recibir",
		handler.Recibir,
	)
}
