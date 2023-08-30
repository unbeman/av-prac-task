package services

import (
	"context"
	log "github.com/sirupsen/logrus"
	"strconv"

	"github.com/unbeman/av-prac-task/internal/database"
	"github.com/unbeman/av-prac-task/internal/model"
)

type UserService struct {
	db database.IDatabase
}

func NewUserService(db database.IDatabase) (*UserService, error) {
	return &UserService{db: db}, nil
}

func (s UserService) UpdateUserSegments(ctx context.Context, input *model.UserSegmentsInput) error {
	user := model.User{}
	user.ID = input.UserID

	return s.db.CreateDeleteUserSegments(ctx, &user, input.SegmentsToAdd, input.SegmentsToDelete)
}

func (s UserService) GetUserActiveSegments(ctx context.Context, input *model.UserInput) (model.Slugs, error) {
	user := &model.User{}
	user.ID = input.UserID

	user, err := s.db.GetUserActiveSegments(ctx, user)
	if err != nil {
		return nil, err
	}

	slugs := make(model.Slugs, 0, len(user.Segments))
	for _, segment := range user.Segments {
		slugs = append(slugs, segment.Slug)
	}

	return slugs, nil
}

func (s UserService) GetUserSegmentsHistory(ctx context.Context, input *model.UserSegmentsHistoryInput) ([][]string, error) {
	user := &model.User{}
	user.ID = input.UserID

	us, err := s.db.GetUserSegmentsHistory(ctx, user, input.FromDate, input.ToDate)
	if err != nil {
		return nil, err
	}

	history := make([][]string, 0, len(us)+1)

	head := []string{"user_id", "segment_slug", "operation", "date"}

	history = append(history, head)

	for _, segment := range us {
		if segment.CreatedAt.After(input.FromDate) && segment.CreatedAt.Before(input.ToDate) {
			row := []string{strconv.Itoa(int(user.ID)), string(segment.Segment.Slug), "add", segment.CreatedAt.String()}
			history = append(history, row)
		}

		if segment.DeletedAt.Valid &&
			segment.DeletedAt.Time.After(input.FromDate) &&
			segment.DeletedAt.Time.Before(input.ToDate) {
			row := []string{strconv.Itoa(int(user.ID)), string(segment.Segment.Slug), "delete", segment.DeletedAt.Time.String()}
			history = append(history, row)
		}
	}
	log.Info(history)
	return history, nil
}
