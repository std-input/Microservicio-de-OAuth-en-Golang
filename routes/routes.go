package routes

import (
	"adbr.xx/auth_microservice/configs"
	"adbr.xx/auth_microservice/handlers"
	"github.com/gofiber/fiber/v3"

	jwtware "github.com/gofiber/contrib/v3/jwt"
)

func Routes(app *fiber.App) {
	api_v1 := app.Group("/api")

	api_v1.Get("/auth", handlers.Login)

	api_v1.Get("/auth/google/callback", handlers.GoogleCallback)

	api_v1.Get("/user/:id", handlers.GetUser)

	api_v1.Post("/auth/refresh_token", handlers.RefreshToken)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(configs.Get("SECRET_KEY")),
		},
	}))

	api_v1.Delete("/me", handlers.DeleteUser)

	api_v1.Put("/me", handlers.UpdateUser)

}
