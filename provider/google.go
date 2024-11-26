package provider

import (
	"context"
	"encoding/json"
	"github.com/guneyin/disgo/internal/google"
	"golang.org/x/oauth2"
)

type Google struct {
	cfg GoogleConfig
	api *google.Api
}

type GoogleConfig struct {
	ApiKey       string `desc:"Google API key"`
	ClientID     string `desc:"OAuth2 2.0 Client ID"`
	ClientSecret string `desc:"OAuth2 2.0 Client Secret"`
	CallBackUrl  string `desc:"Callback URL"`
}

func NewGoogle(config GoogleConfig, accessToken []byte) *Google {
	return &Google{
		cfg: config,
		api: google.NewApi(google.ApiConfig{
			ApiKey:       config.ApiKey,
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			CallBackUrl:  config.CallBackUrl,
		}, string(accessToken)),
	}
}

func (g *Google) InitAuth() string {
	return g.api.InitAuth()
}

func (g *Google) VerifyAuth(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.api.VerifyAuth(ctx, code)
}

func (g *Google) GetAuthorizedUser() ([]byte, error) {
	data, err := g.api.About()
	if err != nil {
		return nil, err
	}

	res, _ := json.MarshalIndent(data, "", "  ")
	return res, nil
}

func (g *Google) GetDirectoryList(parentId string) ([]byte, error) {
	data, err := g.api.FileList(google.MimeTypeFolder, parentId)
	if err != nil {
		return nil, err
	}

	res, _ := json.MarshalIndent(data, "", "  ")
	return res, nil
}

func (g *Google) CreateDirectory(name, parentId string) ([]byte, error) {
	data, err := g.api.CreateDirectory(name, parentId)
	if err != nil {
		return nil, err
	}

	res, _ := json.MarshalIndent(data, "", "  ")
	return res, nil
}

func (g *Google) DeleteDirectory(id string) error {
	return g.api.DeleteDirectory(id)
}
