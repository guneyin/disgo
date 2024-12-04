package provider

type MimeType string

const (
	MimeTypeUnknown MimeType = "unknown"
	MimeTypeFolder           = "folder"
	MimeTypeFile             = "file"
)

type User struct {
	Kind         string `json:"kind"`
	DisplayName  string `json:"displayName"`
	PhotoLink    string `json:"photoLink"`
	Me           bool   `json:"me"`
	PermissionId string `json:"permissionId"`
	EmailAddress string `json:"emailAddress"`
}

type File struct {
	Id   string   `json:"id"`
	Name string   `json:"name"`
	Size int64    `json:"size"`
	Type MimeType `json:"type"`
}

type FileList struct {
	Files []File `json:"files"`
}
