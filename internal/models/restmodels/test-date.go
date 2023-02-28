package restmodels

type CreateTestDateRequest struct {
	Date          string `json:"date"`
	Time          string `json:"time"`
	Location      string `json:"location"`
	MaxPersons    int64  `json:"max_persons"`
	EducationYear int64  `json:"education_year"`
	PubStatus     string `json:"pub_status"`
}

type ListTestDatesRequest struct {
	EducationYear int64 `json:"education_year"`
}

type TestDateResponse struct {
	Id                int64  `json:"id"`
	Date              string `json:"date"`
	Time              string `json:"time"`
	Location          string `json:"location"`
	RegisteredPersons int64  `json:"registered_persons"`
	AttendedPersons   int64  `json:"attended_persons"`
	MaxPersons        int64  `json:"max_persons"`
	EducationYear     int64  `json:"education_year"`
	PubStatus         string `json:"pub_status"`
}
