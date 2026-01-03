package user

import "github.com/gofiber/fiber/v2"



func RegisterUserRoutes(app *fiber.App, h UserHandler){
	api := app.Group("/api")
	
	users := api.Group("/users")
	users.Get("/", h.GetAllUsers) 
	
}
