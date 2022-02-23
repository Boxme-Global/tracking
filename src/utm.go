package omisocial

import (
	"net/http"
	"strings"
)

type utmParams struct {
	source   string
	medium   string
	campaign string
	content  string
	term     string
}

type otmParams struct {
	source   string
	medium   string
	campaign string
	position string
}

func getUTMParams(r *http.Request) utmParams {
	query := r.URL.Query()
	return utmParams{
		source:   strings.TrimSpace(query.Get("utm_source")),
		medium:   strings.TrimSpace(query.Get("utm_medium")),
		campaign: strings.TrimSpace(query.Get("utm_campaign")),
		content:  strings.TrimSpace(query.Get("utm_content")),
		term:     strings.TrimSpace(query.Get("utm_term")),
	}
}

func getOTMParams(r *http.Request) otmParams {
	query := r.URL.Query()
	return otmParams{
		source:   strings.TrimSpace(query.Get("otm_source")),
		medium:   strings.TrimSpace(query.Get("otm_medium")),
		campaign: strings.TrimSpace(query.Get("otm_campaign")),
		position: strings.TrimSpace(query.Get("otm_position")),
	}
}
