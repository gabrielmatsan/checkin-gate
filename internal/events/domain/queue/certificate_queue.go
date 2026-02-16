package queue

import (
	"context"
	"time"
)

// CertificateQueue define a interface (Port) para a fila de certificados
type CertificateQueue interface {
	// Enqueue adiciona um job na fila
	Enqueue(ctx context.Context, job *CertificateJob) error

	// EnqueueBatch adiciona múltiplos jobs na fila (mais eficiente)
	EnqueueBatch(ctx context.Context, jobs []*CertificateJob) error

	// Dequeue remove e retorna o próximo job da fila (blocking)
	Dequeue(ctx context.Context) (*CertificateJob, error)

	// DequeueWithTimeout remove com timeout (retorna nil se timeout)
	DequeueWithTimeout(ctx context.Context, timeout time.Duration) (*CertificateJob, error)

	// Len retorna o tamanho atual da fila
	Len(ctx context.Context) (int64, error)
}
