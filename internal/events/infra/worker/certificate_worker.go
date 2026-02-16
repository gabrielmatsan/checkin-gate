package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/queue"
	"github.com/gabrielmatsan/checkin-gate/internal/events/infra/pdf"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/mail"
	"go.uber.org/zap"
)

type CertificateWorker struct {
	queue        queue.CertificateQueue
	generator    pdf.CertificateGenerator
	emailService mail.EmailService
	logger       *zap.Logger
}

func NewCertificateWorker(
	queue queue.CertificateQueue,
	generator pdf.CertificateGenerator,
	emailService mail.EmailService,
	logger *zap.Logger,
) *CertificateWorker {
	return &CertificateWorker{
		queue:        queue,
		generator:    generator,
		emailService: emailService,
		logger:       logger,
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

			w.processJob(ctx, job)
		}
	}
}

func (w *CertificateWorker) processJob(ctx context.Context, job *queue.CertificateJob) {
	w.logger.Info("processing certificate job",
		zap.String("job_id", job.GetJobID()),
		zap.String("event", fmt.Sprintf("%s (%s)", job.GetEventInfo().EventName, job.GetEventInfo().EventID)),
		zap.String("user", fmt.Sprintf("%s <%s>", job.GetUserInfo().UserName, job.GetUserInfo().UserEmail)),
		zap.String("activity", fmt.Sprintf("%s (%s)", job.GetActivityInfo().ActivityName, job.GetActivityInfo().ActivityID)),
		zap.Time("checked_at", job.GetCheckedAt()),
		zap.Time("enqueued_at", job.GetEnqueuedAt()),
	)

	// Calcula a carga horária
	workload := w.calculateWorkload(job.GetActivityInfo())

	// Monta os dados do certificado
	data := pdf.CertificateData{
		RecipientName:   job.GetUserInfo().UserName,
		EventName:       job.GetEventInfo().EventName,
		EventDate:       w.formatDate(job.GetActivityInfo().ActivityDate),
		Workload:        workload,
		DirectorName:    "Dr. João Silva",
		CoordinatorName: "Dra. Maria Santos",
		CertificateDate: w.formatDate(time.Now()),
	}

	// Gera o PDF
	pdfBytes, err := w.generator.Generate(ctx, data)
	if err != nil {
		w.logger.Error("failed to generate certificate PDF",
			zap.String("job_id", job.GetJobID()),
			zap.Error(err),
		)
		return
	}

	w.logger.Info("certificate generated",
		zap.String("job_id", job.GetJobID()),
		zap.Int("size_bytes", len(pdfBytes)),
	)

	// Envia email com certificado
	emailParams := mail.SendEmailParams{
		To:      job.GetUserInfo().UserEmail,
		Subject: fmt.Sprintf("Certificado - %s", job.GetEventInfo().EventName),
		Body:    w.buildCertificateEmailBody(job.GetUserInfo().UserName, job.GetEventInfo().EventName),
		Attachments: []mail.Attachment{
			{
				Filename:    fmt.Sprintf("certificado-%s.pdf", job.GetJobID()),
				Content:     pdfBytes,
				ContentType: "application/pdf",
			},
		},
	}

	if err := w.emailService.Send(ctx, emailParams); err != nil {
		w.logger.Error("failed to send certificate email",
			zap.String("job_id", job.GetJobID()),
			zap.String("email", job.GetUserInfo().UserEmail),
			zap.Error(err),
		)
		return
	}

	w.logger.Info("certificate email sent successfully",
		zap.String("job_id", job.GetJobID()),
		zap.String("email", job.GetUserInfo().UserEmail),
	)
}

// calculateWorkload calcula a carga horária a partir do horário de início e fim
func (w *CertificateWorker) calculateWorkload(activity queue.ActivityInfo) string {
	duration := activity.EndTime.Sub(activity.StartTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if minutes == 0 {
		if hours == 1 {
			return "1 hora"
		}
		return fmt.Sprintf("%d horas", hours)
	}

	if hours == 0 {
		if minutes == 1 {
			return "1 minuto"
		}
		return fmt.Sprintf("%d minutos", minutes)
	}

	hoursLabel := "hora"
	if hours > 1 {
		hoursLabel = "horas"
	}

	minutesLabel := "minuto"
	if minutes > 1 {
		minutesLabel = "minutos"
	}

	return fmt.Sprintf("%d %s e %d %s", hours, hoursLabel, minutes, minutesLabel)
}

// formatDate formata uma data no padrão brasileiro
func (w *CertificateWorker) formatDate(t time.Time) string {
	months := []string{
		"", "janeiro", "fevereiro", "março", "abril", "maio", "junho",
		"julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
	}
	return fmt.Sprintf("%d de %s de %d", t.Day(), months[t.Month()], t.Year())
}

// buildCertificateEmailBody constrói o corpo HTML do email de certificado
func (w *CertificateWorker) buildCertificateEmailBody(name, eventName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4F46E5; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }
        .footer { text-align: center; margin-top: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Parabéns!</h1>
        </div>
        <div class="content">
            <p>Olá, <strong>%s</strong>!</p>
            <p>É com grande satisfação que enviamos seu certificado de participação no evento <strong>%s</strong>.</p>
            <p>O certificado está em anexo neste email em formato PDF. Guarde-o em um local seguro!</p>
            <p>Agradecemos sua participação e esperamos vê-lo(a) em nossos próximos eventos.</p>
            <p>Atenciosamente,<br>Equipe Checkin Gate</p>
        </div>
        <div class="footer">
            <p>Este é um email automático, por favor não responda.</p>
        </div>
    </div>
</body>
</html>
`, name, eventName)
}
