package app

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/unbeman/av-prac-task/internal/config"
	"github.com/unbeman/av-prac-task/internal/database"
	"github.com/unbeman/av-prac-task/internal/handlers"
	"github.com/unbeman/av-prac-task/internal/services"
)

type SegApp struct {
	server http.Server
}

func GetSegApp(cfg config.AppConfig) (*SegApp, error) {
	db, err := database.GetDatabase(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("coudnt get database: %w", err)
	}

	uServ, err := services.NewUserService(db, cfg.FileDirectory)
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
	application := &SegApp{
		server: http.Server{
			Addr:    cfg.Address,
			Handler: handler,
		},
	}
	return application, nil
}

func (a *SegApp) Run(ctx context.Context) {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.server.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		return a.server.Shutdown(context.Background())
	})

	log.Infof("application started, Addr: %s", a.server.Addr)

	if err := g.Wait(); err != nil {
		log.Infof("application stopped, reason: %s", err)
	}
}
