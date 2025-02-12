package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"

	omisocial "github.com/Boxme-Global/tracking/src"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	omisocial.SetFingerprintKeys(42, 123)
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("ENV")
	log.Println(env)
	var (
		migrate_uri string
		client_uri  string
	)

	if env == "production" {
		clickhouse_host := os.Getenv("CLICKHOUSE_HOST")
		clickhouse_port := os.Getenv("CLICKHOUSE_PORT")

		q := make(url.Values)
		q.Set("username", os.Getenv("CLICKHOUSE_USERNAME"))
		q.Set("password", os.Getenv("CLICKHOUSE_PASSWORD"))
		q.Set("database", os.Getenv("CLICKHOUSE_DATABASE"))

		client_uri = (&url.URL{
			Scheme:   "tcp",
			Host:     clickhouse_host + ":" + clickhouse_port,
			RawQuery: q.Encode(),
		}).String()

		q.Set("x-multi-statement", "true")

		migrate_uri = (&url.URL{
			Scheme:   "clickhouse",
			Host:     clickhouse_host + ":" + clickhouse_port,
			RawQuery: q.Encode(),
		}).String()

	} else {
		migrate_uri = "clickhouse://10.148.0.23:9000?x-multi-statement=true"
		client_uri = "tcp://10.148.0.23:9000"

	}

	// Set the key for SipHash. This should be called on startup (before generating the first fingerprint) and is NOT concurrency save.

	// Migrate the database.
	omisocial.Migrate(migrate_uri)

	// Create a new ClickHouse client to save hits.
	store, _ := omisocial.NewClient(client_uri, nil)
	// db, err := sql.Open("clickhouse", client_uri)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Set up a default tracker with a salt.
	// This will buffer and store hits and generate sessions by default.
	tracker := omisocial.NewTracker(store, "BuS7BsvURhatRPqr", nil)

	// Create a handler to serve traffic.
	// We prevent tracking resources by checking the path. So a file on /my-file.txt won't create a new hit
	// but all page calls will be tracked.
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventName := r.URL.Query().Get("event_name")
		pageloadEvents := []string{"pageload", "pageclose"}

		if eventName != "" && omisocial.Contains(pageloadEvents, eventName) {
			tracker.Hit(r, omisocial.HitOptionsFromRequest(r))
		} else {
			metaJson := r.URL.Query().Get("event_data")
			var eventData map[string]interface{}
			json.Unmarshal([]byte(metaJson), &eventData)
			options := omisocial.EventOptions{
				Name:     eventName,
				Duration: 0, // optional field to save a duration, this will be used to calculate an average time when using the analyzer
				Meta:     eventData,
			}

			tracker.Event(r, options, omisocial.HitOptionsFromRequest(r))
		}
		w.Write([]byte("hi"))
	}))

	http.Handle("/report/visitors", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		analyzer := omisocial.NewAnalyzer(store)

		from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
		site_id, _ := strconv.ParseInt(r.URL.Query().Get("site_id"), 10, 64)
		group_by := r.URL.Query().Get("group_by")

		groups := []string{"day", "week", "month"}
		if from == 0 || to == 0 || site_id == 0 || from > to || (group_by != "" && !omisocial.Contains(groups, group_by)) {
			jData, _ := json.Marshal(&omisocial.Response{
				Message: "Invalid input data",
				Error:   true,
				Data:    nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		visitors, _ := analyzer.Visitors(
			&omisocial.Filter{
				From:     time.Unix(from, 0),
				To:       time.Unix(to, 0),
				ClientID: site_id,
			},
			group_by,
		)

		jData, _ := json.Marshal(&omisocial.ResponseVisitors{
			Message: "",
			Error:   false,
			Data:    visitors,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))

	http.Handle("/report/total-visitors", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		analyzer := omisocial.NewAnalyzer(store)

		from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
		site_id, _ := strconv.ParseInt(r.URL.Query().Get("site_id"), 10, 64)

		if from == 0 || to == 0 || site_id == 0 || from > to {
			jData, _ := json.Marshal(&omisocial.Response{
				Message: "Invalid input data",
				Error:   true,
				Data:    nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		growth, _ := analyzer.TotalVisitors(&omisocial.Filter{
			From:     time.Unix(from, 0),
			To:       time.Unix(to, 0),
			ClientID: site_id,
		})

		jData, _ := json.Marshal(&omisocial.ResponseTotalVisitors{
			Message: "",
			Error:   false,
			Data:    growth,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))

	http.Handle("/report/platforms", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		analyzer := omisocial.NewAnalyzer(store)

		from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
		site_id, _ := strconv.ParseInt(r.URL.Query().Get("site_id"), 10, 64)

		if from == 0 || to == 0 || site_id == 0 || from > to {
			jData, _ := json.Marshal(&omisocial.Response{
				Message: "Invalid input data",
				Error:   true,
				Data:    nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		platforms, _ := analyzer.PlatformVisitors(&omisocial.Filter{
			From:     time.Unix(from, 0),
			To:       time.Unix(to, 0),
			ClientID: site_id,
		})

		jData, _ := json.Marshal(&omisocial.ResponsePlatformVisitors{
			Message: "",
			Error:   false,
			Data:    platforms,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))

	http.Handle("/report/pages", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		analyzer := omisocial.NewAnalyzer(store)

		from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
		site_id, _ := strconv.ParseInt(r.URL.Query().Get("site_id"), 10, 64)
		path_pattern := r.URL.Query().Get("path_pattern")
		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		page_size, _ := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)

		if from == 0 || to == 0 || site_id == 0 || from > to {
			jData, _ := json.Marshal(&omisocial.Response{
				Message: "Invalid input data",
				Error:   true,
				Data:    nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		if page < 1 {
			page = 1
		}

		if page_size <= 0 {
			page_size = 25
		}

		count, _ := analyzer.PageCount(&omisocial.Filter{
			From:        time.Unix(from, 0),
			To:          time.Unix(to, 0),
			ClientID:    site_id,
			PathPattern: path_pattern,
		})

		if count == 0 {
			jData, _ := json.Marshal(&omisocial.ResponsePages{
				Message:    "No data",
				Error:      false,
				TotalPages: 0,
				Count:      count,
				Page:       int(page),
				PageSize:   int(page_size),
				Data:       nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		var total_pages = int64(math.Ceil(float64(count) / float64(page_size)))

		if page > total_pages {
			page = total_pages
		}

		offset := (int(page) - 1) * int(page_size)

		pages, _ := analyzer.Pages(&omisocial.Filter{
			From:        time.Unix(from, 0),
			To:          time.Unix(to, 0),
			ClientID:    site_id,
			PathPattern: path_pattern,
			Limit:       int(page_size),
			Offset:      int(offset),
		})

		jData, _ := json.Marshal(&omisocial.ResponsePages{
			Message:    "",
			Error:      false,
			TotalPages: 0,
			Count:      count,
			Page:       int(page),
			PageSize:   int(page_size),
			Data:       pages,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))

	http.Handle("/report/referrers", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		analyzer := omisocial.NewAnalyzer(store)

		from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
		site_id, _ := strconv.ParseInt(r.URL.Query().Get("site_id"), 10, 64)
		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		page_size, _ := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)

		if from == 0 || to == 0 || site_id == 0 || from > to {
			jData, _ := json.Marshal(&omisocial.Response{
				Message: "Invalid input data",
				Error:   true,
				Data:    nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		if page < 1 {
			page = 1
		}

		if page_size <= 0 {
			page_size = 25
		}

		count, _ := analyzer.ReferrerCount(&omisocial.Filter{
			From:     time.Unix(from, 0),
			To:       time.Unix(to, 0),
			ClientID: site_id,
		})

		if count == 0 {
			jData, _ := json.Marshal(&omisocial.ResponseReferrers{
				Message:    "No data",
				Error:      false,
				TotalPages: 0,
				Count:      count,
				Page:       int(page),
				PageSize:   int(page_size),
				Data:       nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		var total_pages = int64(math.Ceil(float64(count) / float64(page_size)))

		if page > total_pages {
			page = total_pages
		}

		offset := (int(page) - 1) * int(page_size)

		referrers, _ := analyzer.Referrer(&omisocial.Filter{
			From:     time.Unix(from, 0),
			To:       time.Unix(to, 0),
			ClientID: site_id,
			Limit:    int(page_size),
			Offset:   int(offset),
		})

		jData, _ := json.Marshal(&omisocial.ResponseReferrers{
			Message:    "",
			Error:      false,
			TotalPages: 0,
			Count:      count,
			Page:       int(page),
			PageSize:   int(page_size),
			Data:       referrers,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))

	http.Handle("/report/events", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		analyzer := omisocial.NewAnalyzer(store)

		from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
		site_id, _ := strconv.ParseInt(r.URL.Query().Get("site_id"), 10, 64)

		if from == 0 || to == 0 || site_id == 0 || from > to {
			jData, _ := json.Marshal(&omisocial.Response{
				Message: "Invalid input data",
				Error:   true,
				Data:    nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		events, _ := analyzer.Events(&omisocial.Filter{
			From:     time.Unix(from, 0),
			To:       time.Unix(to, 0),
			ClientID: site_id,
		})

		jData, _ := json.Marshal(&omisocial.ResponseEvents{
			Message: "",
			Error:   false,
			Data:    events,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))

	http.Handle("/report/group-events", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		analyzer := omisocial.NewAnalyzer(store)

		from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
		site_id, _ := strconv.ParseInt(r.URL.Query().Get("site_id"), 10, 64)
		group_by := r.URL.Query().Get("group_by")

		groups := []string{"day", "week", "month"}
		if from == 0 || to == 0 || site_id == 0 || from > to || (group_by != "" && !omisocial.Contains(groups, group_by)) {
			jData, _ := json.Marshal(&omisocial.Response{
				Message: "Invalid input data",
				Error:   true,
				Data:    nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		events, _ := analyzer.GroupEvents(
			&omisocial.Filter{
				From:     time.Unix(from, 0),
				To:       time.Unix(to, 0),
				ClientID: site_id,
			},
			group_by,
		)

		group_events := map[string]omisocial.GroupEvents{}
		for _, event := range events {
			if val, ok := group_events[event.Period]; ok {
				val.Events = append(
					val.Events,
					omisocial.GroupEvent{
						Name:     event.Name,
						Visitors: event.Visitors,
						Views:    event.Views,
						Sessions: event.Sessions,
					},
				)
				group_events[val.Period] = val
			} else {
				group_event := omisocial.GroupEvents{
					Period: event.Period,
					Events: []omisocial.GroupEvent{
						{
							Name:     event.Name,
							Visitors: event.Visitors,
							Views:    event.Views,
							Sessions: event.Sessions,
						},
					},
				}
				group_events[event.Period] = group_event
			}
		}

		keys := []string{}
		for key := range group_events {
			keys = append(keys, key)

		}
		sort.Strings(keys)

		data := []omisocial.GroupEvents{}
		for _, key := range keys {
			data = append(data, group_events[key])
		}

		jData, _ := json.Marshal(&omisocial.ResponseGroupEvents{
			Message: "",
			Error:   false,
			Data:    data,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))

	http.Handle("/report/over-time/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		analyzer := omisocial.NewAnalyzer(store)

		from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
		site_id, _ := strconv.ParseInt(r.URL.Query().Get("site_id"), 10, 64)
		group_by := r.URL.Query().Get("group_by")

		groups := []string{"day", "week", "month"}
		if from == 0 || to == 0 || site_id == 0 || from > to || (group_by != "" && !omisocial.Contains(groups, group_by)) {
			jData, _ := json.Marshal(&omisocial.Response{
				Message: "Invalid input data",
				Error:   true,
				Data:    nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		visitors, _ := analyzer.Visitors(
			&omisocial.Filter{
				From:     time.Unix(from, 0),
				To:       time.Unix(to, 0),
				ClientID: site_id,
			},
			group_by,
		)

		events, _ := analyzer.GroupEvents(
			&omisocial.Filter{
				From:     time.Unix(from, 0),
				To:       time.Unix(to, 0),
				ClientID: site_id,
			},
			group_by,
		)

		group_events := map[string]omisocial.GroupEvents{}
		for _, event := range events {
			if val, ok := group_events[event.Period]; ok {
				val.Events = append(
					val.Events,
					omisocial.GroupEvent{
						Name:     event.Name,
						Visitors: event.Visitors,
						Views:    event.Views,
						Sessions: event.Sessions,
					},
				)
				group_events[val.Period] = val
			} else {
				group_event := omisocial.GroupEvents{
					Period: event.Period,
					Events: []omisocial.GroupEvent{
						{
							Name:     event.Name,
							Visitors: event.Visitors,
							Views:    event.Views,
							Sessions: event.Sessions,
						},
					},
				}
				group_events[event.Period] = group_event
			}
		}

		keys := []string{}
		for key := range group_events {
			keys = append(keys, key)

		}
		sort.Strings(keys)

		data := []omisocial.OverTime{}
		for _, key := range keys {
			visitor := 0
			for _, item := range visitors {
				if item.Period == group_events[key].Period {
					visitor = item.Visitors
				}
			}
			data = append(
				data,
				omisocial.OverTime{
					Period:   group_events[key].Period,
					Visitors: visitor,
					Events:   group_events[key].Events,
				},
			)
		}

		jData, _ := json.Marshal(&omisocial.ResponseOverTime{
			Message: "",
			Error:   false,
			Data:    data,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))

	http.Handle("/report/utm-sources", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		analyzer := omisocial.NewAnalyzer(store)

		from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
		site_id, _ := strconv.ParseInt(r.URL.Query().Get("site_id"), 10, 64)
		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		page_size, _ := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)

		if from == 0 || to == 0 || site_id == 0 || from > to {
			jData, _ := json.Marshal(&omisocial.Response{
				Message: "Invalid input data",
				Error:   true,
				Data:    nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		if page < 1 {
			page = 1
		}

		if page_size <= 0 {
			page_size = 25
		}

		count, _ := analyzer.UTMSourceCount(&omisocial.Filter{
			From:     time.Unix(from, 0),
			To:       time.Unix(to, 0),
			ClientID: site_id,
		})

		if count == 0 {
			jData, _ := json.Marshal(&omisocial.ResponseUTMSources{
				Message:    "No data",
				Error:      false,
				TotalPages: 0,
				Count:      count,
				Page:       int(page),
				PageSize:   int(page_size),
				Data:       nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		var total_pages = int64(math.Ceil(float64(count) / float64(page_size)))

		if page > total_pages {
			page = total_pages
		}

		offset := (int(page) - 1) * int(page_size)

		sources, _ := analyzer.UTMSource(&omisocial.Filter{
			From:     time.Unix(from, 0),
			To:       time.Unix(to, 0),
			ClientID: site_id,
			Limit:    int(page_size),
			Offset:   int(offset),
		})

		jData, _ := json.Marshal(&omisocial.ResponseUTMSources{
			Message:    "",
			Error:      false,
			TotalPages: 0,
			Count:      count,
			Page:       int(page),
			PageSize:   int(page_size),
			Data:       sources,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))

	http.Handle("/report/otm-sources", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		analyzer := omisocial.NewAnalyzer(store)

		from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
		site_id, _ := strconv.ParseInt(r.URL.Query().Get("site_id"), 10, 64)
		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		page_size, _ := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 64)

		if from == 0 || to == 0 || site_id == 0 || from > to {
			jData, _ := json.Marshal(&omisocial.Response{
				Message: "Invalid input data",
				Error:   true,
				Data:    nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		if page < 1 {
			page = 1
		}

		if page_size <= 0 {
			page_size = 25
		}

		count, _ := analyzer.OTMSourceCount(&omisocial.Filter{
			From:     time.Unix(from, 0),
			To:       time.Unix(to, 0),
			ClientID: site_id,
		})

		if count == 0 {
			jData, _ := json.Marshal(&omisocial.ResponseOTMSources{
				Message:    "No data",
				Error:      false,
				TotalPages: 0,
				Count:      count,
				Page:       int(page),
				PageSize:   int(page_size),
				Data:       nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		var total_pages = int64(math.Ceil(float64(count) / float64(page_size)))

		if page > total_pages {
			page = total_pages
		}

		offset := (int(page) - 1) * int(page_size)

		sources, _ := analyzer.OTMSource(&omisocial.Filter{
			From:     time.Unix(from, 0),
			To:       time.Unix(to, 0),
			ClientID: site_id,
			Limit:    int(page_size),
			Offset:   int(offset),
		})

		jData, _ := json.Marshal(&omisocial.ResponseOTMSources{
			Message:    "",
			Error:      false,
			TotalPages: 0,
			Count:      count,
			Page:       int(page),
			PageSize:   int(page_size),
			Data:       sources,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}))

	// And finally, start the server.
	// We don't flush hits on shutdown but you should add that in a real application by calling Tracker.Flush().
	log.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)

}
