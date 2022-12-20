package ports

import (
	"tasks/internal/domain/models"
	"context"
)

type MailSender interface {
	SendRequestMail(ctx context.Context, task *models.Task, email string)
	SendDeclineMail(ctx context.Context, task *models.Task, email string)
	SendApproveMail(ctx context.Context, task *models.Task, email string)
}