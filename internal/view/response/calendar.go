package response

type CalendarResponse struct {
	Id         string `json:"id"`
	Summary    string `json:"summary"`
	ColorId    string `json:"color_id"`
	AccessRole string `json:"access_role"`
}

type CalendarListResponse struct {
	Calendars []CalendarResponse `json:"calendars"`
}
