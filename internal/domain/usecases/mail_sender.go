package usecases

import (
	"tasks/internal/domain/models"
	"tasks/pkg/infra/logger"

	"context"
)

type MailSender struct {
	l logger.Logger 
}

func NewSender(l logger.Logger) *MailSender {
	return &MailSender{
		l:l,
	}
}

func (m *MailSender) SendRequestMail(ctx context.Context, task *models.Task, email string) string {
	m.l.Sugar().Infof("Mail text: Request mail for task %s. Send to %s", task.Title, email)
	// согласующий направляется по ссылке /tasks/api/v1/tasks/approve или /tasks/api/v1/tasks/decline
	// для примера возвращаем статус строкой
	return "declined"
}

func (m *MailSender) SendDeclineMail(ctx context.Context, task *models.Task, emailReceiver string, emailDecliner string) {
	m.l.Sugar().Infof("Mail text: Task %s declined by %s. Send to %s", task.Title, emailDecliner, emailReceiver)
}

func (m *MailSender) SendApproveMail(ctx context.Context, task *models.Task, email string) {
	m.l.Sugar().Infof("Mail text: Task %s approved by everyone. Send to %s", task.Title, email)
}