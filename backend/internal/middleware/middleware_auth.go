package middleware

import (
	"errors"
	"github.com/blueship581/gbemr/internal/constants"
	"github.com/blueship581/gbemr/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Auth(s *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if len(h) < 8 || h[:7] != "Bearer " {
			abort(c, http.StatusUnauthorized, constants.CodeUnauthorized, "missing token")
			return
		}
		claims, e := s.Parse(h[7:])
		if e != nil {
			abort(c, http.StatusUnauthorized, constants.CodeUnauthorized, "invalid token")
			return
		}
		id, ok := claims["sub"].(float64)
		if !ok {
			abort(c, http.StatusUnauthorized, constants.CodeUnauthorized, "invalid token subject")
			return
		}
		username, _ := claims["username"].(string)
		role, _ := claims["role"].(string)
		c.Set("user_id", uint(id))
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}
func Roles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, _ := c.Get("role")
		for _, x := range roles {
			if r == x {
				c.Next()
				return
			}
		}
		abort(c, http.StatusForbidden, constants.CodeForbidden, "forbidden")
	}
}
func abort(c *gin.Context, status, code int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"code": code, "message": msg})
}

var _ = errors.New
var _ = strconv.Itoa
