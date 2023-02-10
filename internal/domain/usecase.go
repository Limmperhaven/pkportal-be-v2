package domain

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage"
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage/stpg"
)

type Usecase struct {
	st storage.PGer
}

func NewUsecase() *Usecase {
	return &Usecase{st: stpg.Gist()}
}
