package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"adbr.xx/auth_microservice/configs"
	"adbr.xx/auth_microservice/database"
	"adbr.xx/auth_microservice/models"
	"github.com/golang-jwt/jwt/v5"
)

type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

// Crea un JWT token y un refresh token para el usuario con la ID dada
func CreateJWTToken(user models.User) (UserTokens, error) {

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Minute * 10).Unix(),
		"role": user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(configs.Get("SECRET_KEY")))
	if err != nil {
		return UserTokens{}, err
	}

	db := database.DB
	refreshToken := GenerateToken(user.ID)

	userRefreshToken := models.RefreshToken{
		UserID:     user.ID,
		Token:      refreshToken,
		Expiration: time.Now().Add(time.Hour * 24), // El token de refresco expira en 1 dia
	}
	db.Where("user_id = ?", user.ID).FirstOrCreate(&userRefreshToken)
	userRefreshToken.Token = refreshToken
	db.Save(&userRefreshToken)

	return UserTokens{
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
	}, nil
}

// Crea un token de refresco para el usuario con la ID dada
func GenerateToken(id string) string {
	b := make([]byte, 32)
	rand.Read(b)
	return id + ":" + base64.StdEncoding.EncodeToString(b)
}

// Verifica el token de refresco, si es valido devuelve un nuevo JWT token para el usuario
func RefreshToken(token string) (UserTokens, error) {
	// Verificar el token de refresco
	var refreshToken models.RefreshToken
	db := database.DB
	results := db.Where("token = ?", token).First(&refreshToken).Preload("User")
	if results.RowsAffected == 0 {
		return UserTokens{}, errors.New("Token de refresco no valido")
	}
	if refreshToken.Expiration.Before(time.Now()) {
		return UserTokens{}, errors.New("Token de refresco expirado")
	}
	return CreateJWTToken(refreshToken.User)
}
