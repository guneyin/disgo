package provider

import (
	"context"
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

func (g *Google) CreateDirectory(name, parentId string) (*File, error) {
	data, err := g.api.CreateDirectory(name, parentId)
	if err != nil {
		return nil, err
	}

	return g.toFile(data)
}

func (g *Google) DeleteDirectory(id string) error {
	return g.api.DeleteDirectory(id)
}

func (g *Google) toUserDto(user *google.User) *User {
	return &User{
		Kind:         user.User.Kind,
		DisplayName:  user.User.DisplayName,
		PhotoLink:    user.User.PhotoLink,
		Me:           user.User.Me,
		PermissionId: user.User.PermissionId,
		EmailAddress: user.User.EmailAddress,
	}
}

func (g *Google) toFile(file *google.File) (*File, error) {
	return &File{
		Id:   file.Id,
		Name: file.Name,
		Type: g.toMimeType(file.MimeType),
	}, nil
}

func (g *Google) toFileList(fl *google.FileList) (*FileList, error) {
	list := make([]File, len(fl.Files))
	for i, f := range fl.Files {
		list[i] = File{
			Id:   f.Id,
			Name: f.Name,
			Type: g.toMimeType(f.MimeType),
		}
	}

	return &FileList{Files: list}, nil
}

func (g *Google) toMimeType(mt string) MimeType {
	switch google.MimeType(mt) {
	case google.MimeTypeFolder:
		return MimeTypeFolder
	case google.MimeTypeFile:
		return MimeTypeFile
	default:
		return MimeTypeUnknown
	}
}
