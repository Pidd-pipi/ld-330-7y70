package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func idParam(c *gin.Context) (uint, error) {
	n, e := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(n), e
}
func page(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	s, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if p < 1 {
		p = 1
	}
	if s < 1 || s > 100 {
		s = 10
	}
	return p, s
}
func currentUser(c *gin.Context) (uint, string, string) {
	id, _ := c.Get("user_id")
	name, _ := c.Get("username")
	role, _ := c.Get("role")
	return id.(uint), name.(string), role.(string)
}
