package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type AuthToken struct {
	AccessToken   string
	CookieSession string
	ClientTable   string
	UrlApi        string
}

func GetToken() (AuthToken, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
		return AuthToken{}, err
	}

	access_token := os.Getenv("ACCESS_TOKEN")
	cookie_session := os.Getenv("COOKIES_SESSION")
	client_table := os.Getenv("CLIENT_TABLE")
	url_api := os.Getenv("URLAPI")

	token := AuthToken{
		AccessToken:   access_token,
		CookieSession: cookie_session,
		ClientTable:   client_table,
		UrlApi:        url_api,
	}

	return token, nil
}
