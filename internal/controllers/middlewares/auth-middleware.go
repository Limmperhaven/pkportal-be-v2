package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

func (m *MiddlewareStorage) AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader(body.AuthToken)

	headerParts := strings.Split(tokenString, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		response.NewErrorResponse(c, errs.NewUnauthorized(errors.New("invalid header content")))
		return
	}

	token, err := jwt.ParseWithClaims(headerParts[1], &tpportal.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(body.AppSalt), nil
	})
	if err != nil {
		response.NewErrorResponse(c, errs.NewUnauthorized(err))
		return
	}

	claims, ok := token.Claims.(*tpportal.Claims)
	if !ok || !token.Valid {
		response.NewErrorResponse(c, errs.NewUnauthorized(errors.New("invalid token data")))
		return
	}

	user, err := tpportal.FindUser(c, m.st.DBSX(), claims.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		response.NewErrorResponse(c, errs.NewInternal(err))
		return
	}

	c.Set("user", *user)
	c.Next()
}
