package usuario

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

func (h *Handler) Login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.service.Login(
		c.Request.Context(),
		request,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) List(c *gin.Context) {
	usuarios, err := h.service.ListarTodas(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	dtos := make([]UsuarioDTO, 0, len(usuarios))

	for _, usuario := range usuarios {
		dtos = append(dtos, usuario.ToDTO())
	}

	c.JSON(http.StatusOK, dtos)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	usuario, err := h.service.BuscarPorID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, usuario.ToDTO())
}

func (h *Handler) Create(c *gin.Context) {
	var request CrearUsuarioRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	usuario, err := h.service.Crear(
		c.Request.Context(),
		request,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, usuario.ToDTO())
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var request ActualizarUsuarioRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	usuario, err := h.service.Actualizar(
		c.Request.Context(),
		id,
		request,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, usuario.ToDTO())
}

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
) {
	router.POST(
		"/api/login",
		handler.Login,
	)

	usuarios := router.Group("/api/usuarios")

	{
		usuarios.GET("", handler.List)
		usuarios.GET("/:id", handler.GetByID)
		usuarios.POST("", handler.Create)
		usuarios.PUT("/:id", handler.Update)
	}
}
