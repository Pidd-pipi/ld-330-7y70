package middleware

import (
	"github.com/blueship581/gbemr/internal/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, _ interface{}) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": constants.CodeInternal, "message": "internal server error"})
	})
}
