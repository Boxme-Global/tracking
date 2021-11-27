
package main
import (
	"github.com/Boxme-Global/tracking/src"
	"net/http"
	_ "github.com/lib/pq"
	"log"
)
func main() {
	// Set the key for SipHash. This should be called on startup (before generating the first fingerprint) and is NOT concurrency save.
	omisocial.SetFingerprintKeys(42, 123)

	// Migrate the database.
	omisocial.Migrate("clickhouse://10.148.0.23:9000?x-multi-statement=true")

	// Create a new ClickHouse client to save hits.
	store, _ := omisocial.NewClient("tcp://10.148.0.23:9000", nil)

	// Set up a default tracker with a salt.
	// This will buffer and store hits and generate sessions by default.
	tracker := omisocial.NewTracker(store, "salt", nil)

	// Create a handler to serve traffic.
	// We prevent tracking resources by checking the path. So a file on /my-file.txt won't create a new hit
	// but all page calls will be tracked.
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			go tracker.Hit(r, nil)
		}

		w.Write([]byte("<h1>Hello World!</h1>"))
	}))

	// And finally, start the server.
	// We don't flush hits on shutdown but you should add that in a real application by calling Tracker.Flush().
	log.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)

}