package domain

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage"
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage/stpg"
)

type Usecase struct {
	s3   S3Client
	mail MailClient
	st   storage.PGer
}

func NewUsecase(mail MailClient, s3 S3Client) *Usecase {
	return &Usecase{st: stpg.Gist(), mail: mail, s3: s3}
}
