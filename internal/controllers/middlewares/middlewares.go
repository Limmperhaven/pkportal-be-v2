package middlewares

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage"
	"github.com/Limmperhaven/pkportal-be-v2/internal/storage/stpg"
)

type MiddlewareStorage struct {
	st storage.PGer
}

func NewMiddlewareStorage() *MiddlewareStorage {
	return &MiddlewareStorage{st: stpg.Gist()}
}
