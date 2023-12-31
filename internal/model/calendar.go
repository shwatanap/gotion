package model

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Calendar struct {
	Calendar *calendar.CalendarListEntry
}

type CalendarsService struct {
	Srv *calendar.Service
}

func NewCalendarService(ctx context.Context, token *oauth2.Token) (*CalendarsService, error) {
	o := NewGoogleOAuth()
	srv, err := calendar.NewService(ctx, option.WithTokenSource(o.Config.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}
	return &CalendarsService{
		Srv: srv,
	}, nil
}

func (cs *CalendarsService) CalendarList() ([]*Calendar, error) {
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

func (cs *CalendarsService) GetCalendar(calendarID string) (*Calendar, error) {
	calendarListEntry, err := cs.Srv.CalendarList.Get(calendarID).Do()
	if err != nil {
		return nil, err
	}
	if calendarListEntry == nil {
		return nil, fmt.Errorf("calendar not found")
	}
	calendar := &Calendar{
		Calendar: calendarListEntry,
	}
	return calendar, nil
}
