package main

import (
	"fmt"
	"os"
	"strings"

	"strconv"

	f "github.com/Sotera/go_watchman/follow_along"
	rd "github.com/Sotera/go_watchman/redis_dispatcher"
)

func main() {
	q := os.Getenv("QUEUE_NAME")
	if q == "" {
		q = "genie:followfinder"
	}

	handler := func(job map[string]string) (string, error) {
		sn := f.NewScraper(job["id"])
		max, err := strconv.Atoi(job["max"])
		if err == nil {
			sn.SetMaxFollowees(max)
		}
		_, err = sn.IsFollowing("")
		if err != nil {
			return "", err
		}
		fmt.Println("found", sn.Followees())

		return strings.Join(sn.Followees(), ","), nil
	}

	redis := rd.NewRedisClient()
	w := rd.Watcher{
		QueueName:   q,
		Redis:       redis,
		HandlerFunc: handler,
	}
	w.Watch()
}
