package provider

import (
	"context"
	"github.com/guneyin/disgo/internal/google"
	"golang.org/x/oauth2"
	"io"
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

func NewGoogle(ctx context.Context, config GoogleConfig, oauth2 *oauth2.Token) (*Google, error) {
	api, err := google.NewApi(ctx, google.ApiConfig{
		ApiKey:       config.ApiKey,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		CallBackUrl:  config.CallBackUrl,
	}, oauth2)
	if err != nil {
		return nil, err
	}

	return &Google{
		cfg: config,
		api: api,
	}, nil
}

func (g *Google) InitAuth() string {
	return g.api.InitAuth(false)
}

func (g *Google) VerifyAuth(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.api.VerifyAuth(ctx, code)
}

func (g *Google) GetAuthorizedUser() (*User, error) {
	data, err := g.api.About()
	if err != nil {
		return nil, err
	}

	return g.toUserDto(data), nil
}

func (g *Google) GetDirectoryList(parentId string) (*FileList, error) {
	data, err := g.api.FileList(google.MimeTypeFolder, parentId)
	if err != nil {
		return nil, err
	}

	return g.toFileList(data)
}

func (g *Google) GetDirectory(id string) (*FileList, error) {
	data, err := g.api.FileList(google.MimeTypeNone, id)
	if err != nil {
		return nil, err
	}

	return g.toFileList(data)
}

func (g *Google) CreateDirectory(name, parentId string) (*File, error) {
	data, err := g.api.CreateDirectory(name, parentId)
	if err != nil {
		return nil, err
	}

	return g.toFile(data), nil
}

func (g *Google) DeleteDirectory(id string) error {
	return g.api.DeleteDirectory(id)
}

func (g *Google) GetFileMeta(id string) (*File, error) {
	data, err := g.api.GetFileMeta(id)
	if err != nil {
		return nil, err
	}

	return g.toFile(data), nil
}

func (g *Google) DownloadFile(id string, w io.Writer) error {
	return g.api.DownloadFile(id, w)
}

func (g *Google) UploadFile(name, parentId string, media io.Reader) (*File, error) {
	data, err := g.api.UploadFile(name, parentId, media)
	if err != nil {
		return nil, err
	}

	return g.toFile(data), nil
}
