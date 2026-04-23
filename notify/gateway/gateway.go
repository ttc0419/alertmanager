package gateway

import (
	"context"
	"log/slog"

	"github.com/prometheus/alertmanager/alert"
	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/notify/gateway/im"
	"github.com/prometheus/alertmanager/template"
)

type Notifier struct {
	imNotifier notify.Notifier
}

func New(cfg *config.GatewayConfig, _ *template.Template, _ *slog.Logger) *Notifier {
	return &Notifier{
		imNotifier: im.NewLarkNotifier(cfg),
	}
}

func (n *Notifier) Notify(ctx context.Context, alerts ...*alert.Alert) (bool, error) {
	return n.imNotifier.Notify(ctx, alerts...)
}
