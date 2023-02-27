package controllers

import (
	"fmt"
	"github.com/Limmperhaven/pkportal-be-v2/internal/body"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/response"
	"github.com/Limmperhaven/pkportal-be-v2/internal/errs"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/mapper"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func (s *ControllerStorage) CreateUser(c *gin.Context) {
	var req restmodels.CreateUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.CreateUser(c, mapper.NewCreateUserRequestFromRest(&req))
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) GetUser(c *gin.Context) {
	userIdParam := c.Param("id")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный userId: %s", userIdParam)))
		return
	}
	user, err := s.uc.GetUser(c, userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	res := mapper.NewGetUserResponseToRest(&user)
	c.JSON(http.StatusOK, *res)
}

func (s *ControllerStorage) ListUsers(c *gin.Context) {
	var req restmodels.UserFilter
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	res, err := s.uc.ListUsers(c, *mapper.NewListUsersRequestFromRest(&req))
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.NewGetUserResponseArrayToRest(res))
}

func (s *ControllerStorage) UpdateUser(c *gin.Context) {
	userIdParam := c.Param("id")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный userId: %s", userIdParam)))
		return
	}
	var req restmodels.UpdateUserRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	err = s.uc.UpdateUser(c, *mapper.NewUpdateUserRequestFromRest(&req), userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) GetMe(c *gin.Context) {
	userIdCtx, ok := c.Get(body.UserCtx)
	if !ok {
		response.NewErrorResponse(c, errs.NewInternal(errors.New("в контексте отсутствует userId")))
		return
	}
	userId := userIdCtx.(tpportal.User).ID
	user, err := s.uc.GetUser(c, userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	res := mapper.NewGetUserResponseToRest(&user)
	c.JSON(http.StatusOK, *res)
}

func (s *ControllerStorage) ListStatuses(c *gin.Context) {
	var req restmodels.ListStatusesRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	statuses, err := s.uc.ListStatuses(c, *mapper.NewListStatusesRequestFromRest(&req))
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.NewIdNameArrayToRest(statuses))
}

func (s *ControllerStorage) SetUserStatus(c *gin.Context) {
	userIdParam := c.Param("userId")
	statusIdParam := c.Param("statusId")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный userId: %s", userIdParam)))
		return
	}
	statusId, err := strconv.ParseInt(statusIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный statusId: %s", userIdParam)))
		return
	}
	err = s.uc.SetUserStatus(c, userId, statusId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) UploadScreenshot(c *gin.Context) {
	screenType := c.PostForm("type")
	fmt.Println(screenType)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(err))
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		response.NewErrorResponse(c, errs.NewInternal(fmt.Errorf("ошибка при открытии файла: %s", err.Error())))
		return
	}
	fileData, err := io.ReadAll(file)
	if err != nil {
		response.NewErrorResponse(c, errs.NewInternal(fmt.Errorf("ошибка при чтении файла: %s", err.Error())))
		return
	}
	err = s.uc.UploadScreenshot(c, *mapper.NewUploadScreenshotRequestFromRest(
		fileHeader.Filename, screenType, fileHeader.Size, fileData))
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) DownloadScreenshot(c *gin.Context) {
	userIdParam := c.Param("userId")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный userId: %s", userIdParam)))
		return
	}
	res, err := s.uc.DownloadScreenshot(c, userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	//c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", res.FileName))
	//c.Data(http.StatusOK, res.ContentType, res.FileContent)
	c.JSON(http.StatusOK, *mapper.NewDownloadFileResponseToRest(&res))
}

func (s *ControllerStorage) DownloadMyScreenshot(c *gin.Context) {
	userIdCtx, ok := c.Get(body.UserCtx)
	if !ok {
		response.NewErrorResponse(c, errs.NewInternal(errors.New("в контексте отсутствует userId")))
		return
	}
	userId := userIdCtx.(tpportal.User).ID
	res, err := s.uc.DownloadScreenshot(c, userId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, *mapper.NewDownloadFileResponseToRest(&res))
}

func (s *ControllerStorage) DownloadRegistrationList(c *gin.Context) {
	tdIdParam := c.Param("tdId")
	tdId, err := strconv.ParseInt(tdIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id: %s", tdIdParam)))
		return
	}
	res, err := s.uc.DownloadRegistrationList(c, tdId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, *mapper.NewDownloadFileResponseToRest(&res))
}

func (s *ControllerStorage) ExportToXlsx(c *gin.Context) {
	tdIdParam := c.Param("tdId")
	tdId, err := strconv.ParseInt(tdIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id: %s", tdIdParam)))
		return
	}
	res, err := s.uc.ExportTestDateToXlsx(c, tdId)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, *mapper.NewDownloadFileResponseToRest(&res))
}

func (s *ControllerStorage) SetUserRole(c *gin.Context) {
	userIdParam := c.Param("userId")
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		response.NewErrorResponse(c, errs.NewBadRequest(fmt.Errorf("невалидный id: %s", userIdParam)))
		return
	}
	role := c.Param("role")
	err = s.uc.SetUserRole(c, userId, role)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (s *ControllerStorage) ResendActivationEmail(c *gin.Context) {
	err := s.uc.ResendActivationEmail(c)
	if err != nil {
		response.NewErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}
