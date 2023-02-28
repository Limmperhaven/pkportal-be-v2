package tpportal

type CreateTestDateRequest struct {
	Date          string
	Time          string
	Location      string
	MaxPersons    int64
	EducationYear int64
	PubStatus     string
}

type ListTestDatesRequest struct {
	EducationYear int64
}

type TestDateResponse struct {
	Id                int64
	Date              string
	Time              string
	Location          string
	RegisteredPersons int64
	AttendedPersons   int64
	MaxPersons        int64
	EducationYear     int64
	PubStatus         string
}
