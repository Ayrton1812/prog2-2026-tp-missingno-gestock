package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		if header == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "token requerido",
				},
			)
			return
		}

		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "formato de token invalido",
				},
			)
			return
		}

		token, err := jwt.Parse(
			parts[1],
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "token invalido",
				},
			)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "claims invalidos",
				},
			)
			return
		}

		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}

		if rol, ok := claims["rol"].(string); ok {
			c.Set("rol", rol)
		}

		if depositoID, ok := claims["deposito_id"].(string); ok {
			c.Set("deposito_id", depositoID)
		}

		c.Next()
	}
}
