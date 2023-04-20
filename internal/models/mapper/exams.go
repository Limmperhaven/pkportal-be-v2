package mapper

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
)

func NewSetGradesRequestFromRest(in *restmodels.SetGradesRequest) *tpportal.SetGradesRequest {
	return &tpportal.SetGradesRequest{
		UserId:               in.UserId,
		TestDateId:           in.TestDateId,
		RussianLanguageGrade: in.RussianLanguageGrade,
		MathGrade:            in.MathGrade,
		ForeignLanguageGrade: in.ForeignLanguageGrade,
		FirstProfileGrade:    in.FirstProfileGrade,
		SecondProfileGrade:   in.SecondProfileGrade,
	}
}
