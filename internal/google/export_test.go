package google

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"testing"
)

const (
	callbackUrl = "/auth/google/callback"
	port        = ":8080"
)

func NewTestApi(ctx context.Context, t *testing.T) (*Api, error) {
	err := godotenv.Load("../../.env.local")
	if err != nil {
		t.Fatalf("Error loading .env file")
	}

	return NewApi(ctx, ApiConfig{
		ApiKey:       os.Getenv("DSG_GOOGLE_API_KEY"),
		ClientID:     os.Getenv("DSG_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("DSG_GOOGLE_SECRET"),
		CallBackUrl:  fmt.Sprintf("http://localhost%s%s", port, callbackUrl),
	}, ReadToken())
}

func ServeHTTP(ctx context.Context, api *Api, tokenCh chan<- *oauth2.Token) {
	mux := http.NewServeMux()
	mux.HandleFunc(callbackUrl, func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := api.VerifyAuth(ctx, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		tokenCh <- token
	})

	log.Fatal(http.ListenAndServe(port, mux))
}

func WriteToken(t *oauth2.Token) error {
	f, err := os.Create("testdata/token.json")
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

func ReadToken() *oauth2.Token {
	t := &oauth2.Token{}

	f, err := os.ReadFile("testdata/token.json")
	if err != nil {
		return t
	}

	_ = json.Unmarshal(f, t)

	return t
}
