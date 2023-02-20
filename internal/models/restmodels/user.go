package restmodels

type GetUserResponse struct {
	Id                   int64  `json:"id"`
	Role                 string `json:"role"`
	Fio                  string `json:"fio"`
	DateOfBirth          string `json:"date_of_birth"`
	Gender               string `json:"gender"`
	Email                string `json:"email"`
	PhoneNumber          string `json:"phone_number"`
	ParentPhoneNumber    string `json:"parent_phone_number"`
	CurrentSchool        string `json:"current_school"`
	EducationYear        int64  `json:"education_year"`
	Status               IdName `json:"status"`
	FirstProfile         IdName `json:"first_profile"`
	SecondProfile        IdName `json:"second_profile"`
	FirstProfileSubject  IdName `json:"first_profile_subject"`
	SecondProfileSubject IdName `json:"second_profile_subject"`
	ForeignLanguage      IdName `json:"foreign_language"`
	IsActivated          bool   `json:"is_activated"`
}

type ListStatusesRequest struct {
	AvailableFor10thClass bool `json:"available_for_10_th_class"`
	AvailableFor9thClass  bool `json:"available_for_9_th_class"`
}

type UpdateUserRequest struct {
	Email             string `json:"email"`
	Fio               string `json:"fio"`
	DateOfBirth       string `json:"date_of_birth"`
	Gender            string `json:"gender"`
	PhoneNumber       string `json:"phone_number"`
	ParentPhoneNumber string `json:"parent_phone_number"`
	CurrentSchool     string `json:"current_school"`
	EducationYear     int64  `json:"education_year"`
}
