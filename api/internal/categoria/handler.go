package categoria

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

func (h *Handler) List(c *gin.Context) {
	categorias, err := h.service.ListarTodas(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	dtos := make([]CategoriaDTO, 0, len(categorias))

	for _, categoria := range categorias {
		dtos = append(dtos, categoria.ToDTO())
	}

	c.JSON(http.StatusOK, dtos)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	categoria, err := h.service.BuscarPorID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, categoria.ToDTO())
}

func (h *Handler) Create(c *gin.Context) {
	var dto CategoriaDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	categoria, err := dto.ToModel()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	creada, err := h.service.Crear(
		c.Request.Context(),
		categoria,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, creada.ToDTO())
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var dto CategoriaDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	categoria, err := dto.ToModel()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	actualizada, err := h.service.Actualizar(
		c.Request.Context(),
		id,
		categoria,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, actualizada.ToDTO())
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Eliminar(
		c.Request.Context(),
		id,
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	categorias := router.Group("/api/categorias")

	{
		categorias.GET("", handler.List)
		categorias.GET("/:id", handler.GetByID)
		categorias.POST("", handler.Create)
		categorias.PUT("/:id", handler.Update)
		categorias.DELETE("/:id", handler.Delete)
	}
}
