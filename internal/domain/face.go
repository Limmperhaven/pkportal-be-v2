package domain

import (
	"context"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
)

type MailClient interface {
	SendTextEmail(subject, message string, to []string) error
	SendHTMLEmail(subject, message string, to []string) error
}

type S3Client interface {
	UploadFile(ctx context.Context, req tpportal.UploadFileRequest) (string, error)
	DownloadFile(ctx context.Context, fileKey string) ([]byte, error)
	DeleteFile(ctx context.Context, fileKey string) error
}
