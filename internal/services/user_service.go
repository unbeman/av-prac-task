package services

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/av-prac-task/internal/database"
	"github.com/unbeman/av-prac-task/internal/model"
)

type UserService struct {
	db      database.IDatabase
	fileDir string
}

func NewUserService(db database.IDatabase, fileDir string) (*UserService, error) {
	if _, err := os.Stat(fileDir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(fileDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return &UserService{db: db, fileDir: fileDir}, nil
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

// todo: async load and gen file
func (s UserService) GenerateUserSegmentsHistoryFile(ctx context.Context, input *model.UserSegmentsHistoryInput) (string, error) {
	user := &model.User{}
	user.ID = input.UserID

	filename := fmt.Sprintf(
		"user-%d_%s_%s.csv",
		user.ID,
		input.FromDate.Format(time.DateOnly),
		input.ToDate.Format(time.DateOnly),
	)

	filePath := fmt.Sprintf("%s/%s", s.fileDir, filename)

	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		//todo: check file, rewrite it if new records exists
		return filename, err
	}

	userSegments, err := s.db.GetUserSegmentsHistory(ctx, user, input.FromDate, input.ToDate)
	if err != nil {
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Error("GenerateUserSegmentsHistoryFile: ", err)
		}
	}(file)

	csvFile := csv.NewWriter(file)

	head := []string{"user_id", "segment_slug", "operation", "date"}
	if err = csvFile.Write(head); err != nil {
		return "", err
	}

	for _, segment := range userSegments {
		if segment.CreatedAt.After(input.FromDate) && segment.CreatedAt.Before(input.ToDate) {
			row := []string{strconv.Itoa(int(user.ID)), string(segment.Segment.Slug), "add", segment.CreatedAt.String()}

			if err = csvFile.Write(row); err != nil {
				return "", err
			}
		}

		if segment.DeletedAt.Valid &&
			segment.DeletedAt.Time.After(input.FromDate) &&
			segment.DeletedAt.Time.Before(input.ToDate) {
			row := []string{strconv.Itoa(int(user.ID)), string(segment.Segment.Slug), "delete", segment.DeletedAt.Time.String()}

			if err = csvFile.Write(row); err != nil {
				return "", err
			}
		}
	}

	csvFile.Flush()

	return filename, nil
}

func (s UserService) DownloadUserSegmentsHistory(filename string) (string, error) {
	filePath := fmt.Sprintf("%s/%s", s.fileDir, filename)
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return "", ErrFileNotFound
	}

	return filePath, nil
}
