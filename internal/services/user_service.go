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

	//todo: what to do if segments to add intersect to delete? should i check it before operations?
	return s.db.CreateDeleteUserSegments(ctx, &user, input.SegmentsToAdd, input.SegmentsToDelete)
}

func (s UserService) GetActiveUserSegments(ctx context.Context, user *model.User) (*model.User, error) {
	user, err := s.db.GetUserSegments(ctx, user)
	return user, err
}
