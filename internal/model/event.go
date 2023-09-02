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

func (e *EventsService) EventList(calendarID string) ([]*Event, error) {
	events, err := e.Srv.List(calendarID).Do()
	if err != nil {
		return nil, err
	}
	var ems []*Event
	for _, event := range events.Items {
		if event == nil {
			continue
		}
		ems = append(ems, &Event{Event: event})
	}
	return ems, nil
}
