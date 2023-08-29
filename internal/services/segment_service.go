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

func (s SegmentService) CreateSegment(ctx context.Context, input *model.CreateSegment) (*model.Segment, error) {
	var err error
	segment := &model.Segment{Slug: input.Slug}

	segment, err = s.db.CreateSegment(ctx, segment)
	if err != nil {
		return nil, err
	}

	if input.Selection == nil {
		return segment, nil
	}

	err = s.db.AddSegmentToRandomUsers(ctx, segment, *input.Selection)
	if err != nil {
		return nil, err
	}

	return segment, nil
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
