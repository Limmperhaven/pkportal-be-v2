package restmodels

type IdName struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type DownloadFileResponse struct {
	FileName    string `json:"file_name"`
	FileContent string `json:"file_content"`
	ContentType string `json:"content_type"`
}

type NullInt64 struct {
	Val     int64 `json:"val"`
	IsValid bool  `json:"is_valid"`
}
