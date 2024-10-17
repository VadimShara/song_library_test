package user

import "context"

type Repository interface {
	Create(ctx context.Context, s *Song) error
	FindAll(ctx context.Context) ([]Song, error)
	FindWithFilter(ctx context.Context, s *Song) ([]Song, error)
	FindOne(ctx context.Context, s *Song) (*Song, error)
	Delete(ctx context.Context, s *Song) error
	Update(ctx context.Context, s *Song) error
}