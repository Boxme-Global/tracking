package omisocial

type Response struct {
	Message string      `json:"message"`
	Error   bool        `json:"error"`
	Data    interface{} `json:"data"`
}

type ResponseVisitors struct {
	Message string         `json:"message"`
	Error   bool           `json:"error"`
	Data    []VisitorStats `json:"data"`
}

type ResponsePages struct {
	Message    string      `json:"message"`
	Error      bool        `json:"error"`
	TotalPages int         `json:"total_pages"`
	Count      int         `json:"count"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Data       []PageStats `json:"data"`
}

type ResponseTotalVisitors struct {
	Message string             `json:"message"`
	Error   bool               `json:"error"`
	Data    *TotalVisitorStats `json:"data"`
}

type ResponseReferrers struct {
	Message    string          `json:"message"`
	Error      bool            `json:"error"`
	TotalPages int             `json:"total_pages"`
	Count      int             `json:"count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	Data       []ReferrerStats `json:"data"`
}

type ResponsePlatformVisitors struct {
	Message string                 `json:"message"`
	Error   bool                   `json:"error"`
	Data    []PlatformVisitorStats `json:"data"`
}

type ResponseEvents struct {
	Message string       `json:"message"`
	Error   bool         `json:"error"`
	Data    []EventStats `json:"data"`
}

type ResponseUTMSources struct {
	Message    string           `json:"message"`
	Error      bool             `json:"error"`
	TotalPages int              `json:"total_pages"`
	Count      int              `json:"count"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	Data       []UTMSourceStats `json:"data"`
}

type ResponseOTMSources struct {
	Message    string           `json:"message"`
	Error      bool             `json:"error"`
	TotalPages int              `json:"total_pages"`
	Count      int              `json:"count"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	Data       []OTMSourceStats `json:"data"`
}

type GroupEvent struct {
	Name     string `json:"name"`
	Visitors int    `json:"visitors"`
	Views    int    `json:"views"`
	Sessions int    `json:"sessions"`
}

type GroupEvents struct {
	Period string       `json:"period"`
	Events []GroupEvent `json:"events"`
}

type ResponseGroupEvents struct {
	Message string        `json:"message"`
	Error   bool          `json:"error"`
	Data    []GroupEvents `json:"data"`
}

type OverTime struct {
	Period   string       `json:"period"`
	Visitors int          `json:"visitors"`
	Events   []GroupEvent `json:"events"`
}

type ResponseOverTime struct {
	Message string     `json:"message"`
	Error   bool       `json:"error"`
	Data    []OverTime `json:"data"`
}
