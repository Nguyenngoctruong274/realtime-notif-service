package mail

import (
	"context"
	"os"
	"yes4all/ads-noti-api/pkg/xhttp"

	"go.elastic.co/apm"
)

type IemailAlert interface {
	SendEmailAlert(ctx context.Context, req WarningEmail) (resp interface{}, err error)
}

type emailAlert struct {
	client xhttp.Client
	url    string
}

func NewClient(client xhttp.Client) IemailAlert {
	return &emailAlert{
		client: client,
		url:    os.Getenv("PORTAL_API_URL"),
	}
}

type WarningEmail struct {
	TicketID              int     `json:"ticket_id" form:"ticket_id" uri:"ticket_id"`
	TeamName              string  `json:"team_name" form:"team_name"`
	BudgetSpendPercentage float64 `json:"budget_spend_percentage" form:"budget_spend_percentage"`
	BudgetSpend           float64 `json:"budget_spend" form:"budget_spend"`
	WarningForDate        int64   `json:"warning_for_date" form:"warning_for_date"`
}

func (c *emailAlert) SendEmailAlert(
	ctx context.Context,
	req WarningEmail,
) (resp interface{}, err error) {
	path := c.url + "/api/internal/portal/budget/warning-email"
	xopt := xhttp.RequestOption{
		GroupPath: "/api/internal/portal/budget/warning-email",
	}
	if _, err = c.client.PostJSON(ctx, path, &req, &resp, xopt); err != nil {
		apm.CaptureError(context.Background(), err).Send()
		return
	}
	return
}
