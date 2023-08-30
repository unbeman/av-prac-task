package database

import (
	"context"
	"time"

	"github.com/unbeman/av-prac-task/internal/config"
	"github.com/unbeman/av-prac-task/internal/model"
)

type IDatabase interface {
	CreateSegment(ctx context.Context, segment *model.Segment) (*model.Segment, error)
	AddSegmentToRandomUsers(ctx context.Context, segment *model.Segment, selection float64) error
	DeleteSegment(ctx context.Context, segment *model.Segment) error
	GetSegment(ctx context.Context, segment *model.Segment) (*model.Segment, error)
	GetSegments(ctx context.Context, slugs []model.Slug) ([]*model.Segment, error)
	CreateDeleteUserSegments(ctx context.Context, user *model.User, SegSlugsForCreate []model.Slug, SegSlugsForDelete []model.Slug) error
	GetUserActiveSegments(ctx context.Context, input *model.User) (*model.User, error)
	GetUserSegmentsHistory(ctx context.Context, user *model.User, from time.Time, to time.Time) ([]model.UserSegment, error)
}

func GetDatabase(cfg config.PostgresConfig) (IDatabase, error) {
	return NewPGDatabase(cfg)
}
