package tpportal

import "github.com/golang-jwt/jwt/v4"

type SignUpRequest struct {
	Email             string
	Password          string
	Fio               string
	DateOfBirth       string
	Gender            string
	PhoneNumber       string
	ParentPhoneNumber string
	CurrentSchool     string
	EducationYear     int64
}

type SignInRequest struct {
	Email    string
	Password string
}

type SignInResponse struct {
	User      User
	AuthToken string
}

type Claims struct {
	jwt.RegisteredClaims
	Id int64
}

type CreateUserRequest struct {
	Email             string
	Fio               string
	Password          string
	DateOfBirth       string
	Gender            string
	PhoneNumber       string
	ParentPhoneNumber string
	CurrentSchool     string
	EducationYear     int64
	IsActivated       bool
	Role              string
	StatusId          int64
}

type GetUserResponse struct {
	Id                   int64
	Role                 string
	Fio                  string
	DateOfBirth          string
	Gender               string
	Email                string
	PhoneNumber          string
	ParentPhoneNumber    string
	CurrentSchool        string
	EducationYear        int64
	Status               IdName
	FirstProfile         IdName
	SecondProfile        IdName
	FirstProfileSubject  IdName
	SecondProfileSubject IdName
	ForeignLanguage      IdName
	IsActivated          bool
}

type ListStatusesRequest struct {
	AvailableFor10thClass bool
	AvailableFor9thClass  bool
}

type UpdateUserRequest struct {
	Email             string
	Fio               string
	DateOfBirth       string
	Gender            string
	PhoneNumber       string
	ParentPhoneNumber string
	CurrentSchool     string
	EducationYear     int64
}
