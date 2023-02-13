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
	tokenString, err := c.Cookie(body.AuthToken)
	if err != nil {
		if err == http.ErrNoCookie {
			response.NewErrorResponse(c, errs.NewUnauthorized(err))
			return
		}
		response.NewErrorResponse(c, errs.NewInternal(err))
		return
	}

	cookieParts := strings.Split(tokenString, " ")
	if len(cookieParts) != 2 || cookieParts[0] != "Bearer" {
		response.NewErrorResponse(c, errs.NewUnauthorized(errors.New("invalid cookie content")))
		return
	}

	token, err := jwt.ParseWithClaims(cookieParts[1], &tpportal.Claims{}, func(token *jwt.Token) (interface{}, error) {
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

	user, err := tpportal.Users(tpportal.UserWhere.ID.EQ(claims.Id)).One(c, m.st.DBSX())
	if err != nil {
		if err == sql.ErrNoRows {
			response.NewErrorResponse(c, errs.NewUnauthorized(err))
			return
		}
		response.NewErrorResponse(c, errs.NewInternal(err))
		return
	}

	c.Set("userId", claims.Id)
	c.Set("userRole", user.Role.String())
	c.Next()
}
