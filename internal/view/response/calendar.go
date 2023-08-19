package response

type CalendarResponse struct {
	ID         string `json:"id"`
	Summary    string `json:"summary"`
	ColorID    string `json:"color_id"`
	AccessRole string `json:"access_role"`
}

type CalendarListResponse struct {
	Calendars []CalendarResponse `json:"calendars"`
}
