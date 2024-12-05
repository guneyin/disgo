package disgo

import (
	"context"
	"github.com/guneyin/disgo/provider"
	"golang.org/x/oauth2"
	"io"
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

func NewGoogle(ctx context.Context, config provider.GoogleConfig, oauth2 *oauth2.Token) (Provider, error) {
	return provider.NewGoogle(ctx, config, oauth2)
}
