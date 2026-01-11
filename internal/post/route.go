package post

import (
	"rest-fiber/internal/middleware"
	"rest-fiber/pkg/httpx"
	"rest-fiber/utils/enums"

	"github.com/gofiber/fiber/v2"
)

type PostRouteParams struct {
	httpx.RouteParams
	PostHandler PostHandler
}
type postRouteImpl struct {
	postHandler PostHandler
}

func NewPostRoute(params PostRouteParams) httpx.ProtectedRoute {
	return &postRouteImpl{postHandler: params.PostHandler}
}
func (r *postRouteImpl) RegisterProtectedRoute(route fiber.Router) {
	posts := route.Group("/posts")
	posts.Get("/", r.postHandler.GetAllPosts)
	posts.Get("/:id", r.postHandler.GetPostByID)
	posts.Post("/", middleware.AuthAllowRoleAccess(enums.Admin), r.postHandler.CreatePost)
	posts.Patch("/:id", middleware.AuthAllowRoleAccess(enums.Admin), r.postHandler.UpdatePostByID)
	posts.Delete("/:id", middleware.AuthAllowRoleAccess(enums.Admin), r.postHandler.DeletePostByID)
}
