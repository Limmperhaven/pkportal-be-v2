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
	Id                int64
	Email             string
	Fio               string
	DateOfBirth       string
	Gender            string
	PhoneNumber       string
	ParentPhoneNumber string
	CurrentSchool     string
	EducationYear     int64
	IsActivated       bool
	Role              string
	Status            IdName
	AuthToken         string
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
	TestDate             GetUserResponseTestDate
	IsActivated          bool
}

type GetUserResponseTestDate struct {
	TestDateId int64
	IsAttended bool
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

type UploadScreenshotRequest struct {
	FileName    string
	FileSize    int64
	FileContent []byte
}

type UploadFileRequest struct {
	FileKey     string
	FileSize    int64
	FileContent []byte
	ContentType string
}

type DownloadScreenshotResponse struct {
	FileName    string
	FileContent []byte
	ContentType string
}

type ListUsersRequest struct {
	ProfileIds     []int64
	EducationYears []int64
	StatusIds      []int64
	TestDateIds    []int64
}
