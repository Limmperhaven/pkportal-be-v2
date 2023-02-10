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
