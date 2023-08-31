package worker

import (
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/av-prac-task/internal/model"
)

type ITask interface {
	Do()
}

type GenHistoryTask struct {
	input    model.UserSegmentsHistoryInput
	filePath string
	doFunc   func(input model.UserSegmentsHistoryInput, filePath string) error
}

func NewGenHistoryTask(
	input model.UserSegmentsHistoryInput,
	filePath string,
	doFunc func(input model.UserSegmentsHistoryInput, filePath string) error) *GenHistoryTask {
	return &GenHistoryTask{input: input, filePath: filePath, doFunc: doFunc}
}

func (t GenHistoryTask) Do() {
	err := t.doFunc(t.input, t.filePath)
	if err != nil {
		log.Errorf("GenHistoryTask.Do got error: %v", err)
	}
}
