package omisocial

import (
	"os"
	"testing"
	"time"
)

var dbClient *Client

func TestMain(m *testing.M) {
	if err := Migrate("clickhouse://10.148.15.231:9000?x-multi-statement=true"); err != nil {
		panic(err)
	}

	c, err := NewClient("tcp://10.148.15.231:9000", nil)

	if err != nil {
		panic(err)
	}

	dbClient = c
	defer func() {
		if err := dbClient.DB.Close(); err != nil {
			panic(err)
		}
	}()
	os.Exit(m.Run())
}

func cleanupDB() {
	dbClient.MustExec(`ALTER TABLE "page_view" DELETE WHERE 1=1`)
	dbClient.MustExec(`ALTER TABLE "session" DELETE WHERE 1=1`)
	dbClient.MustExec(`ALTER TABLE "event" DELETE WHERE 1=1`)
	dbClient.MustExec(`ALTER TABLE "user_agent" DELETE WHERE 1=1`)
	time.Sleep(time.Millisecond * 20)
}
