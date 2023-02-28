package mapper

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
)

func NewCreateTestDateFromRest(in *restmodels.CreateTestDateRequest) *tpportal.CreateTestDateRequest {
	return &tpportal.CreateTestDateRequest{
		Date:          in.Date,
		Time:          in.Time,
		Location:      in.Location,
		MaxPersons:    in.MaxPersons,
		EducationYear: in.EducationYear,
		PubStatus:     in.PubStatus,
	}
}

func NewListTestDatesRequestFromRest(in *restmodels.ListTestDatesRequest) *tpportal.ListTestDatesRequest {
	return &tpportal.ListTestDatesRequest{
		EducationYear: in.EducationYear,
	}
}

func NewTestDateResponseToRest(in *tpportal.TestDateResponse) *restmodels.TestDateResponse {
	return &restmodels.TestDateResponse{
		Id:                in.Id,
		Date:              in.Date,
		Time:              in.Time,
		Location:          in.Location,
		AttendedPersons:   in.AttendedPersons,
		RegisteredPersons: in.RegisteredPersons,
		MaxPersons:        in.MaxPersons,
		EducationYear:     in.EducationYear,
		PubStatus:         in.PubStatus,
	}
}

func NewTestDateResponseArrayToRest(in []tpportal.TestDateResponse) []restmodels.TestDateResponse {
	res := make([]restmodels.TestDateResponse, len(in))
	for i, item := range in {
		res[i] = *NewTestDateResponseToRest(&item)
	}
	return res
}
