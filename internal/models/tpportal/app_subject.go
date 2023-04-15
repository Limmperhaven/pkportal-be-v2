package tpportal

type ListSubjectsRequest struct {
	ProfileId int64
}

type SetSubjectsRequest struct {
	FirstSubjectId  int64
	SecondSubjectId int64
}
