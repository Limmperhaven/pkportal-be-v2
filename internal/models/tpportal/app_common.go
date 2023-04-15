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
