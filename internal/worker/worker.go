package worker

import (
	log "github.com/sirupsen/logrus"
	"sync"

	"github.com/unbeman/av-prac-task/internal/config"
)

type WorkersPool struct {
	wokersCount int
	tasks       chan ITask
	tasksSize   int
	waitGroup   sync.WaitGroup
}

func NewWorkersPool(cfg config.WorkerPoolConfig) *WorkersPool {
	return &WorkersPool{
		wokersCount: cfg.WorkersCount,
		tasks:       make(chan ITask, cfg.TasksSize),
	}
}

func (wp *WorkersPool) Run() {
	log.Infof("starting worker pool %d workers", wp.wokersCount)
	for idx := 0; idx < wp.wokersCount; idx++ {
		wp.waitGroup.Add(1)

		go func(idx int, tasks chan ITask) {
			defer wp.waitGroup.Done()
			log.Infof("worker %d started", idx)

			for task := range tasks {
				log.Infof("worker %d: starting task", idx)

				task.Do() //todo: try N times if err

				log.Infof("worker %d: ends task", idx)
			}
			log.Infof("worker %d: finished", idx)
		}(idx, wp.tasks)
	}
	wp.waitGroup.Wait()
}

func (wp *WorkersPool) AddTask(task ITask) {
	wp.tasks <- task
}

func (wp *WorkersPool) Shutdown() {
	close(wp.tasks)
}
