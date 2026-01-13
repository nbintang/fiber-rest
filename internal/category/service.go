package category

import (
	"context"
	"errors"
	"rest-fiber/pkg/slice"
)

type categoryServiceimpl struct {
	categoryRepo CategoryRepository
}

func NewCategoryService(categoryRepo CategoryRepository) CategoryService {
	return &categoryServiceimpl{categoryRepo}
}

func (s *categoryServiceimpl) FindAllCategories(ctx context.Context, page, limit, offset int) ([]CategoryResponseDTO, int64, error) {
	categories, total, err := s.categoryRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	categoriesResponse := slice.Map[Category, CategoryResponseDTO](
		categories,
		func(c Category) CategoryResponseDTO {
			return CategoryResponseDTO{
				ID:   c.ID,
				Name: c.Name,
			}
		},
	)

	return categoriesResponse, total, nil
}

func (s *categoryServiceimpl) FindCategoryByID(ctx context.Context, id string) (*CategoryResponseDTO, error) {
	existed, err := s.categoryRepo.ExistsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, errors.New("Category Not Found")
	}

	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &CategoryResponseDTO{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

func (s *categoryServiceimpl) CreateCategory(ctx context.Context, dto *CategoryRequestDTO) (*CategoryResponseDTO, error) {
	category := &Category{Name: dto.Name}
	created, err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	category, err = s.categoryRepo.FindByID(ctx, created.String())
	if err != nil {
		return nil, err
	}

	return &CategoryResponseDTO{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}
func (s *categoryServiceimpl) UpdateCategoryByID(ctx context.Context, id string, dto *CategoryRequestDTO) (*CategoryResponseDTO, error) {
	existed, err := s.categoryRepo.ExistsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, errors.New("Category Not Found")
	}

	category := &Category{Name: dto.Name}
	updated, err := s.categoryRepo.Update(ctx, id, category)
	if err != nil {
		return nil, err
	}

	category, err = s.categoryRepo.FindByID(ctx, updated.String())
	if err != nil {
		return nil, err
	}

	return &CategoryResponseDTO{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}
func (s *categoryServiceimpl) DeleteCategoryByID(ctx context.Context, id string) error {
	existed, err := s.categoryRepo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !existed {
		return errors.New("Category Not Found")
	}
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
