package client

import (
	"bytes"
	"context"
	"github.com/Limmperhaven/pkportal-be-v2/internal/config"
	"github.com/Limmperhaven/pkportal-be-v2/internal/models/tpportal"
	"github.com/friendsofgo/errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

type S3Client struct {
	cl         *minio.Client
	bucketName string
	location   string
}

func InitMinio(cfg *config.S3) (*S3Client, error) {
	minioClient, err := minio.New(cfg.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	exists, err := minioClient.BucketExists(context.Background(), cfg.BucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), cfg.BucketName, minio.MakeBucketOptions{Region: cfg.ServerLocation})
		if err != nil {
			return nil, err
		}
	}
	return &S3Client{bucketName: cfg.BucketName, cl: minioClient}, nil
}

func (c *S3Client) UploadFile(ctx context.Context, req tpportal.UploadFileRequest) (string, error) {
	dataReader := bytes.NewReader(req.FileContent)

	info, err := c.cl.PutObject(ctx, c.bucketName, req.FileKey, dataReader,
		req.FileSize, minio.PutObjectOptions{ContentType: req.ContentType})
	if err != nil {
		return "", err
	}

	return info.Key, nil
}

func (c *S3Client) DownloadFile(ctx context.Context, fileKey string) ([]byte, error) {
	obj, err := c.cl.GetObject(ctx, c.bucketName, fileKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	fileData, err := io.ReadAll(obj)
	if err != nil {
		return nil, errors.New("не получилось прочитать полученный файл")
	}

	return fileData, nil
}

func (c *S3Client) DeleteFile(ctx context.Context, fileKey string) error {
	return c.cl.RemoveObject(ctx, c.bucketName, fileKey, minio.RemoveObjectOptions{ForceDelete: true})
}
