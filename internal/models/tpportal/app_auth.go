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

type UserWithAuth struct {
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
