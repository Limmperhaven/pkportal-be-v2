package worker

import (
	"context"
	"github.com/Limmperhaven/pkportal-be-v2/internal/domain"
)

type Pool struct {
	s3        domain.S3Client
	mail      domain.MailClient
	ctx       context.Context
	cancel    context.CancelFunc
	errChan   chan error
	workers   []Worker
	isStarted bool
}

func NewPool(mail domain.MailClient, s3 domain.S3Client) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	return &Pool{
		s3:        s3,
		mail:      mail,
		ctx:       ctx,
		cancel:    cancel,
		errChan:   errChan,
		isStarted: false,
	}
}

func (p *Pool) Start() {
	for i := range p.workers {
		go p.workers[i](p.ctx, p.mail, p.s3)
	}
	p.isStarted = true
}

func (p *Pool) Stop() {
	p.cancel()
	p.isStarted = false
}

func (p *Pool) AddWorker(worker Worker) {
	p.workers = append(p.workers, worker)
	if p.isStarted {
		go worker(p.ctx, p.mail, p.s3)
	}
}
