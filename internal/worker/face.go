package worker

import (
	"context"
	"github.com/Limmperhaven/pkportal-be-v2/internal/domain"
)

type Worker func(ctx context.Context, mail domain.MailClient, s3 domain.S3Client)
