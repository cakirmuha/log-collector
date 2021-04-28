package log_collector

import "cloud.google.com/go/civil"

type CodewayLog struct {
	Type      string  `json:"type"`
	AppID     string  `json:"app_id"`
	SessionID string  `json:"session_id"`
	EventName string  `json:"event_name"`
	EventTime *int64  `json:"event_time,omitempty"`
	EventTs   *string `json:"event_ts,omitempty"`
	Page      string  `json:"page"`
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	UserID    string  `json:"user_id"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}

type TotalUser struct {
	Cnt int64 `json:"count"`
}

type DailyActiveUser struct {
	Dte civil.Date `json:"-"`
	Day string     `json:"day"`
	Cnt int64      `json:"count"`
}

type DailyAverageDuration struct {
	Dte      civil.Date `json:"-"`
	Day      string     `json:"day"`
	Duration int64      `json:"duration"`
}

type CodewayAnalytics struct {
	TotalUser             TotalUser              `json:"total_user"`
	DailyActiveUsers      []DailyActiveUser      `json:"daily_active_users"`
	DailyAverageDurations []DailyAverageDuration `json:"daily_average_durations"`
}
