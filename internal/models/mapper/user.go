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
		ShortFio:             in.ShortFIO,
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
		TestDate:             *NewGetUserResponseTestDateToRest(&in.TestDate),
		IsActivated:          in.IsActivated,
	}
}

func NewGetUserResponseArrayToRest(in []tpportal.GetUserResponse) []restmodels.GetUserResponse {
	res := make([]restmodels.GetUserResponse, len(in))
	for i, item := range in {
		res[i] = *NewGetUserResponseToRest(&item)
	}
	return res
}

func NewGetUserResponseTestDateToRest(in *tpportal.GetUserResponseTestDate) *restmodels.GetUserResponseTestDate {
	return &restmodels.GetUserResponseTestDate{
		Id:            in.Id,
		Date:          in.Date,
		Time:          in.Time,
		Location:      in.Location,
		MaxPersons:    in.MaxPersons,
		EducationYear: in.EducationYear,
		PubStatus:     in.PubStatus,
		IsAttended:    in.IsAttended,
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

func NewUploadScreenshotRequestFromRest(fileName string, fileSize int64, fileContent []byte) *tpportal.UploadScreenshotRequest {
	return &tpportal.UploadScreenshotRequest{
		FileName:    fileName,
		FileSize:    fileSize,
		FileContent: fileContent,
	}
}

func NewListUsersRequestFromRest(in *restmodels.UserFilter) *tpportal.UserFilter {
	return &tpportal.UserFilter{
		ProfileIds:     in.ProfileIds,
		EducationYears: in.EducationYears,
		StatusIds:      in.StatusIds,
		TestDateIds:    in.TestDateIds,
	}
}
