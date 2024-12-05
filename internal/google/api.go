package google

import (
	"context"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type MimeType string

const (
	baseUrl = "https://www.googleapis.com/drive/v3"

	MimeTypeNone   MimeType = ""
	MimeTypeFolder          = "application/vnd.google-apps.folder"
	MimeTypeFile            = "application/vnd.google-apps.file"
)

var (
	authScopes = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/drive.appdata",
		"https://www.googleapis.com/auth/drive.file",
		"https://www.googleapis.com/auth/drive.metadata",
		"https://www.googleapis.com/auth/drive.metadata.readonly",
		"https://www.googleapis.com/auth/drive.photos.readonly",
		"https://www.googleapis.com/auth/drive.readonly",
	}
)

type Api struct {
	ds           *drive.Service
	rc           *req.Client
	apiKey       string
	oauth2token  *oauth2.Token
	oauth2config *oauth2.Config
}

type ApiConfig struct {
	ApiKey       string
	ClientID     string
	ClientSecret string
	CallBackUrl  string
}

func NewApi(ctx context.Context, config ApiConfig, ts *oauth2.Token) (*Api, error) {
	oauth2Config := &oauth2.Config{
		RedirectURL:  config.CallBackUrl,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       authScopes,
		Endpoint:     google.Endpoint,
	}

	ds, err := drive.NewService(ctx, option.WithTokenSource(oauth2Config.TokenSource(ctx, ts)))
	if err != nil {
		return nil, err
	}

	return &Api{
		rc:           req.C(),
		apiKey:       config.ApiKey,
		oauth2token:  ts,
		oauth2config: oauth2Config,
		ds:           ds,
	}, nil
}

func (a *Api) apiUrl(route string, params ...string) *url.URL {
	u, _ := url.Parse(baseUrl)
	u.Path = path.Join(u.Path, route)

	q := u.Query()
	q.Add("key", a.apiKey)

	for i, _ := range params {
		if (i+1)%2 == 0 {
			q.Add(params[i-1], params[i])
		}
	}

	u.RawQuery = q.Encode()

	return u
}

func (a *Api) InitAuth(force bool) string {
	u := a.oauth2config.AuthCodeURL("state")
	au, _ := url.Parse(u)

	if force {
		q := au.Query()
		q.Add("prompt", "consent")
		q.Add("access_type", "offline")
		au.RawQuery = q.Encode()
	}
	
	return au.String()
}

func (a *Api) VerifyAuth(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := a.oauth2config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	a.oauth2token = token

	return token, nil
}

func (a *Api) About() (*drive.User, error) {
	about, err := a.ds.About.Get().Fields("*").Do()
	if err != nil {
		return nil, err
	}

	return about.User, nil
}

func (a *Api) FileList(mimeType MimeType, parentId string) (*drive.FileList, error) {
	q := "not name = ''"
	if mimeType != MimeTypeNone {
		q = fmt.Sprintf("mimeType = '%s'", mimeType)
	}

	parentId = strings.TrimSpace(parentId)
	if parentId != "" {
		q = fmt.Sprintf(`%s and '%s' in parents`, q, parentId)
	}

	return a.ds.Files.List().Q(q).Do()
}

func (a *Api) CreateDirectory(name, parentId string) (*drive.File, error) {
	return a.ds.Files.Create(&drive.File{
		Name:     name,
		MimeType: MimeTypeFolder,
		Parents:  []string{parentId},
	}).Do()
}

func (a *Api) DeleteDirectory(id string) error {
	return a.ds.Files.Delete(id).Do()
}

func (a *Api) GetFileMeta(id string) (*drive.File, error) {
	return a.ds.Files.Get(id).Fields("kind,id,name,mimeType,size").Do()
}

func (a *Api) DownloadFile(id string, w io.Writer) error {
	u := a.apiUrl("/files", "alt", "media").JoinPath(id)

	res, err := req.R().
		SetOutput(w).
		SetBearerAuthToken(a.oauth2token.AccessToken).
		Get(u.String())
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return nil
}

func (a *Api) UploadFile(name, parentId string, media io.Reader) (*drive.File, error) {
	return a.ds.Files.Create(
		&drive.File{Name: name, Parents: []string{parentId}}).
		Media(media).Do()
}
