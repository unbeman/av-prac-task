package database

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/unbeman/av-prac-task/internal/config"
	"github.com/unbeman/av-prac-task/internal/model"
)

type pg struct {
	conn *gorm.DB
}

// NewPGDatabase returns the initialized pg object that implements IDatabase interface.
func NewPGDatabase(cfg config.PostgresConfig) (*pg, error) {
	db := &pg{}
	if err := db.connect(cfg.DSN); err != nil {
		return nil, err
	}
	if err := db.migrate(); err != nil {
		return nil, err
	}

	return db, nil
}

// connect initialize database session connection instance with dsn.
func (p *pg) connect(dsn string) error {
	log.Info("PG DSN: ", dsn)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true}) //todo: use custom logger based on logrus
	if err != nil {
		return err
	}
	p.conn = conn
	return nil
}

// migrate prepares database.
func (p *pg) migrate() error {
	err := p.conn.AutoMigrate(
		&model.User{},
		&model.Segment{},
		&model.UserSegment{},
	)
	if err != nil {
		return err
	}

	err = p.conn.SetupJoinTable(&model.User{}, "Segments", &model.UserSegment{})
	if err != nil {
		return err
	}

	err = p.conn.SetupJoinTable(&model.Segment{}, "Users", &model.UserSegment{})
	if err != nil {
		return err
	}
	return nil
}

// Ping checks database connection.
func (p *pg) Ping(ctx context.Context) error {
	sqlDB, err := p.conn.DB()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDB, err)
	}
	if err = sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrDB, err)
	}
	return nil
}

func (p *pg) getUsersCount(ctx context.Context) (int64, error) {
	var count int64
	result := p.conn.WithContext(ctx).Model(&model.User{}).Count(&count)
	if result.Error != nil {
		return count, result.Error
	}
	return count, nil
}

func (p *pg) createSegment(ctx context.Context, tx *gorm.DB, segment *model.Segment) (*model.Segment, error) {
	result := tx.WithContext(ctx).Create(segment)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fmt.Errorf("segment with slug (%s) %w", segment.Slug, ErrAlreadyExists)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return segment, nil
}

// CreateSegment inserts new segment with given unique slug if not exists.
func (p *pg) CreateSegment(ctx context.Context, segment *model.Segment) (*model.Segment, error) {
	return p.createSegment(ctx, p.conn, segment)
}

// DeleteSegment hard deletes segment by slug.
func (p *pg) DeleteSegment(ctx context.Context, segment *model.Segment) error {
	result := p.conn.WithContext(ctx).Delete(segment, "slug = ?", segment.Slug)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("segment with slug (%s) %w", segment.Slug, ErrNotFound)
	}
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return nil
}

// GetSegment returns segment with all users in it.
func (p *pg) GetSegment(ctx context.Context, segment *model.Segment) (*model.Segment, error) {
	result := p.conn.WithContext(ctx).Preload("Users").First(segment, "slug = ?", segment.Slug)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("segment with slug (%s) %w", segment.Slug, ErrNotFound)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return segment, nil
}

// GetSegments returns segments by given slugs.
func (p *pg) GetSegments(ctx context.Context, slugs []model.Slug) ([]model.Segment, error) {
	return p.getSegments(ctx, p.conn, slugs)
}

// getSegments returns segments by given slugs.
func (p *pg) getSegments(ctx context.Context, tx *gorm.DB, slugs []model.Slug) ([]model.Segment, error) {
	var segments []model.Segment
	result := tx.WithContext(ctx).Find(&segments, "slug IN ?", slugs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("segments %w", ErrNotFound)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return segments, nil
}

// CreateDeleteUserSegments insert and delete relation by specified segments (represented by slugs) for given user.
func (p *pg) CreateDeleteUserSegments(ctx context.Context, user *model.User, toInSegments []model.Slug, toDelSegments []model.Slug) error {
	var insertSegments []model.Segment
	var deleteSegments []model.Segment
	err := p.conn.Transaction(func(tx *gorm.DB) error { //todo: check gorm's tx errors
		var txErr error
		insertSegments, txErr = p.getSegments(ctx, tx, toInSegments)
		if txErr != nil {
			return txErr
		}

		deleteSegments, txErr = p.getSegments(ctx, tx, toDelSegments)
		if txErr != nil {
			return txErr
		}

		txErr = p.insertUserSegments(ctx, tx, user, insertSegments)
		if txErr != nil {
			return txErr
		}

		txErr = p.deleteUserSegments(ctx, tx, user, deleteSegments)
		if txErr != nil {
			return txErr
		}
		return nil
	})

	return err
}

// insertUserSegments connects user with given segments.
func (p *pg) insertUserSegments(ctx context.Context, tx *gorm.DB, user *model.User, segments []model.Segment) error {
	log.Info(user.ID)
	err := tx.WithContext(ctx).Model(user).Association("Segments").Append(segments)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// deleteUserSegments removes user relation to given segments.
func (p *pg) deleteUserSegments(ctx context.Context, tx *gorm.DB, user *model.User, segments []model.Segment) error {
	err := tx.WithContext(ctx).Model(user).Association("Segments").Delete(segments)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// GetUserActiveSegments returns user with related segments.
func (p *pg) GetUserActiveSegments(ctx context.Context, user *model.User) (*model.User, error) {
	result := p.conn.WithContext(ctx).Preload("Segments").First(user, "ID = ?", user.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user with ID (%d) %w", user.ID, ErrNotFound)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return user, nil
}

// GetUserSegmentsHistory returns history of user's segments changes for the specified time interval.
func (p *pg) GetUserSegmentsHistory(ctx context.Context, user *model.User, from time.Time, to time.Time) (*model.User, error) {
	result := p.conn.WithContext(ctx).
		Preload("Segments",
			"created_at > ? AND created_at < ? AND (deleted_at = NULL OR (deleted_at > ? AND deleted_at < ?))",
			from, to, from, to).
		First(user, "ID = ?", user.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user with ID (%d) %w", user.ID, ErrNotFound)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return user, nil
}

func (p *pg) getRandomUsers(ctx context.Context, count int) ([]*model.User, error) {
	var users []*model.User
	result := p.conn.WithContext(ctx).
		Find(&users).
		Order("random()").
		Limit(count)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (p *pg) addSegmentToUsers(ctx context.Context, segment model.Segment, users []*model.User) error {
	for _, user := range users {
		if err := p.insertUserSegments(ctx, p.conn, user, []model.Segment{segment}); err != nil {
			return err
		}
	}
	return nil
}

func (p *pg) AddSegmentToRandomUsers(ctx context.Context, segment *model.Segment, selection float64) error {
	count, err := p.getUsersCount(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		log.Info("AddSegmentToRandomUsers: no users in database to add segment")
		return nil
	}

	selectionCount := int(math.Ceil(float64(count) * selection)) //todo: wrap

	users, err := p.getRandomUsers(ctx, selectionCount)
	if err != nil {
		return err
	}

	err = p.addSegmentToUsers(ctx, *segment, users)
	if err != nil {
		return err
	}
	return nil
}
