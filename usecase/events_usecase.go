package usecase

import (
	"errors"
	"fmt"
	"gatherly-app/models"
	"gatherly-app/models/dto"
	"gatherly-app/repositories"
	"gatherly-app/utils"

	"time"
)

type eventsUsecase struct {
	repo         repositories.EventsRepository
	attendeeRepo repositories.EventAttendeeRepository
}

type EventsUsecase interface {
	CreateEvent(request dto.CreateEventRequestDTO) (*models.Event, error)
	GetAllEvent() ([]dto.EventResponseDTO, error)
	GetEventByID(id int) (*dto.EventResponseDTO, error)
	UpdateEvent(id int, request dto.UpdateEventRequestDTO) (*models.Event, error)
	DeleteEvent(id int) error
	GetEventByDistance(latitude, longitude, radius float64) ([]dto.EventNearbyDistanceResponseDTO, error)
}

// Update NewEventUsecase
func NewEventUsecase(
	repo repositories.EventsRepository,
	attendeeRepo repositories.EventAttendeeRepository, // Sesuai dengan nama di server.go
) EventsUsecase {
	return &eventsUsecase{
		repo:         repo,
		attendeeRepo: attendeeRepo,
	}
}

func (uc *eventsUsecase) CreateEvent(request dto.CreateEventRequestDTO) (*models.Event, error) {
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, fmt.Errorf("format harus YYYY-MM-DD: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return nil, fmt.Errorf("format harus YYYY-MM-DD: %w", err)
	}

	coordinate, err := utils.GetCoordinatesFromAddress(request.Address)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan koordinat: %w", err)
	}

	events := &models.Event{
		Name: request.Name,
		Category: request.Category,
		Description: request.Description,
		StartDate: startDate,
		EndDate: endDate,
		IsPaid: request.IsPaid,
		Capacity: request.Capacity,
		Latitude: coordinate.Latitude,
		Longitude: coordinate.Longitude,
		PosterURL: request.PosterURL,
		Status: request.Status,
	}

	create, err := uc.repo.CreateEvent(events)
	if err != nil {
		return nil, err
	}
	return create, nil
}

func (uc *eventsUsecase) GetAllEvent() ([]dto.EventResponseDTO, error) {
	events, err := uc.repo.FindEvent()
	if err != nil {
		return nil, err
	}

	var response []dto.EventResponseDTO
	now := time.Now()

	for _, event := range events {
		startTime := event.StartDate
		endTime := event.EndDate

		status := "Unknown"
		if startTime.After(now) {
			status = "Up Coming"
		} else if now.After(startTime) && now.Before(endTime) {
			status = "On Going"
		} else if now.After(endTime) {
			status = "Ended"
		}

		var ticketResponse *dto.TicketResponseDTO
		if len(event.Tickets) > 0 {
			ticket := event.Tickets[0]
			ticketStatus := "Available"
			if ticket.Quota <= 0 || ticket.Quota >= event.Capacity {
				ticketStatus = "Not Available"
			}
			ticketResponse = &dto.TicketResponseDTO{
				ID: ticket.Id,
				TicketType: ticket.TicketType,
				Price: ticket.Price,
				Quota: ticket.Quota,
				Status: ticketStatus,
			}
		}
		eventResponse := dto.EventResponseDTO{
			Name: event.Name,
			Category: event.Category,
			Description: event.Description,
			StartDate: event.StartDate.Format(time.RFC3339),
			EndDate: event.EndDate.Format(time.RFC3339),
			IsPaid:      event.IsPaid,
			Ticket:      ticketResponse,
			Capacity:    event.Capacity,
			Latitude:    event.Latitude,
			Longitude:   event.Longitude,
			PosterURL:   event.PosterURL,
			Status:      status,
		}
		response = append(response, eventResponse)
	}
	return response, nil
}

func (uc *eventsUsecase) GetEventByID(id int) (*dto.EventResponseDTO, error) {
	event, err := uc.repo.FindEventByID(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	status := "Unknown"
	if event.StartDate.After(now) {
		status = "Upcoming"
	} else if now.After(event.StartDate) && now.Before(event.EndDate) {
		status = "Ongoing"
	} else if now.After(event.EndDate) {
		status = "Ended"
	}

	var ticketResponse *dto.TicketResponseDTO
	if len(event.Tickets) > 0 {
		ticket := event.Tickets[0]
		ticketStatus := "Available"
		if ticket.Quota <= 0 || ticket.Quota >= event.Capacity {
			ticketStatus = "Not Available"
		}
		ticketResponse = &dto.TicketResponseDTO{
			ID:         ticket.Id,
			TicketType: ticket.TicketType,
			Price:      ticket.Price,
			Quota:      event.Capacity,
			Status:     ticketStatus,
		}
	}

	response := &dto.EventResponseDTO{
		Name:        event.Name,
		Category:    event.Category,
		Description: event.Description,
		StartDate:   event.StartDate.Format(time.RFC3339),
		EndDate:     event.EndDate.Format(time.RFC3339),
		IsPaid:      event.IsPaid,
		Ticket:      ticketResponse,
		Capacity:    event.Capacity,
		Latitude:    event.Latitude,
		Longitude:   event.Longitude,
		PosterURL:   event.PosterURL,
		Status:      status,
	}
	return response, nil
}

func (uc *eventsUsecase) UpdateEvent(id int, request dto.UpdateEventRequestDTO) (*models.Event, error) {
	var startDate, endDate time.Time
	var err error

	if request.StartDate != nil {
		startDate, err = time.Parse("2006-01-02", *request.StartDate)
		if err != nil {
			return nil, fmt.Errorf("format harus YYYY-MM-DD: %w", err)
		}
	}

	if request.EndDate != nil {
		endDate, err = time.Parse("2006-01-02", *request.EndDate)
		if err != nil {
			return nil, fmt.Errorf("format harus YYYY-MM-DD: %w", err)
		}
	}

	isExist, err := uc.repo.FindEventByID(id)
	if err != nil {
		return nil, err
	}

	if request.Name != nil {
		isExist.Name = *request.Name
	}
	if request.Category != nil {
		isExist.Category = *request.Category
	}
	if request.Description != nil {
		isExist.Description = *request.Description
	}
	if request.StartDate != nil {
		isExist.StartDate = startDate
	}
	if request.EndDate != nil {
		isExist.EndDate = endDate
	}
	if request.IsPaid != nil {
		isExist.IsPaid = *request.IsPaid
	}
	if request.Capacity != nil {
		isExist.Capacity = *request.Capacity
	}

	if request.Address != nil {
		coordinate, err := utils.GetCoordinatesFromAddress(*request.Address)
		if err == nil {
			isExist.Latitude = coordinate.Latitude
			isExist.Longitude = coordinate.Longitude
		}
	}

	if request.PosterURL != nil {
		isExist.PosterURL = *request.PosterURL
	}

	now := time.Now()
	startTime := isExist.StartDate
	endTime := isExist.EndDate

	if startTime.After(now) {
		isExist.Status = "Upcoming"
	} else if now.After(startTime) && now.Before(endTime) {
		isExist.Status = "Ongoing"
	} else if now.After(endTime) {
		isExist.Status = "Ended"
	}

	updatedEvent, err := uc.repo.UpdateEvent(id, isExist)
	if err != nil {
		return nil, err
	}
	return updatedEvent, nil
}

func (uc *eventsUsecase) DeleteEvent(id int) error {
	_, err := uc.repo.FindEventByID(id)
	if err != nil {

		return errors.New("event tidak ditemukan atau terjadi kesalahan saat mencari: " + err.Error())
	}

	err = uc.repo.DeleteEvent(id)
	if err != nil {
		return errors.New("data berhasil dihapus")
	}

	return nil
}

func (uc *eventsUsecase) GetEventByDistance(latitude, longitude, radius float64) ([]dto.EventNearbyDistanceResponseDTO, error) {
	events, err := uc.repo.FindEventByDistance(latitude, longitude, radius)
	if err != nil {
		return nil, err
	}

	var results []dto.EventNearbyDistanceResponseDTO
	now := time.Now()

	for _, event := range events {
		status := "Unknown"
		if event.StartDate.After(now) {
			status = "Up coming"
		} else if now.After(event.StartDate) && now.Before(event.EndDate) {
			status = "On Going"
		} else if now.After(event.EndDate) {
			status = "Ended"
		}

		var ticketDTO *dto.TicketResponseDTO
		if event.Ticket != nil {
			ticketStatus := "available"
			if event.Ticket.Quota <= 0 || event.Ticket.Quota >= event.Capacity {
				ticketStatus = "sold out"
			}

			ticketDTO = &dto.TicketResponseDTO{
				ID:         event.Ticket.ID,
				TicketType: event.Ticket.TicketType,
				Price:      event.Ticket.Price,
				Quota:      event.Ticket.Quota,
				Status:     ticketStatus,
			}
		}

		results = append(results, dto.EventNearbyDistanceResponseDTO{
			Name:        event.Name,
			Category:    event.Category,
			Description: event.Description,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			IsPaid:      event.IsPaid,
			Capacity:    event.Capacity,
			Latitude:    event.Latitude,
			Longitude:   event.Longitude,
			Distance:    event.Distance,
			PosterURL:   event.PosterURL,
			Status:      status,
			Ticket:      ticketDTO, // tambahkan ke DTO
		})
	}
	return results, nil
}
