package handlers

import (
	"adbr.xx/auth_microservice/auth"
	"adbr.xx/auth_microservice/database"
	"adbr.xx/auth_microservice/models"
	"adbr.xx/auth_microservice/utils"
	"github.com/gofiber/fiber/v3"
)

func Login(c fiber.Ctx) {
	url := auth.GetConfig().AuthCodeURL("state")
	c.Redirect().Status(fiber.StatusTemporaryRedirect).To(url)
}

func GoogleCallback(c fiber.Ctx) {
	token, err := auth.GetConfig().Exchange(c.RequestCtx(), c.FormValue("code"))
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString("Error intercambiando el token")
		return
	}
	user, err := auth.GetUser(token.AccessToken)
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
		return
	}

	database.GetOrCreateUser(&user)

	// Guarda o crea el usuario en tu base de datos y genera un JWT para el usuario
	tokens, err := auth.CreateJWTToken(user)
	if err != nil {
		c.SendStatus(500)
		return
	}

	c.JSON(fiber.Map{
		"message": "Autenticación exitosa",
		"user":    user,
		"tokens":  tokens,
	})
}

func GetUser(c fiber.Ctx) {
	id := c.Params("id", "-1")
	exists, user := database.GetUserById(id)
	if exists < 1 {
		c.JSON(fiber.Map{
			"message": "Usuario no encontrado",
		})
		return
	}
	c.JSON(fiber.Map{
		"user": user,
	})
}

func DeleteUser(c fiber.Ctx) {
	subject, err := utils.GetUserFromContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
		return
	}
	affectedRows := database.DeleteUser(subject)
	if affectedRows < 1 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Usuario no encontrado",
		})
		return
	}
	c.JSON(fiber.Map{
		"message": "Usuario eliminado exitosamente",
	})
}

func UpdateUser(c fiber.Ctx) {
	var user models.User
	if err := c.Bind().JSON(&user); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Datos de usuario inválidos",
		})
		return
	}
	subject, err := utils.GetUserFromContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error al obtener el usuario",
		})
		return
	}
	user.ID = subject
	affectedRows := database.UpdateUser(&user)
	if affectedRows < 1 {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Usuario no encontrado",
		})
		return
	}
}

func RefreshToken(c fiber.Ctx) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind().JSON(&request); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Datos de solicitud inválidos",
		})
		return
	}
	tokens, err := auth.RefreshToken(request.RefreshToken)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al refrescar el token",
		})
		return
	}
	c.JSON(fiber.Map{
		"tokens": tokens,
	})
}
