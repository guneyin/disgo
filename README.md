
# Current Status: W.I.P

Please DO NOT use in production

# disgo: A very basic cloud disk manager.
Add cloud disk capability to your golang applications

## Usage

```go
	gc := provider.GoogleConfig{
		ApiKey:       os.Getenv("DSG_GOOGLE_API_KEY"),
		ClientID:     os.Getenv("DSG_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("DSG_GOOGLE_SECRET"),
		CallBackUrl:  "http://localhost:8080/auth/google/callback",
	}

	accessToken := os.Getenv("DSG_GOOGLE_OAUTH2_TOKEN")

	gp := provider.NewGoogle(gc, []byte(accessToken))
	user, err := gp.GetAuthorizedUser()
	if err != nil {
		log.Fatal(err)
	}
	userData, _ := json.MarshalIndent(user, "", " ")
	_ = fmt.Sprintf("%-v", string(userData))

	dirs, err := gp.GetDirectoryList("")
	if err != nil {
		log.Fatal(err)
	}
	dirsData, _ := json.MarshalIndent(dirs, "", " ")
	_ = fmt.Sprintf("%-v", string(dirsData))

```

  
