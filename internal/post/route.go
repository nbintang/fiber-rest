package post

import (
	"rest-fiber/internal/enums"
	"rest-fiber/internal/http/middleware"
	"rest-fiber/internal/http/router"

	"github.com/gofiber/fiber/v2"
)

type PostRouteParams struct {
	router.RouteParams
	PostHandler PostHandler
}
type postRouteImpl struct {
	postHandler PostHandler
}

func NewPostRoute(params PostRouteParams) router.ProtectedRoute {
	return &postRouteImpl{postHandler: params.PostHandler}
}
func (r *postRouteImpl) RegisterProtectedRoute(route fiber.Router) {
	posts := route.Group("/posts")
	posts.Get("/", r.postHandler.GetAllPosts)
	posts.Get("/:id", r.postHandler.GetPostByID)
	posts.Post("/", middleware.AllowRoleAccess(enums.Admin), r.postHandler.CreatePost)
	posts.Patch("/:id", middleware.AllowRoleAccess(enums.Admin), r.postHandler.UpdatePostByID)
	posts.Delete("/:id", middleware.AllowRoleAccess(enums.Admin), r.postHandler.DeletePostByID)
}
