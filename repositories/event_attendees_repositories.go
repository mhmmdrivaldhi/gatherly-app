package repositories

import (
	"context"
	"errors"
	"gatherly-app/models"

	"gorm.io/gorm"
)

type EventAttendeeRepository interface {
	Create(ctx context.Context, attendee *models.EventAttendee) error
	FindByUserAndEvent(ctx context.Context, userID, eventID int) (*models.EventAttendee, error)
	Update(ctx context.Context, attendee *models.EventAttendee) error
	Delete(ctx context.Context, userID, eventID int) error
	ListByEventID(ctx context.Context, eventID int) ([]*models.EventAttendee, error)
	ListByUserID(ctx context.Context, userID int) ([]*models.EventAttendee, error)
	GetFavoriteCategory(userID int) (string, error)
	// Optional methods like Exists or CountByEventID could be added here too
}

type eventAttendeeRepositoryImpl struct {
	db *gorm.DB
}

func MakeNewEventAttendeeRepository(db *gorm.DB) *eventAttendeeRepositoryImpl {
	return &eventAttendeeRepositoryImpl{db: db}
}

func (r *eventAttendeeRepositoryImpl) Create(ctx context.Context, attendee *models.EventAttendee) error {
	result := r.db.WithContext(ctx).Create(attendee)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *eventAttendeeRepositoryImpl) FindByUserAndEvent(ctx context.Context, userID, eventID int) (*models.EventAttendee, error) {
	var attendee models.EventAttendee
	result := r.db.WithContext(ctx).Where("user_id = ? AND event_id = ?", userID, eventID).First(&attendee)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &attendee, nil
}

func (r *eventAttendeeRepositoryImpl) Update(ctx context.Context, attendee *models.EventAttendee) error {
	result := r.db.WithContext(ctx).Save(attendee)
	// result := r.db.WithContext(ctx).Model(&models.EventAttendee{}).Where("user_id = ? AND event_id = ?", attendee.UserID, attendee.EventID).Updates(attendee)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *eventAttendeeRepositoryImpl) Delete(ctx context.Context, userID, eventID int) error {
	result := r.db.WithContext(ctx).Where("user_id = ? AND event_id = ?", userID, eventID).Delete(&models.EventAttendee{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *eventAttendeeRepositoryImpl) ListByEventID(ctx context.Context, eventID int) ([]*models.EventAttendee, error) {
	var attendees []*models.EventAttendee
	result := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&attendees)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []*models.EventAttendee{}, nil
		}
		return nil, result.Error
	}
	return attendees, nil
}

func (r *eventAttendeeRepositoryImpl) ListByUserID(ctx context.Context, userID int) ([]*models.EventAttendee, error) {
	var attendees []*models.EventAttendee
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&attendees)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []*models.EventAttendee{}, nil
		}
		return nil, result.Error
	}
	return attendees, nil
}

// Tambahkan method baru untuk cari kategori favorit user
func (r *eventAttendeeRepositoryImpl) GetFavoriteCategory(userID int) (string, error) {
	var category string
	err := r.db.Model(&models.EventAttendee{}).
		Select("events.category").
		Joins("JOIN events ON events.id = event_attendees.event_id").
		Where("event_attendees.user_id = ?", userID).
		Group("events.category").
		Order("COUNT(*) DESC").
		Limit(1).
		Pluck("events.category", &category).Error
	return category, err
}
