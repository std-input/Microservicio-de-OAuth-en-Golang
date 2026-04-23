package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"adbr.xx/auth_microservice/configs"
	"adbr.xx/auth_microservice/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     configs.Get("GOOGLE_CLIENT_ID"),
		ClientSecret: configs.Get("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  configs.Get("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		Endpoint: google.Endpoint,
	}

	// Obtiene el URL de redirección para iniciar el proceso de autenticación con Google
	return conf
}

func GetUser(token string) (models.User, error) {
	reqURL, err := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		return models.User{}, err
	}
	ptoken := fmt.Sprintf("Bearer %s", token)
	res := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {ptoken},
		},
	}
	req, err := http.DefaultClient.Do(res)
	if err != nil {
		return models.User{}, err
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return models.User{}, err
	}
	var data models.User
	err = json.Unmarshal(body, &data)
	if err != nil {
		return models.User{}, err
	}
	return data, nil
}
