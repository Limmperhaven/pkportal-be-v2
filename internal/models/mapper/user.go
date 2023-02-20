package mapper

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
)

func NewGetUserResponseToRest(in *tpportal.GetUserResponse) *restmodels.GetUserResponse {
	return &restmodels.GetUserResponse{
		Id:                   in.Id,
		Role:                 in.Role,
		Fio:                  in.Fio,
		DateOfBirth:          in.DateOfBirth,
		Gender:               in.Gender,
		Email:                in.Email,
		PhoneNumber:          in.PhoneNumber,
		ParentPhoneNumber:    in.ParentPhoneNumber,
		CurrentSchool:        in.CurrentSchool,
		EducationYear:        in.EducationYear,
		Status:               *NewIdNameToRest(&in.Status),
		FirstProfile:         *NewIdNameToRest(&in.FirstProfile),
		SecondProfile:        *NewIdNameToRest(&in.SecondProfile),
		FirstProfileSubject:  *NewIdNameToRest(&in.FirstProfileSubject),
		SecondProfileSubject: *NewIdNameToRest(&in.SecondProfileSubject),
		ForeignLanguage:      *NewIdNameToRest(&in.ForeignLanguage),
		IsActivated:          in.IsActivated,
	}
}

func NewListStatusesRequestFromRest(in *restmodels.ListStatusesRequest) *tpportal.ListStatusesRequest {
	return &tpportal.ListStatusesRequest{
		AvailableFor10thClass: in.AvailableFor10thClass,
		AvailableFor9thClass:  in.AvailableFor9thClass,
	}
}

func NewUpdateUserRequestFromRest(in *restmodels.UpdateUserRequest) *tpportal.UpdateUserRequest {
	return &tpportal.UpdateUserRequest{
		Email:             in.Email,
		Fio:               in.Fio,
		DateOfBirth:       in.DateOfBirth,
		Gender:            in.Gender,
		PhoneNumber:       in.PhoneNumber,
		ParentPhoneNumber: in.ParentPhoneNumber,
		CurrentSchool:     in.CurrentSchool,
		EducationYear:     in.EducationYear,
	}
}
