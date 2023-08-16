package model

import "google.golang.org/api/calendar/v3"

type Event struct {
	Event *calendar.Event
}

type EventsService struct {
	Srv *calendar.EventsService
}

func NewEventService(srv *calendar.EventsService) *EventsService {
	return &EventsService{
		Srv: srv,
	}
}

func (e *EventsService) List(calendarId string, timeMin string) ([]*Event, error) {
	events, err := e.Srv.List(calendarId).TimeMin(timeMin).Do()
	if err != nil {
		return nil, err
	}
	ems := make([]*Event, len(events.Items))
	for i, event := range events.Items {
		ems[i] = &Event{Event: event}
	}
	return ems, nil
}
