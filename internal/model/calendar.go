package model

import "google.golang.org/api/calendar/v3"

type Calendar struct {
	Calendar *calendar.CalendarListEntry
}

type CalendarsService struct {
	Srv *calendar.Service
}

func NewCalendarService(srv *calendar.Service) *CalendarsService {
	return &CalendarsService{
		Srv: srv,
	}
}

func (cs *CalendarsService) List() ([]*Calendar, error) {
	calendarList, err := cs.Srv.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}
	cms := make([]*Calendar, len(calendarList.Items))
	for i, calendarListEntry := range calendarList.Items {
		cms[i] = &Calendar{Calendar: calendarListEntry}
	}
	return cms, nil
}
