package restmodels

type SignUpRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	Fio               string `json:"fio"`
	DateOfBirth       string `json:"date_of_birth"`
	Gender            string `json:"gender"`
	PhoneNumber       string `json:"phone_number"`
	ParentPhoneNumber string `json:"parent_phone_number"`
	CurrentSchool     string `json:"current_school"`
	EducationYear     int64  `json:"education_year"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Id                int64  `json:"id"`
	Email             string `json:"email"`
	Fio               string `json:"fio"`
	DateOfBirth       string `json:"date_of_birth"`
	Gender            string `json:"gender"`
	PhoneNumber       string `json:"phone_number"`
	ParentPhoneNumber string `json:"parent_phone_number"`
	CurrentSchool     string `json:"current_school"`
	EducationYear     int64  `json:"education_year"`
	IsActivated       bool   `json:"is_activated"`
	Role              string `json:"role"`
	Status            IdName `json:"status"`
	AuthToken         string `json:"auth_token"`
}

type CreateUserRequest struct {
	Email             string `json:"email"`
	Fio               string `json:"fio"`
	Password          string `json:"password"`
	DateOfBirth       string `json:"date_of_birth"`
	Gender            string `json:"gender"`
	PhoneNumber       string `json:"phone_number"`
	ParentPhoneNumber string `json:"parent_phone_number"`
	CurrentSchool     string `json:"current_school"`
	EducationYear     int64  `json:"education_year"`
	IsActivated       bool   `json:"is_activated"`
	Role              string `json:"role"`
	StatusId          int64  `json:"status_id"`
}

type ConfirmRecoverRequest struct {
	Password string `json:"password"`
}
