package model

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jomei/notionapi"
	"golang.org/x/oauth2"
)

func GCalendarExport(ctx context.Context, token *oauth2.Token, notionAPIkey, pageID, dbTitle string, calendarIDs []string) (*notionapi.Database, error) {
	cs, _ := NewCalendarService(ctx, token)
	es := NewEventService(cs.Srv.Events)

	// calendarIDsからカレンダーを取得
	calendars := make([]*Calendar, len(calendarIDs))
	for i, calendarID := range calendarIDs {
		calendar, err := cs.GetCalendar(calendarID)
		if err != nil {
			return nil, err
		}
		if calendar == nil {
			return nil, fmt.Errorf("calendar not found")
		}
		calendars[i] = calendar
	}

	// カレンダーのオプションを作成
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

	// データベースの作成
	nc := NewNotionClient(notionAPIkey)
	db, err := nc.CreateDatabase(ctx, CreateDatabaseRequest{
		PageID:          pageID,
		Title:           dbTitle,
		CalendarOptions: calendarOptions,
	})
	if err != nil {
		return nil, err
	}

	// ここからイベント作成処理
	for i, c := range calendars {
		events, err := es.EventList(c.Calendar.Id)
		if err != nil {
			return nil, err
		}
		for _, e := range events {
			log.Println(e.Event.Id, e.Event.Start, e.Event.Summary)
			start, end, err := ConvertGCalendarToNotionTimeFormat(e.Event.Start, e.Event.End)
			if err != nil {
				return nil, err
			}
			if err := nc.AddEvent(ctx, AddEventRequest{
				DatabaseID:       notionapi.DatabaseID(db.ID),
				Title:            e.Event.Summary,
				SelectedCalendar: db.Properties["calendar"].(*notionapi.SelectPropertyConfig).Select.Options[i],
				DateStart:        start,
				DateEnd:          end,
			}); err != nil {
				return nil, err
			}
		}
	}
	return db, nil
}
