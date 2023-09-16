package request

type ExportRequest struct {
	DBName      string   `json:"db_name"`
	PageID      string   `json:"page_id"`
	CalendarIDs []string `json:"calendar_ids"`
}
