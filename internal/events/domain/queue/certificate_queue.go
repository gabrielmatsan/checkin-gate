package queue

import (
	"context"
	"time"
)

type CertificateQueue struct {
	Enqueue            func(ctx context.Context, job *CertificateJob) error
	EnqueueBatch       func(ctx context.Context, jobs []*CertificateJob) error
	Dequeue            func(ctx context.Context) (*CertificateJob, error)
	DequeueBatch       func(ctx context.Context, count int) ([]*CertificateJob, error)
	DequeueWithTimeout func(ctx context.Context, timeout time.Duration) (*CertificateJob, error)
	Len                func(ctx context.Context) (int, error)
}
