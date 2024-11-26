package google

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"testing"
)

const (
	callbackUrl = "/auth/google/callback"
	port        = ":8080"
)

func NewTestApi(t *testing.T) *Api {
	t.Helper()
	err := godotenv.Load("../../.env.local")
	if err != nil {
		t.Fatalf("Error loading .env file")
	}

	return NewApi(ApiConfig{
		ApiKey:       os.Getenv("DSG_GOOGLE_API_KEY"),
		ClientID:     os.Getenv("DSG_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("DSG_GOOGLE_SECRET"),
		CallBackUrl:  fmt.Sprintf("http://localhost%s%s", port, callbackUrl),
	}, ReadToken())
}

func ServeHTTP(ctx context.Context, api *Api, tokenCh chan<- string) {
	mux := http.NewServeMux()
	mux.HandleFunc(callbackUrl, func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := api.VerifyAuth(ctx, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		tokenCh <- token.AccessToken
	})

	log.Fatal(http.ListenAndServe(port, mux))
}

func WriteToken(t string) error {
	f, err := os.Create("../../token.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(t))
	return err
}

func ReadToken() string {
	f, err := os.ReadFile("../../token.txt")
	if err != nil {
		return ""
	}

	return string(f)
}
