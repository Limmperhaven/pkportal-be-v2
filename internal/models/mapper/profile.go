package mapper

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
)

func NewListProfilesResponseItemToRest(in *tpportal.ListProfilesResponseItem) *restmodels.ListProfilesResponseItem {
	return &restmodels.ListProfilesResponseItem{
		Id:            in.Id,
		Name:          in.Name,
		EducationYear: in.EducationYear,
		Subjects:      NewIdNameArrayToRest(in.Subjects),
	}
}

func NewListProfilesResponseToRest(in []tpportal.ListProfilesResponseItem) []restmodels.ListProfilesResponseItem {
	out := make([]restmodels.ListProfilesResponseItem, len(in))
	for i, item := range in {
		out[i] = *NewListProfilesResponseItemToRest(&item)
	}
	return out
}

func NewSetProfileToUserRequestFromRest(in *restmodels.SetProfilesToUserRequest) *tpportal.SetProfilesToUserRequest {
	return &tpportal.SetProfilesToUserRequest{
		FirstProfileId:  in.FirstProfileId,
		SecondProfileId: in.SecondProfileId,
	}
}
