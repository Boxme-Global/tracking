package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	omisocial "github.com/Boxme-Global/tracking/src"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

type Response struct {
	Message string                   `json:"message"`
	Error   bool                     `json:"error"`
	Data    []omisocial.VisitorStats `json:"data"`
}

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

		if eventName != "" && Contains(pageloadEvents, eventName) {
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

		log.Println(from, to, site_id, from > to)
		if from == 0 || to == 0 || site_id == 0 || from > to {
			jData, _ := json.Marshal(&Response{
				"Invalid input data",
				true,
				nil,
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		visitors, _ := analyzer.Visitors(&omisocial.Filter{
			From:     time.Unix(from, 0),
			To:       time.Unix(to, 0),
			ClientID: site_id,
		})

		jData, _ := json.Marshal(&Response{
			"",
			false,
			visitors,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
		return
	}))

	// And finally, start the server.
	// We don't flush hits on shutdown but you should add that in a real application by calling Tracker.Flush().
	log.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)

}
