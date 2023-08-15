package model

import "google.golang.org/api/calendar/v3"

type Event struct {
	Event *calendar.Event
}

type EventsService struct {
	Srv        *calendar.EventsService
	CalendarId string
}

func NewEventService(srv *calendar.EventsService, calendarId string) *EventsService {
	return &EventsService{
		Srv:        srv,
		CalendarId: calendarId,
	}
}

func (e *EventsService) List(timeMin string) ([]*Event, error) {
	events, err := e.Srv.List(e.CalendarId).TimeMin(timeMin).Do()
	if err != nil {
		return nil, err
	}
	ems := make([]*Event, len(events.Items))
	for i, event := range events.Items {
		ems[i] = &Event{Event: event}
	}
	return ems, nil
}
