package model

import (
	"errors"
	"time"

	"github.com/jomei/notionapi"
	"google.golang.org/api/calendar/v3"
)

// GCalendarでは終日の場合、当日の00:00:00と次の日の00:00:00を設定するため、
// Dateが空の文字列でないかつ時間差が1日だけの場合に、終了時間をnilとして返す。
// 終日にも関わらず、NotionのDBの表示に00:00されてしまうのは、notionapiパッケージで、
// time.TimeをJson化する時にRFC3339のフォーマットでしか対応してないため。
func ConvertGCalendarToNotionTimeFormat(gcalendarStart, gcalendarEnd *calendar.EventDateTime) (*notionapi.Date, *notionapi.Date, error) {
	if gcalendarStart == nil || gcalendarEnd == nil {
		return nil, nil, errors.New("nil pointer error")
	}
	timeStart, err := convertGCalenarToTime(gcalendarStart)
	if err != nil {
		return nil, nil, err
	}
	timeEnd, err := convertGCalenarToTime(gcalendarEnd)
	if err != nil {
		return nil, nil, err
	}
	notionStart := notionapi.Date(timeStart)
	notionEnd := notionapi.Date(timeEnd)
	// 終日処理
	if gcalendarStart.Date != "" && isAllDay(timeStart, timeEnd) {
		return &notionStart, nil, nil
	}
	return &notionStart, &notionEnd, nil
}

// GCalendar
// 日付のみ、Date: "2000-01-01", DateTime: ""
// 時間あり、Date: "", DateTime: "2023-09-08T09:00:00+09:00"
// TODO: 日本時間以外も対応
func convertGCalenarToTime(cd *calendar.EventDateTime) (time.Time, error) {
	if cd == nil {
		return time.Now(), errors.New("nil pointer error")
	}
	target := cd.Date
	layout := time.RFC3339
	if cd.Date == "" {
		target = cd.DateTime
	} else {
		target = target + "T00:00:00+09:00"
	}
	t, err := time.Parse(layout, target)
	if err != nil {
		return time.Now(), err
	}
	// nd := notionapi.Date(t)
	return t, nil
}

func isAllDay(start, end time.Time) bool {
	sub := end.Sub(start)
	if sub.Hours()/24 == 1 {
		return true
	}
	return false
}
