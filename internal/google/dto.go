package google

type User struct {
	User struct {
		Kind         string `json:"kind"`
		DisplayName  string `json:"displayName"`
		PhotoLink    string `json:"photoLink"`
		Me           bool   `json:"me"`
		PermissionId string `json:"permissionId"`
		EmailAddress string `json:"emailAddress"`
	} `json:"user"`
}

type File struct {
	Kind     string `json:"kind"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
}

type FileList struct {
	Kind             string `json:"kind"`
	IncompleteSearch bool   `json:"incompleteSearch"`
	Files            []File `json:"files"`
}
