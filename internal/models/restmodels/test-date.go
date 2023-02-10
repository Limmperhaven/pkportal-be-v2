package restmodels

type CreateTestDateRequest struct {
	Date          string `json:"date"`
	Time          string `json:"time"`
	Location      string `json:"location"`
	MaxPersons    int64  `json:"max_persons"`
	EducationYear int64  `json:"education_year"`
	PubStatus     string `json:"pub_status"`
}
