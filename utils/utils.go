package utils

import (
	"fmt"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
)

func GetUserFromContext(c fiber.Ctx) (string, error) {
	token := jwtware.FromContext(c)
	if token == nil {
		return "", fmt.Errorf("Token no encontrado en el contexto")
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	return subject, nil
}
