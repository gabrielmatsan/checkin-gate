package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/queue"
	"go.uber.org/zap"
)

type CertificateWorker struct {
	queue  queue.CertificateQueue
	logger *zap.Logger
}

func NewCertificateWorker(queue queue.CertificateQueue, logger *zap.Logger) *CertificateWorker {
	return &CertificateWorker{
		queue:  queue,
		logger: logger,
	}
}

func (w *CertificateWorker) Start(ctx context.Context) error {
	w.logger.Info("certificate worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("certificate worker stopping")
			return nil
		default:
			job, err := w.queue.DequeueWithTimeout(ctx, 5*time.Second)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				w.logger.Error("failed to dequeue job", zap.Error(err))
				continue
			}

			if job == nil {
				continue
			}

			w.processJob(job)
		}
	}
}

func (w *CertificateWorker) processJob(job *queue.CertificateJob) {
	w.logger.Info("FILAAAAA")
	w.logger.Info("processing certificate job",
		zap.String("job_id", job.GetJobID()),
		zap.String("event", fmt.Sprintf("%s (%s)", job.GetEventInfo().EventName, job.GetEventInfo().EventID)),
		zap.String("user", fmt.Sprintf("%s <%s>", job.GetUserInfo().UserName, job.GetUserInfo().UserEmail)),
		zap.String("activity", fmt.Sprintf("%s (%s)", job.GetActivityInfo().ActivityName, job.GetActivityInfo().ActivityID)),
		zap.Time("checked_at", job.GetCheckedAt()),
		zap.Time("enqueued_at", job.GetEnqueuedAt()),
	)
}
