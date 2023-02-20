package mapper

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
)

func NewSignUpRequestFromRest(in *restmodels.SignUpRequest) *tpportal.SignUpRequest {
	return &tpportal.SignUpRequest{
		Email:             in.Email,
		Password:          in.Password,
		Fio:               in.Fio,
		DateOfBirth:       in.DateOfBirth,
		Gender:            in.Gender,
		PhoneNumber:       in.PhoneNumber,
		ParentPhoneNumber: in.ParentPhoneNumber,
		CurrentSchool:     in.CurrentSchool,
		EducationYear:     in.EducationYear,
	}
}

func NewSignInRequestFromRest(in *restmodels.SignInRequest) *tpportal.SignInRequest {
	return &tpportal.SignInRequest{
		Email:    in.Email,
		Password: in.Password,
	}
}

func NewSignInResponseToRest(in *tpportal.SignInResponse) *restmodels.SignInResponse {
	return &restmodels.SignInResponse{
		Id:                in.User.ID,
		Email:             in.User.Email,
		Fio:               in.User.Fio,
		DateOfBirth:       in.User.DateOfBirth.Format("02.01.2006"),
		Gender:            in.User.Gender.String(),
		PhoneNumber:       in.User.PhoneNumber,
		ParentPhoneNumber: in.User.ParentPhoneNumber,
		CurrentSchool:     in.User.CurrentSchool.String,
		EducationYear:     int64(in.User.EducationYear),
		Role:              in.User.Role.String(),
		StatusId:          in.User.StatusID,
		IsActivated:       in.User.IsActivated,
		AuthToken:         in.AuthToken,
	}
}

func NewCreateUserRequestFromRest(in *restmodels.CreateUserRequest) *tpportal.CreateUserRequest {
	return &tpportal.CreateUserRequest{
		Email:             in.Email,
		Fio:               in.Fio,
		Password:          in.Password,
		DateOfBirth:       in.DateOfBirth,
		Gender:            in.Gender,
		PhoneNumber:       in.PhoneNumber,
		ParentPhoneNumber: in.ParentPhoneNumber,
		CurrentSchool:     in.CurrentSchool,
		EducationYear:     in.EducationYear,
		IsActivated:       in.IsActivated,
		Role:              in.Role,
		StatusId:          in.StatusId,
	}
}
