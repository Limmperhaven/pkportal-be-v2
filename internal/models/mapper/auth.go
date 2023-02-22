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
		Id:                in.Id,
		Email:             in.Email,
		Fio:               in.Fio,
		DateOfBirth:       in.DateOfBirth,
		Gender:            in.Gender,
		PhoneNumber:       in.PhoneNumber,
		ParentPhoneNumber: in.ParentPhoneNumber,
		CurrentSchool:     in.CurrentSchool,
		EducationYear:     in.EducationYear,
		Role:              in.Role,
		Status:            *NewIdNameToRest(&in.Status),
		IsActivated:       in.IsActivated,
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
