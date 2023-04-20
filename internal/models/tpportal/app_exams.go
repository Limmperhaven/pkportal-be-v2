package tpportal

type SetGradesRequest struct {
	UserId               int64
	TestDateId           int64
	RussianLanguageGrade int64
	MathGrade            int64
	ForeignLanguageGrade int64
	FirstProfileGrade    int64
	SecondProfileGrade   int64
}
