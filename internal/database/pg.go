package database

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
	"math"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true, Logger: logger.Default.LogMode(logger.Info)}) //todo: use custom logger based on logrus
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

// getUsersCount returns count of user.
func (p *pg) getUsersCount(ctx context.Context) (int64, error) {
	var count int64
	result := p.conn.WithContext(ctx).Model(&model.User{}).Count(&count)
	if result.Error != nil {
		return count, result.Error
	}
	return count, nil
}

// createSegment returns new saved model.Segment.
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

// deleteSegment soft deletes segment by slug.
func (p *pg) deleteSegment(ctx context.Context, tx *gorm.DB, segment *model.Segment) error {
	result := tx.WithContext(ctx).Clauses(clause.Returning{}).Delete(segment, "slug = ?", segment.Slug)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrDB, result.Error)
	} else if result.RowsAffected < 1 {
		return fmt.Errorf("segment with slug (%s) is %w for delete", segment.Slug, ErrNotFound)
	}
	return nil
}

// deleteSegment soft deletes user segment relation by segment id.
func (p *pg) deleteSegmentFromUsers(ctx context.Context, tx *gorm.DB, segment *model.Segment) error {
	result := tx.WithContext(ctx).Delete(&model.UserSegment{}, "segment_id = ?", segment.ID)
	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return nil
}

// DeleteSegment soft deletes segment by slug from its table and user_segments.
func (p *pg) DeleteSegment(ctx context.Context, segment *model.Segment) error {
	err := p.conn.Transaction(func(tx *gorm.DB) error {
		err := p.deleteSegment(ctx, tx, segment)
		if err != nil {
			return err
		}

		err = p.deleteSegmentFromUsers(ctx, tx, segment)
		if err != nil {
			return err
		}
		return nil
	})

	return err
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
func (p *pg) GetSegments(ctx context.Context, slugs []model.Slug) ([]*model.Segment, error) {
	return p.getSegments(ctx, p.conn, slugs)
}

// getSegments returns segments by given slugs.
func (p *pg) getSegments(ctx context.Context, tx *gorm.DB, slugs []model.Slug) ([]*model.Segment, error) {
	var segments []*model.Segment
	result := tx.WithContext(ctx).Find(&segments, "slug IN ?", slugs)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrDB, result.Error)
	} else if len(segments) != len(slugs) {
		return nil, fmt.Errorf("some segments %w", ErrNotFound)
	}
	return segments, nil
}

// CreateDeleteUserSegments insert and delete relation by specified segments (represented by slugs) for given user.
func (p *pg) CreateDeleteUserSegments(ctx context.Context, user *model.User, toInSegments []model.Slug, toDelSegments []model.Slug) error {
	var insertSegments []*model.Segment
	var deleteSegments []*model.Segment
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

		if len(insertSegments) > 0 {
			txErr = p.insertUserSegments(ctx, tx, user, insertSegments)
			if txErr != nil {
				return txErr
			}
		}

		if len(deleteSegments) > 0 {
			txErr = p.deleteUserSegments(ctx, tx, user, deleteSegments)
			if txErr != nil {
				return txErr
			}
		}

		return nil
	})

	return err
}

// insertUserSegments connects user with given segments.
func (p *pg) insertUserSegments(ctx context.Context, tx *gorm.DB, user *model.User, segments []*model.Segment) error {
	log.Info(user.ID)
	err := tx.WithContext(ctx).Model(user).Omit("Segments.*").Association("Segments").Append(segments)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDB, err)
	}
	return nil
}

// deleteUserSegments soft deletes user relation to given segments,
// just updates deleted_at column.
func (p *pg) deleteUserSegments(ctx context.Context, tx *gorm.DB, user *model.User, segments []*model.Segment) error {
	err := tx.WithContext(ctx).Model(user).Association("Segments").Delete(segments)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDB, err)
	}
	return nil
}

// GetUserWithActiveSegments returns user with related segments.
func (p *pg) GetUserWithActiveSegments(ctx context.Context, user *model.User) (*model.User, error) {
	result := p.conn.WithContext(ctx).Preload("Segments").First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user with ID (%d) %w", user.ID, ErrNotFound)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}
	return user, nil
}

// GetUserSegmentsHistory returns user relation to segments for the specified time interval.
func (p *pg) GetUserSegmentsHistory(ctx context.Context, user *model.User, from time.Time, to time.Time) ([]model.UserSegment, error) {
	var userSegments []model.UserSegment

	result := p.conn.WithContext(ctx).Unscoped().Preload("Segment").Model(&userSegments).
		Find(&userSegments, "user_segments.user_id = ? AND "+
			"((user_segments.created_at >= ? AND user_segments.created_at < ?)"+
			" OR (user_segments.deleted_at = null "+
			" OR (user_segments.deleted_at >= ? AND user_segments.deleted_at < ?)))"+
			" ORDER BY user_segments.created_at ASC", user.ID, from, to, from, to)
	if result.Error != nil {
		log.Error(result)
		return nil, fmt.Errorf("%w: %v", ErrDB, result.Error)
	}

	return userSegments, nil
}

// getRandomUsers returns random users rows.
func (p *pg) getRandomUsers(ctx context.Context, count int) ([]*model.User, error) {
	var users []*model.User
	result := p.conn.WithContext(ctx).
		Order("random()").
		Limit(count).
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// addSegmentToUsers adds relation for given segment and users.
func (p *pg) addSegmentToUsers(ctx context.Context, segment model.Segment, users []*model.User) error {
	for _, user := range users {
		if err := p.insertUserSegments(ctx, p.conn, user, []*model.Segment{&segment}); err != nil {
			return err
		}
	}
	return nil
}

// AddSegmentToRandomUsers adds relation for given segment to random users depending on selection percent.
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

// GetUser returns user with given user.ID.
func (p *pg) GetUser(ctx context.Context, user *model.User) (*model.User, error) {
	result := p.conn.WithContext(ctx).First(user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user with id (%d) %w", user.ID, ErrNotFound)
	}
	return user, nil
}
