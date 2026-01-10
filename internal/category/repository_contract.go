package category

import "context"

type CategoryRepository interface {
	ExistsByID(ctx context.Context, id string) (bool, error)
}
