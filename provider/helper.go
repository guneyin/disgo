package provider

import (
	"github.com/guneyin/disgo/internal/google"
	"strconv"
)

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
	size, _ := strconv.ParseInt(file.Size, 10, 64)
	return &File{
		Id:   file.Id,
		Name: file.Name,
		Size: size,
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
