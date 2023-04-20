package restmodels

type SetGradesRequest struct {
	UserId               int64 `json:"user_id"`
	TestDateId           int64 `json:"test_date_id"`
	RussianLanguageGrade int64 `json:"russian_language_grade"`
	MathGrade            int64 `json:"math_grade"`
	ForeignLanguageGrade int64 `json:"foreign_language_grade"`
	FirstProfileGrade    int64 `json:"first_profile_grade"`
	SecondProfileGrade   int64 `json:"second_profile_grade"`
}
