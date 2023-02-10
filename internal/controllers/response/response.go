package response

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/gin-gonic/gin"
	"net/http"
)

type apiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, err error) {
	myErr, ok := err.(errs.IApiError)
	if ok {
		c.AbortWithStatusJSON(myErr.Status(), apiError{
			Status:  myErr.Status(),
			Message: myErr.Error(),
		})
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
}
