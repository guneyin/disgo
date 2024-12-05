package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/guneyin/disgo"
	"github.com/guneyin/disgo/provider"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type handler struct {
	ctx context.Context
	p   disgo.Provider
}

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gc := provider.GoogleConfig{
		ApiKey:       os.Getenv("DSG_GOOGLE_API_KEY"),
		ClientID:     os.Getenv("DSG_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("DSG_GOOGLE_SECRET"),
		CallBackUrl:  "http://localhost:8080/auth/google/callback",
	}

	gp, err := provider.NewGoogle(ctx, gc, nil)
	if err != nil {
		log.Fatal(err)
	}

	hnd := &handler{
		ctx: ctx,
		p:   gp,
	}

	fmt.Println("http://localhost:8080/auth/google/init")
	fmt.Println("http://localhost:8080/user")

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/google/init", hnd.authInit)
	mux.HandleFunc("/auth/google/callback", hnd.authCallback)
	mux.HandleFunc("/user", hnd.user)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (h *handler) authInit(w http.ResponseWriter, r *http.Request) {
	url := h.p.InitAuth()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handler) authCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := h.p.VerifyAuth(h.ctx, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data, _ := json.Marshal(token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (h *handler) user(w http.ResponseWriter, r *http.Request) {
	user, err := h.p.GetAuthorizedUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data, _ := json.Marshal(user)
	_, _ = w.Write(data)
}
