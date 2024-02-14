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

func NewDownloadFileResponseToRest(in *tpportal.DownloadFileResponse) *restmodels.DownloadFileResponse {
	return &restmodels.DownloadFileResponse{
		FileName:    in.FileName,
		FileContent: in.FileContent,
		ContentType: in.ContentType,
	}
}

func NewNullInt64ToRest(in *tpportal.NullInt64) *restmodels.NullInt64 {
	return &restmodels.NullInt64{
		Val:     in.Val,
		IsValid: in.IsValid,
	}
}

func NewNullInt64FromRest(in *restmodels.NullInt64) *tpportal.NullInt64 {
	return &tpportal.NullInt64{
		Val:     in.Val,
		IsValid: in.IsValid,
	}
}
