package tpportal

type ListProfilesResponseItem struct {
	Id            int64
	Name          string
	EducationYear int64
	Subjects      []IdName
}

type SetProfilesToUserRequest struct {
	FirstProfileId  int64
	SecondProfileId int64
}
