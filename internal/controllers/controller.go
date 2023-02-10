package controllers

import "github.com/Limmperhaven/pkportal-be-v2/internal/domain"

type ControllerStorage struct {
	uc *domain.Usecase
}

func NewController(uc *domain.Usecase) *ControllerStorage {
	return &ControllerStorage{uc: uc}
}
