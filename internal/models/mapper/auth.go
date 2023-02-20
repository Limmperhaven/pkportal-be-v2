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

func NewUserToRest(in *tpportal.User) *restmodels.User {
	return &restmodels.User{
		Id:                in.ID,
		Email:             in.Email,
		Fio:               in.Fio,
		DateOfBirth:       in.DateOfBirth.Format("02.01.2006"),
		Gender:            in.Gender.String(),
		PhoneNumber:       in.PhoneNumber,
		ParentPhoneNumber: in.ParentPhoneNumber,
		CurrentSchool:     in.CurrentSchool.String,
		EducationYear:     int64(in.EducationYear),
		Role:              in.Role.String(),
		StatusId:          in.StatusID,
		IsActivated:       in.IsActivated,
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
