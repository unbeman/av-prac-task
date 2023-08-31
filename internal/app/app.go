package app

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/av-prac-task/internal/config"
	"github.com/unbeman/av-prac-task/internal/database"
	"github.com/unbeman/av-prac-task/internal/handlers"
	"github.com/unbeman/av-prac-task/internal/services"
	"github.com/unbeman/av-prac-task/internal/worker"
)

type SegApp struct {
	server      http.Server
	workersPool *worker.WorkersPool
}

func GetSegApp(cfg config.AppConfig) (*SegApp, error) {
	db, err := database.GetDatabase(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("coudnt get database: %w", err)
	}

	wp := worker.NewWorkersPool(cfg.WorkersPool)

	uServ, err := services.NewUserService(db, wp, cfg.FileDirectory)
	if err != nil {
		return nil, fmt.Errorf("coudnt get user service: %w", err)
	}

	sServ, err := services.NewSegmentService(db)
	if err != nil {
		return nil, fmt.Errorf("coudnt get segment service: %w", err)
	}

	handler, err := handlers.GetHandler(uServ, sServ)
	if err != nil {
		return nil, fmt.Errorf("coudnt get handler: %w", err)
	}

	server := http.Server{
		Addr:    cfg.Address,
		Handler: handler,
	}

	application := &SegApp{
		server:      server,
		workersPool: wp,
	}
	return application, nil
}

func (a *SegApp) Run() {
	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		a.workersPool.Run()
		log.Info("worker pool finished")
	}()

	go func() {
		defer wg.Done()
		err := a.server.ListenAndServe()
		log.Infof("server closed: %s", err)
	}()

	log.Infof("application started, Addr: %s", a.server.Addr)

	wg.Wait()
}

func (a *SegApp) Stop() {
	log.Infof("shutting down")
	err := a.server.Shutdown(context.TODO())
	if err != nil {
		log.Error(err)
	}
	a.workersPool.Shutdown()
}
