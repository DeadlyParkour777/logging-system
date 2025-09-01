package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DeadlyParkour777/logging-system/pkg/models"
)

type Notifier interface {
	SendNotification(ctx context.Context, message string) error
}

type AlertService struct {
	notifier Notifier
}

func NewAlertService(notifier Notifier) *AlertService {
	return &AlertService{notifier: notifier}
}

func (s *AlertService) ProcessAlert(ctx context.Context, entry models.LogEntry) error {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("*🚨 Alert: %s*\n\n", escapeMarkdown(entry.Level)))
	builder.WriteString(fmt.Sprintf("*Service:* `%s`\n", escapeMarkdown(entry.ServiceName)))
	builder.WriteString(fmt.Sprintf("*Time:* `%s`\n", escapeMarkdown(entry.Timestamp.Format(time.RFC3339))))
	builder.WriteString(fmt.Sprintf("*Message:* \n```\n%s\n```\n", escapeMarkdown(entry.Message)))

	if len(entry.Metadata) > 0 {
		builder.WriteString("*Metadata:*\n")
		for k, v := range entry.Metadata {
			builder.WriteString(fmt.Sprintf(" \\- `%s`: `%s`\n", escapeMarkdown(k), escapeMarkdown(v)))
		}
	}

	return s.notifier.SendNotification(ctx, builder.String())
}

func escapeMarkdown(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
		"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>",
		"#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|",
		"\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!",
	)
	return replacer.Replace(s)
}
