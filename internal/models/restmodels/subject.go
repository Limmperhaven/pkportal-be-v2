package restmodels

type ListSubjectsRequest struct {
	ProfileId int64 `json:"profile_id"`
}

type SetSubjectsRequest struct {
	FirstSubjectId  int64 `json:"first_subject_id"`
	SecondSubjectId int64 `json:"second_subject_id"`
}
