package restmodels

type SetGradesRequest struct {
	UserId               int64     `json:"user_id"`
	TestDateId           int64     `json:"test_date_id"`
	RussianLanguageGrade NullInt64 `json:"russian_language_grade"`
	MathGrade            NullInt64 `json:"math_grade"`
	ForeignLanguageGrade NullInt64 `json:"foreign_language_grade"`
	FirstProfileGrade    NullInt64 `json:"first_profile_grade"`
	SecondProfileGrade   NullInt64 `json:"second_profile_grade"`
}
