package database

import (
	"context"

	"github.com/unbeman/av-prac-task/internal/config"
	"github.com/unbeman/av-prac-task/internal/model"
)

type IDatabase interface {
	CreateSegment(ctx context.Context, segment *model.Segment) (*model.Segment, error)
	DeleteSegment(ctx context.Context, segment *model.Segment) error
	GetSegment(ctx context.Context, segment *model.Segment) (*model.Segment, error)
	GetSegments(ctx context.Context, slugs []model.Slug) ([]model.Segment, error)
	CreateDeleteUserSegments(ctx context.Context, user *model.User, SegSlugsForCreate []model.Slug, SegSlugsForDelete []model.Slug) error
	GetUserSegments(ctx context.Context, input *model.User) (*model.User, error)
}

func GetDatabase(cfg config.PostgresConfig) (IDatabase, error) {
	return NewPGDatabase(cfg)
}
