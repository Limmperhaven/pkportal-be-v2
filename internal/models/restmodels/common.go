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
