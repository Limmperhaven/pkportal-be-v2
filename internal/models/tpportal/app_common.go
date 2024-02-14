package tpportal

type IdName struct {
	Id   int64
	Name string
}

type DownloadFileResponse struct {
	FileName    string
	FileContent string
	ContentType string
}

type NullInt64 struct {
	Val     int64
	IsValid bool
}
