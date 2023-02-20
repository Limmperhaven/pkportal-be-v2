package mapper

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
)

func NewIdNameToRest(in *tpportal.IdName) *restmodels.IdName {
	return &restmodels.IdName{
		Id:   in.Id,
		Name: in.Name,
	}
}

func NewIdNameArrayToRest(in []tpportal.IdName) []restmodels.IdName {
	res := make([]restmodels.IdName, len(in))
	for i, item := range in {
		res[i] = *NewIdNameToRest(&item)
	}
	return res
}
