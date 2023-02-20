package restmodels

type ListProfilesResponseItem struct {
	Id            int64    `json:"id"`
	Name          string   `json:"name"`
	EducationYear int64    `json:"education_year"`
	Subjects      []IdName `json:"subjects"`
}

type SetProfilesToUserRequest struct {
	FirstProfileId  int64 `json:"first_profile_id"`
	SecondProfileId int64 `json:"second_profile_id"`
}
