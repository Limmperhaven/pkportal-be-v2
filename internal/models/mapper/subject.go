package mapper

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
)

func NewListSubjectsRequestFromRest(in *restmodels.ListSubjectsRequest) *tpportal.ListSubjectsRequest {
	return &tpportal.ListSubjectsRequest{ProfileId: in.ProfileId}
}

func NewSetSubjectsRequestFromRest(in *restmodels.SetSubjectsRequest) *tpportal.SetSubjectsRequest {
	return &tpportal.SetSubjectsRequest{
		FirstSubjectId:  in.FirstSubjectId,
		SecondSubjectId: in.SecondSubjectId,
	}
}
