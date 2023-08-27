package services

import (
	"context"

	"github.com/unbeman/av-prac-task/internal/database"
	"github.com/unbeman/av-prac-task/internal/model"
)

type SegmentService struct {
	db database.IDatabase
}

func NewSegmentService(db database.IDatabase) (*SegmentService, error) {
	return &SegmentService{db: db}, nil
}

func (s SegmentService) CreateSegment(ctx context.Context, input *model.SegmentInput) (*model.Segment, error) {
	segment := model.Segment{Slug: input.Slug, Selection: input.Selection}
	return s.db.CreateSegment(ctx, &segment)
}

func (s SegmentService) DeleteSegment(ctx context.Context, input *model.SegmentInput) error {
	segment := model.Segment{Slug: input.Slug}
	return s.db.DeleteSegment(ctx, &segment)
}

// todo: remove
func (s SegmentService) GetSegment(ctx context.Context, input *model.SegmentInput) (*model.Segment, error) {
	segment := model.Segment{Slug: input.Slug}
	return s.db.GetSegment(ctx, &segment)
}
