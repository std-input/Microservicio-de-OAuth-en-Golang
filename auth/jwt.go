package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"adbr.xx/auth_microservice/configs"
	"adbr.xx/auth_microservice/database"
	"adbr.xx/auth_microservice/models"
	"github.com/golang-jwt/jwt/v5"
)

type UserTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func hash_secret(secret string) string {
	hasher := sha256.New()
	hasher.Write([]byte(secret + configs.Get("SECRET_KEY")))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Crea un JWT token y un refresh token para el usuario con la ID dada
func CreateJWTToken(user models.User) (UserTokens, error) {
	if user.ID == "" {
		return UserTokens{}, errors.New("Usuario invalido")
	}

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
	hashedToken := hash_secret(refreshToken)

	userRefreshToken := models.RefreshToken{
		UserID:     user.ID,
		Token:      hashedToken,
		Expiration: time.Now().Add(time.Hour * 24), // El token de refresco expira en 1 dia
	}
	db.Where("user_id = ?", user.ID).FirstOrCreate(&userRefreshToken)
	userRefreshToken.Token = hashedToken
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
	results := db.Preload("User").Where("token = ?", hash_secret(token)).First(&refreshToken)
	if results.RowsAffected == 0 {
		return UserTokens{}, errors.New("Token de refresco no valido")
	}
	if refreshToken.Expiration.Before(time.Now()) {
		return UserTokens{}, errors.New("Token de refresco expirado")
	}
	return CreateJWTToken(refreshToken.User)
}
