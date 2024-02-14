package controllers

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/mapper"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *ControllerStorage) SetGrades(c *gin.Context) {
	var req restmodels.SetGradesRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.SetGrades(c, *mapper.NewSetGradesRequestFromRest(&req))
	if err != nil {
		if err != nil {
			response.NewErrorResponse(c, err)
			return
		}
	}
	c.Status(http.StatusOK)
}
