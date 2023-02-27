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

func NewListTestDateResponseItemToRest(in *tpportal.ListTestDatesResponseItem) *restmodels.ListTestDatesResponseItem {
	return &restmodels.ListTestDatesResponseItem{
		Id:                in.Id,
		Date:              in.Date,
		Time:              in.Time,
		Location:          in.Location,
		RegisteredPersons: in.RegisteredPersons,
		MaxPersons:        in.MaxPersons,
		EducationYear:     in.EducationYear,
		PubStatus:         in.PubStatus,
	}
}

func NewListTestDateResponseToRest(in []tpportal.ListTestDatesResponseItem) []restmodels.ListTestDatesResponseItem {
	res := make([]restmodels.ListTestDatesResponseItem, len(in))
	for i, item := range in {
		res[i] = *NewListTestDateResponseItemToRest(&item)
	}
	return res
}
