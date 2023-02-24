package middlewares

import (
	"errors"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *MiddlewareStorage) CheckAdminRoleMiddleware(c *gin.Context) {
	userCtx, ok := c.Get(body.UserCtx)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	user, ok := userCtx.(tpportal.User)
	if !ok {
		response.NewErrorResponse(c, errs.NewInternal(errors.New("неправильный формат пользователя в контексте")))
		return
	}
	if user.Role != tpportal.UserRoleAdmin {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}

func (m *MiddlewareStorage) CheckActivationMiddleware(c *gin.Context) {
	userCtx, ok := c.Get(body.UserCtx)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	user, ok := userCtx.(tpportal.User)
	if !ok {
		response.NewErrorResponse(c, errs.NewInternal(errors.New("неправильный формат пользователя в контексте")))
		return
	}
	if !user.IsActivated {
		response.NewErrorResponse(c, errs.NewForbidden(errors.New("для доступа к ресурсу активируйте аккаунт")))
		return
	}
	c.Next()
}
