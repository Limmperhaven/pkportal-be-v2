package middlewares

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *MiddlewareStorage) CheckAdminRoleMiddleware(c *gin.Context) {
	if role, ok := c.Get(body.UserRole); !ok || role != tpportal.UserRoleAdmin.String() {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}
