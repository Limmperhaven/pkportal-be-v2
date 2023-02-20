package tpportal

type CreateTestDateRequest struct {
	Date          string
	Time          string
	Location      string
	MaxPersons    int64
	EducationYear int64
	PubStatus     string
}

type ListTestDatesResponseItem struct {
	Id                int64
	Date              string
	Time              string
	Location          string
	RegisteredPersons int64
	MaxPersons        int64
	EducationYear     int64
	PubStatus         string
}
