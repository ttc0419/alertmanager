package im

import (
	"context"
	"encoding/json"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/prometheus/alertmanager/alert"
	"github.com/prometheus/alertmanager/config"
)

type LarkNotifier struct {
	larkClient *lark.Client
}

func NewLarkNotifier(cfg *config.GatewayConfig) *LarkNotifier {
	return &LarkNotifier{
		larkClient: lark.NewClient(cfg.AppID, cfg.AppSecret),
	}
}

type templateCardDataMember struct {
	ID       string `json:"id,omitempty"`
	Selected bool   `json:"selected,omitempty"`
}

type templateCardData struct {
	TemplateID       string         `json:"template_id,omitempty"`
	TemplateVariable map[string]any `json:"template_variable,omitempty"`
}

type templateCard struct {
	Type string           `json:"type,omitempty"`
	Data templateCardData `json:"data,omitempty"`
}

func (n *LarkNotifier) Notify(ctx context.Context, alerts ...*alert.Alert) (bool, error) {
	color := "red"
	switch alerts[0].Labels["priority"] {
	case "2":
		color = "orange"
	case "3":
		color = "yellow"
	}

	priority := "P" + alerts[0].Labels["priority"]
	if alerts[0].Resolved() {
		color = "green"
		priority = "Resolved"
	}

	c := templateCard{
		Type: "template",
		Data: templateCardData{
			TemplateID: "AAqe9aD0Zq0iz",
			TemplateVariable: map[string]any{
				"color":    color,
				"priority": priority,
				"summary":  alerts[0].Annotations["summary"],
				"service":  alerts[0].Labels["container"],
				"start_at": alerts[0].StartsAt.Format("2006-01-02 15:04"),
				"members": []templateCardDataMember{{
					ID: "ou_2a32693bd91108ce7bd13984be9b9def",
				}},
			},
		},
	}

	cb, err := json.Marshal(c)
	if err != nil {
		return false, err
	}

	_, err = n.larkClient.Im.V1.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).Body(
		larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId("ou_2a32693bd91108ce7bd13984be9b9def").
			MsgType(larkim.MsgTypeInteractive).
			Content(string(cb)).Build(),
	).Build())
	return false, err
}
