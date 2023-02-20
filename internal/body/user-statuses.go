package body

type UserStatus int64

const (
	Registered                UserStatus = 1
	SignedUpForTest           UserStatus = 2
	AttachedScreenshot        UserStatus = 3
	RecommendedFirstCommon    UserStatus = 4
	RecommendedFirstLyceum    UserStatus = 5
	RecommendedFirstOlympiad  UserStatus = 6
	RecommendedFirstMCKO      UserStatus = 7
	RecommendedSecondCommon   UserStatus = 8
	RecommendedSecondLyceum   UserStatus = 9
	RecommendedSecondOlympiad UserStatus = 10
	RecommendedSecondMCKO     UserStatus = 11
	RecommendedBothCommon     UserStatus = 12
	RecommendedBothLyceum     UserStatus = 13
	RecommendedBothOlympiad   UserStatus = 14
	RecommendedBothMCKO       UserStatus = 15
	EnrolledFirstCommon       UserStatus = 16
	EnrolledFirstLyceum       UserStatus = 17
	EnrolledFirstOlympiad     UserStatus = 18
	EnrolledFirstMCKO         UserStatus = 19
	EnrolledSecondCommon      UserStatus = 10
	EnrolledSecondLyceum      UserStatus = 21
	EnrolledSecondOlympiad    UserStatus = 22
	EnrolledSecondMCKO        UserStatus = 23
)

func (u UserStatus) Int64() int64 {
	return int64(u)
}
