package model

import (
	"context"
	"log"
	"strconv"

	"github.com/jomei/notionapi"
	"golang.org/x/oauth2"
)

func GCalendarExport(ctx context.Context, token *oauth2.Token, notionAPIkey, pageID, dbTitle string) {
	cs, _ := NewCalendarService(ctx, token)
	es := NewEventService(cs.Srv.Events)
	calendars, _ := cs.CalendarList()
	calendarOptions := make([]notionapi.Option, len(calendars))
	for i, c := range calendars {
		color := GCalendaToNotionColorMap[c.Calendar.ColorId]
		summary := c.Calendar.Summary
		option := notionapi.Option{
			// IDがNotionで自動で書き換えられている
			// TODO: 対処法を考える
			ID:    notionapi.PropertyID(strconv.Itoa(i)),
			Name:  summary,
			Color: color,
		}
		calendarOptions[i] = option
	}

	nc := NewNotionClient(notionAPIkey)
	db, err := nc.CreateDatabase(ctx, CreateDatabaseRequest{
		PageID:          pageID,
		Title:           dbTitle,
		CalendarOptions: calendarOptions,
	})
	if err != nil {
		log.Printf("error from new notion client😡: %v", err)
	}

	// ここからイベント作成処理
	for i, c := range calendars {
		events, err := es.EventList(c.Calendar.Id)
		if err != nil {
			log.Printf("error from event list😡: %v", err)
			return
		}
		for _, e := range events {
			start, end, err := ConvertGCalendarToNotionTimeFormat(e.Event.Start, e.Event.End)
			if err != nil {
				log.Fatalln(err)
				return
			}
			if err := nc.AddEvent(ctx, AddEventRequest{
				DatabaseID:       notionapi.DatabaseID(db.ID),
				Title:            e.Event.Summary,
				SelectedCalendar: db.Properties["calendar"].(*notionapi.SelectPropertyConfig).Select.Options[i],
				DateStart:        start,
				DateEnd:          end,
			}); err != nil {
				log.Printf("error from add event😡: %v", err)
				return
			}
		}
	}
}
