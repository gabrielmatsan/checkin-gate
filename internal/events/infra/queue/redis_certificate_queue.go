package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domainqueue "github.com/gabrielmatsan/checkin-gate/internal/events/domain/queue"
	"github.com/redis/go-redis/v9"
)

const certificateQueueKey = "certificate:jobs"

type RedisCertificateQueue struct {
	client *redis.Client
}

func NewRedisCertificateQueue(client *redis.Client) *RedisCertificateQueue {
	return &RedisCertificateQueue{
		client: client,
	}
}

// Enqueue adiciona um job na fila usando LPUSH
func (q *RedisCertificateQueue) Enqueue(ctx context.Context, job *domainqueue.CertificateJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	if err := q.client.LPush(ctx, certificateQueueKey, data).Err(); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	return nil
}

// EnqueueBatch adiciona múltiplos jobs na fila usando pipeline
func (q *RedisCertificateQueue) EnqueueBatch(ctx context.Context, jobs []*domainqueue.CertificateJob) error {
	if len(jobs) == 0 {
		return nil
	}

	pipe := q.client.Pipeline()

	for _, job := range jobs {
		data, err := json.Marshal(job)
		if err != nil {
			return fmt.Errorf("failed to marshal job: %w", err)
		}
		pipe.LPush(ctx, certificateQueueKey, data)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute batch enqueue: %w", err)
	}

	return nil
}

// Dequeue remove e retorna o próximo job da fila usando BRPOP (blocking)
func (q *RedisCertificateQueue) Dequeue(ctx context.Context) (*domainqueue.CertificateJob, error) {
	return q.DequeueWithTimeout(ctx, 0) // 0 = block indefinitely
}

// DequeueWithTimeout remove com timeout usando BRPOP
func (q *RedisCertificateQueue) DequeueWithTimeout(ctx context.Context, timeout time.Duration) (*domainqueue.CertificateJob, error) {
	result, err := q.client.BRPop(ctx, timeout, certificateQueueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // timeout, no job available
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	// BRPop returns [key, value]
	if len(result) < 2 {
		return nil, fmt.Errorf("unexpected BRPop result: %v", result)
	}

	var job domainqueue.CertificateJob
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

// Len retorna o tamanho atual da fila
func (q *RedisCertificateQueue) Len(ctx context.Context) (int64, error) {
	length, err := q.client.LLen(ctx, certificateQueueKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}
	return length, nil
}

// Compile-time check to ensure RedisCertificateQueue implements CertificateQueue
var _ domainqueue.CertificateQueue = (*RedisCertificateQueue)(nil)
