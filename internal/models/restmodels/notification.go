package restmodels

type Notification struct {
	UserIds []int64 `json:"user_ids"`
	Topic   string  `json:"topic"`
	Message string  `json:"message"`
}
