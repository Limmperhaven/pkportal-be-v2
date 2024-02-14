package mapper

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
)

func NewSetGradesRequestFromRest(in *restmodels.SetGradesRequest) *tpportal.SetGradesRequest {
	return &tpportal.SetGradesRequest{
		UserId:               in.UserId,
		TestDateId:           in.TestDateId,
		RussianLanguageGrade: *NewNullInt64FromRest(&in.RussianLanguageGrade),
		MathGrade:            *NewNullInt64FromRest(&in.MathGrade),
		ForeignLanguageGrade: *NewNullInt64FromRest(&in.ForeignLanguageGrade),
		FirstProfileGrade:    *NewNullInt64FromRest(&in.FirstProfileGrade),
		SecondProfileGrade:   *NewNullInt64FromRest(&in.SecondProfileGrade),
	}
}
