package auth_cloak //nolint

import (
	"context"
	"fmt"
	"net/http"
	"yes4all/ads-noti-api/pkg/config"
	"yes4all/ads-noti-api/pkg/logger"
	"yes4all/ads-noti-api/pkg/xhttp"
)

type Auth interface {
	GetUserInfo(ctx context.Context, authenToken string, keycloakType string) (err error)
}
type auth struct {
	client    xhttp.Client
	domain    string
	domainYSC string
}

func NewAuth(client xhttp.Client) Auth {
	return &auth{
		client:    client,
		domain:    config.KeycloakConfig().KeycloakDomain,
		domainYSC: config.KeycloakConfig().KeycloakYSCDomain,
	}
}

func (m *auth) GetUserInfo(ctx context.Context, authenToken string, keycloakType string) (err error) {
	domain := m.domain
	if keycloakType != "" {
		domain = m.domainYSC
	}
	url := fmt.Sprintf("%s/protocol/openid-connect/userinfo", domain)
	opts := xhttp.RequestOption{
		Header: map[string]string{"Authorization": "bearer " + authenToken},
	}
	var response interface{}
	entry := logger.NewLogger().WithKeyword(ctx, "keycloak-getuserinfo")
	entry.WithOutputField("url", url).WithOutputField("opts", opts)
	var status int
	status, err = m.client.Get(
		ctx,
		url,
		&response,
		opts,
	)
	entry.WithOutputField("response", response)
	if err != nil {
		entry.WithError(err).Error()
		return err
	}
	if status != http.StatusOK {
		err = fmt.Errorf("statusCode failed")
		entry.WithError(err).Error()
		return err
	}

	return nil
}
