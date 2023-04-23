package tpportal

type SetGradesRequest struct {
	UserId               int64
	TestDateId           int64
	RussianLanguageGrade NullInt64
	MathGrade            NullInt64
	ForeignLanguageGrade NullInt64
	FirstProfileGrade    NullInt64
	SecondProfileGrade   NullInt64
}
