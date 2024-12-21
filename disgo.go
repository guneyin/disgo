package disgo

import (
	"context"
	"fmt"
	"github.com/guneyin/disgo/provider"
	"golang.org/x/oauth2"
	"io"
)

type ProviderType string

const (
	ProviderTypeGoogle ProviderType = "google"
)

type Provider interface {
	// InitAuth initialize login url
	InitAuth() string
	// VerifyAuth verify login code and generate oauth2 token
	VerifyAuth(ctx context.Context, code string) (*oauth2.Token, error)
	// GetAuthorizedUser get oauth2 user info
	GetAuthorizedUser() (*provider.User, error)
	// GetDirectoryList list directories
	GetDirectoryList(parentId string) (*provider.FileList, error)
	// GetDirectory list directory content
	GetDirectory(id string) (*provider.FileList, error)
	// CreateDirectory create a directory on target
	CreateDirectory(name, parentId string) (*provider.File, error)
	// DeleteDirectory delete directory by id
	DeleteDirectory(id string) error
	// GetFileMeta get file meta
	GetFileMeta(id string) (*provider.File, error)
	// DownloadFile download file content
	DownloadFile(id string, w io.Writer) error
	// UploadFile upload file
	UploadFile(name, parentId string, media io.Reader) (*provider.File, error)
}

var _ Provider = (*provider.Google)(nil)

func NewProviderType(pt string) (ProviderType, error) {
	switch ProviderType(pt) {
	case ProviderTypeGoogle:
		return ProviderType(pt), nil
	}
	return "", fmt.Errorf("unknown provider type: %s", pt)
}

func New(ctx context.Context, pt ProviderType, config, token []byte) (Provider, error) {
	switch pt {
	case ProviderTypeGoogle:
		return NewGoogleDrive(ctx, config, token)
	default:
		return nil, fmt.Errorf("unknown provider: %s", pt)
	}
}

func NewGoogleDrive(ctx context.Context, config, token []byte) (Provider, error) {
	cfg, err := provider.NewGoogleConfig(config)
	if err != nil {
		return nil, err
	}
	tkn, _ := provider.NewOAuth2Token(token)

	return provider.NewGoogleDrive(ctx, cfg, tkn)
}
