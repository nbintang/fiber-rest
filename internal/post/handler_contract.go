package post

import "github.com/gofiber/fiber/v2"

type PostHandler interface {
	GetAllPosts(c *fiber.Ctx) error
	GetPostByID(c *fiber.Ctx) error
	CreatePost(c *fiber.Ctx) error
}
