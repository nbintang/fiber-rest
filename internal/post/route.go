package post

import (
	"rest-fiber/internal/contract"
	"rest-fiber/pkg/httpx"

	"github.com/gofiber/fiber/v2"
)

type PostRouteParams struct {
	httpx.RouteParams
	PostHandler PostHandler
}
type postRouteImpl struct {
	postHandler PostHandler
}

func NewPostRoute(params PostRouteParams) contract.ProtectedRoute {
	return &postRouteImpl{postHandler: params.PostHandler}
}
func (r *postRouteImpl) RegisterProtectedRoute(api fiber.Router) {
	posts := api.Group("/posts")
	posts.Get("/", r.postHandler.GetAllPosts)
	posts.Get("/:id", r.postHandler.GetPostByID)
	posts.Post("/", r.postHandler.CreatePost)
}
