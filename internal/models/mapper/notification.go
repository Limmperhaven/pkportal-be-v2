package mapper

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/restmodels"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
)

func NewNotificationFromRest(in *restmodels.Notification) *tpportal.Notification {
	return &tpportal.Notification{
		UserIds: in.UserIds,
		Topic:   in.Topic,
		Message: in.Message,
	}
}
