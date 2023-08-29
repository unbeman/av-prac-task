package services

import (
	"context"

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

func (s UserService) GetUserWithActiveSegments(ctx context.Context, input *model.UserInput) (*model.User, error) {
	user := &model.User{}
	user.ID = input.UserID

	return s.db.GetUserActiveSegments(ctx, user)
}

func (s UserService) GetUserWithSegmentsHistory(ctx context.Context, input *model.UserSegmentsHistoryInput) (*model.User, error) {
	user := &model.User{}
	user.ID = input.UserID
	return s.db.GetUserSegmentsHistory(ctx, user, input.FromDate, input.ToDate)
}
