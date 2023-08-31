package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/unbeman/av-prac-task/internal/database"
	"github.com/unbeman/av-prac-task/internal/model"
	"github.com/unbeman/av-prac-task/internal/utils"
	"github.com/unbeman/av-prac-task/internal/worker"
)

type UserService struct {
	db      database.IDatabase
	wp      *worker.WorkersPool
	fileDir string
}

func NewUserService(db database.IDatabase, wp *worker.WorkersPool, fileDir string) (*UserService, error) {
	if _, err := os.Stat(fileDir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(fileDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return &UserService{db: db, wp: wp, fileDir: fileDir}, nil
}

func (s UserService) UpdateUserSegments(ctx context.Context, input *model.UserSegmentsInput) error {
	user := model.User{}
	user.ID = input.UserID

	return s.db.CreateDeleteUserSegments(ctx, &user, input.SegmentsToAdd, input.SegmentsToDelete)
}

func (s UserService) GetUserActiveSegments(ctx context.Context, input *model.UserInput) (model.Slugs, error) {
	user := &model.User{}
	user.ID = input.UserID

	user, err := s.db.GetUserWithActiveSegments(ctx, user)
	if err != nil {
		return nil, err
	}

	slugs := make(model.Slugs, 0, len(user.Segments))
	for _, segment := range user.Segments {
		slugs = append(slugs, segment.Slug)
	}

	return slugs, nil
}

func (s UserService) generateUserSegmentsHistoryFile(input model.UserSegmentsHistoryInput, filePath string) error {
	user := &model.User{}
	user.ID = input.UserID

	ctx := context.TODO()
	userSegments, err := s.db.GetUserSegmentsHistory(ctx, user, input.FromDate, input.ToDate)
	if err != nil {
		return err
	}

	return utils.SaveCSVHistory(input, filePath, userSegments)
}

func (s UserService) GenerateUserSegmentsHistoryFile(ctx context.Context, input *model.UserSegmentsHistoryInput) (string, error) {
	user := &model.User{}
	user.ID = input.UserID

	_, err := s.db.GetUser(ctx, user)
	if err != nil {
		return "", err
	}

	filename := utils.FormatCSVFileName(input.UserID, input.FromDate, input.ToDate)

	filePath := utils.FormatCSVFilePath(s.fileDir, filename)

	// if the history is requested for the interval, including today, then history will gen again
	if input.ToDate.After(time.Now()) {
		s.wp.AddTask(worker.NewGenHistoryTask(*input, filePath, s.generateUserSegmentsHistoryFile))
		return filename, nil
	}

	// if file already exist, no need gen the new one
	if err = utils.CheckFileExists(filePath); err == nil {
		return filename, nil
	}

	// adding task to workers for gen history
	s.wp.AddTask(worker.NewGenHistoryTask(*input, filePath, s.generateUserSegmentsHistoryFile))

	return filename, nil
}

func (s UserService) DownloadUserSegmentsHistory(filename string) (string, error) {
	filePath := utils.FormatCSVFilePath(s.fileDir, filename)
	if err := utils.CheckFileExists(filePath); err != nil {
		return "", err
	}
	return filePath, nil
}
