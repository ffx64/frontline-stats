package middleware

import (
	"log"

	"github.com/ffx64/gamestats-backend/internal/errors"
	"github.com/gin-gonic/gin"
)

func Auth(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			log.Println("[middleware:auth] Authorization header não encontrado")
			c.JSON(errors.ErrHeaderMissing.Status, errors.ErrHeaderMissing)
			c.Abort()
			return
		}

		if authHeader != key {
			log.Println("[middleware:auth] Chave de autorização inválida")
			c.JSON(errors.ErrUnauthorized.Status, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}
