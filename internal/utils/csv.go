package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"

	"github.com/unbeman/av-prac-task/internal/model"
)

const (
	OperationAdd    = "add"
	OperationDelete = "delete"
)

var ErrFileNotFound = errors.New("file not found")

func FormatCSVFileName(userID uint64, from, to time.Time) string {
	return fmt.Sprintf(
		"user-%d_%s_%s.csv",
		userID,
		from.Format(time.DateOnly),
		to.Format(time.DateOnly),
	)
}

func FormatCSVFilePath(saveDir, fileName string) string {
	return fmt.Sprintf("%s/%s", saveDir, fileName)
}

func CheckFileExists(filePath string) error {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return ErrFileNotFound
	}
	return nil
}

func SaveCSVHistory(input model.UserSegmentsHistoryInput, filePath string, userSegments []model.UserSegment) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Error("SaveCSVHistory: ", err)
		}
	}(file)

	csvFile := csv.NewWriter(file)

	head := []string{"user_id", "segment_slug", "operation", "date"}
	if err = csvFile.Write(head); err != nil {
		return err
	}

	for _, segment := range userSegments {
		if segment.CreatedAt.After(input.FromDate) && segment.CreatedAt.Before(input.ToDate) {
			row := []string{
				strconv.Itoa(int(input.UserID)),
				string(segment.Segment.Slug),
				OperationAdd,
				segment.CreatedAt.String(),
			}

			if err = csvFile.Write(row); err != nil {
				return err
			}
		}

		if segment.DeletedAt.Valid &&
			segment.DeletedAt.Time.After(input.FromDate) &&
			segment.DeletedAt.Time.Before(input.ToDate) {
			row := []string{
				strconv.Itoa(int(input.UserID)),
				string(segment.Segment.Slug),
				OperationDelete,
				segment.DeletedAt.Time.String(),
			}

			if err = csvFile.Write(row); err != nil {
				return err
			}
		}
	}

	csvFile.Flush()
	return nil
}
