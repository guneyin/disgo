package google

import (
	"context"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type MimeType string

const (
	baseUrl = "https://www.googleapis.com/drive/v3"

	MimeTypeFolder MimeType = "application/vnd.google-apps.folder"
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
	rc           *req.Client
	apiKey       string
	oauth2token  oauth2.Token
	oauth2config oauth2.Config
}

type ApiConfig struct {
	ApiKey       string
	ClientID     string
	ClientSecret string
	CallBackUrl  string
}

func NewApi(config ApiConfig, accessToken string) *Api {
	return &Api{
		rc:          req.C(),
		apiKey:      config.ApiKey,
		oauth2token: oauth2.Token{AccessToken: accessToken},
		oauth2config: oauth2.Config{
			RedirectURL:  config.CallBackUrl,
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Scopes:       authScopes,
			Endpoint:     google.Endpoint,
		},
	}
}

func (a *Api) apiUrl(route string, params ...string) *url.URL {
	u, _ := url.Parse(baseUrl)
	u.Path = path.Join(u.Path, route)

	q := u.Query()
	q.Add("key", a.apiKey)
	q.Add("oauth_token", a.oauth2token.AccessToken)

	for i, _ := range params {
		if (i+1)%2 == 0 {
			q.Add(params[i-1], params[i])
		}
	}

	u.RawQuery = q.Encode()

	return u
}

func (a *Api) InitAuth() string {
	return a.oauth2config.AuthCodeURL("state")
}

func (a *Api) VerifyAuth(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := a.oauth2config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	a.oauth2token = *token

	return token, nil
}

func (a *Api) About() (*User, error) {
	data := &User{}

	u := a.apiUrl("/about", "fields", "user")
	res, err := req.R().
		SetSuccessResult(data).
		Get(u.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	return data, nil
}

func (a *Api) FileList(mimeType MimeType, parentId string) (*FileList, error) {
	data := &FileList{}

	q := fmt.Sprintf("mimeType = '%s'", mimeType)
	parentId = strings.TrimSpace(parentId)
	if parentId != "" {
		q = fmt.Sprintf(`%s and '%s' in parents`, q, parentId)
	}

	u := a.apiUrl("/files", "q", q)
	res, err := req.R().
		SetSuccessResult(data).
		Get(u.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	return data, nil
}

func (a *Api) CreateDirectory(name, parentId string) (*File, error) {
	data := &File{}

	body := struct {
		Name     string   `json:"name"`
		MimeType string   `json:"mimeType"`
		Parents  []string `json:"parents"`
	}{
		Name:     name,
		MimeType: string(MimeTypeFolder),
		Parents:  []string{parentId},
	}

	u := a.apiUrl("/files")
	res, err := req.DevMode().R().
		SetBody(body).
		SetSuccessResult(data).
		Post(u.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	return data, nil
}

func (a *Api) DeleteDirectory(id string) error {
	u := a.apiUrl("/files").JoinPath(id)

	res, err := req.R().
		SetBody(map[string]any{"trashed": true}).
		Patch(u.String())
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return nil
}
