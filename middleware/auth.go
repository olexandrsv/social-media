package middleware

import (
	"log"
	"social-media/auth"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.String(403, "no credentials")
		c.Abort()
		return
	}

	if err := auth.ValidateToken(token); err != nil {
		log.Println(err)
		c.String(403, "invalid credentials")
		c.Abort()
		return
	}
	c.Next()
}
