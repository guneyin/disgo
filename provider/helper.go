package provider

import (
	"encoding/json"
	"github.com/guneyin/disgo/internal/google"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

func NewOAuth2Token(d []byte) (*oauth2.Token, error) {
	token := &oauth2.Token{}
	err := json.Unmarshal(d, token)
	return token, err
}

func (g *Google) toUserDto(user *drive.User) *User {
	return &User{
		Kind:         user.Kind,
		DisplayName:  user.DisplayName,
		PhotoLink:    user.PhotoLink,
		Me:           user.Me,
		PermissionId: user.PermissionId,
		EmailAddress: user.EmailAddress,
	}
}

func (g *Google) toFile(file *drive.File) *File {
	return &File{
		Id:   file.Id,
		Name: file.Name,
		Size: file.Size,
		Type: g.toMimeType(file.MimeType),
	}
}

func (g *Google) toFileList(fl *drive.FileList) (*FileList, error) {
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
