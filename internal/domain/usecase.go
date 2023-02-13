package domain

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/client"
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage"
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage/stpg"
)

type Usecase struct {
	mail *client.MailClient
	st   storage.PGer
}

func NewUsecase(mail *client.MailClient) *Usecase {
	return &Usecase{st: stpg.Gist(), mail: mail}
}
