package disgo

import (
	"context"
	"github.com/guneyin/disgo/provider"
	"golang.org/x/oauth2"
)

type Provider interface {
	// InitAuth initialize login url
	InitAuth() string
	// VerifyAuth verify login code and generate oauth2 token
	VerifyAuth(ctx context.Context, code string) (*oauth2.Token, error)
	// GetAuthorizedUser get oauth2 user info
	GetAuthorizedUser() ([]byte, error)
	// GetDirectoryList list directories
	GetDirectoryList(parentId string) ([]byte, error)
	// CreateDirectory create a directory on target
	CreateDirectory(name, parentId string) ([]byte, error)
	// DeleteDirectory delete directory by id
	DeleteDirectory(id string) error
}

var _ Provider = (*provider.Google)(nil)

func NewGoogle(config provider.GoogleConfig, accessToken []byte) (Provider, error) {
	gp := provider.NewGoogle(config, accessToken)
	return gp, nil
}
