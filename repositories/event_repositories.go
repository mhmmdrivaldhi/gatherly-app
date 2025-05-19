package repositories

import (
	"errors"
	"gatherly-app/models"
	"gatherly-app/models/dto"

	"gorm.io/gorm"
)

type eventsRepository struct {
	db *gorm.DB
}

type EventsRepository interface {
	CreateEvent(event *models.Event) (*models.Event, error) 
	FindEvent() ([]models.Event, error)
	FindEventByID(id int) (*models.Event, error)
	UpdateEvent(id int, updatedEvent *models.Event) (*models.Event, error)
	DeleteEvent(id int) error
	FindEventByDistance(latitude, longitude, radius float64) ([]dto.EventNearbyDistanceResponseDTO, error) 
}

func NewEventsRepository(db *gorm.DB) *eventsRepository {
	return &eventsRepository{db: db}
}

func (e *eventsRepository) CreateEvent(event *models.Event) (*models.Event, error) {
	err := e.db.Create(event).Error
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (e *eventsRepository) FindEvent() ([]models.Event, error) {
	var event []models.Event

	err := e.db.Preload("Tickets").Find(&event).Error
	if err != nil {
		return nil, err
	}

	if len(event) == 0 {
		return nil, errors.New("tidak ada acara yang ditemukan")
	}
	return event, nil
}

func (e *eventsRepository) FindEventByID(id int) (*models.Event, error) {
	var event models.Event

	err := e.db.Preload("Tickets").First(&event, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("event tidak ditemukan")
	} else if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e *eventsRepository) UpdateEvent(id int, updatedEvent *models.Event) (*models.Event, error) {
	var event models.Event

	err := e.db.First(&event, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("event tidak ditemukan")
	} else if err != nil {
		return nil, err
	}
	
	err = e.db.Model(&event).Updates(updatedEvent).Error
	if err != nil {
		return nil, err
	}
	return &event, err
}

func (e *eventsRepository) DeleteEvent(id int) error {
	err := e.db.Delete(&models.Event{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *eventsRepository) FindEventByDistance(latitude, longitude, radius float64) ([]dto.EventNearbyDistanceResponseDTO, error) {
	var results []dto.EventNearbyDistanceResponseDTO

	query := `
	SELECT * FROM (
		SELECT 
			id, name, category, description, start_date, end_date, is_paid, price, capacity,
			latitude, longitude, poster_url, status,
			CAST(6371 * acos(
				cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) +
				sin(radians(?)) * sin(radians(latitude))
			) AS FLOAT) AS distance
		FROM events
		WHERE latitude IS NOT NULL AND longitude IS NOT NULL
	) AS events_with_distance
	WHERE distance <= ?
	ORDER BY distance ASC
	`
	
	err := e.db.Raw(query, latitude, longitude, latitude, radius).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for i := range results {
		var ticket models.Ticket
		err := e.db.Where("id = ?", results[i].ID).Limit(1).First(&ticket).Error
		if err == nil {
			status := "Available"
			if ticket.Quota <=  0 || ticket.Quota >= results[i].Capacity {
				status = "Not Available"
			}

			results[i].Ticket = &dto.TicketResponseDTO{
				ID: ticket.Id,
				TicketType: ticket.TicketType,
				Price: ticket.Price,
				Quota: ticket.Quota,
				Status: status,
			}
		}
	}

	return results, nil
}