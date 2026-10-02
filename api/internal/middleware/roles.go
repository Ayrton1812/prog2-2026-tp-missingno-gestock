package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Roles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		rolValue, exists := c.Get("rol")

		if !exists {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "rol no encontrado",
				},
			)
			return
		}

		rol, ok := rolValue.(string)

		if !ok {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "rol invalido",
				},
			)
			return
		}

		for _, rolPermitido := range roles {
			if rol == rolPermitido {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error": "no tiene permisos para realizar esta operacion",
			},
		)
	}
}
