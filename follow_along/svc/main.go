package main

import (
	"fmt"
	"os"
	"strings"

	f "github.com/Sotera/go_watchman/follow_along"
	rd "github.com/Sotera/go_watchman/redis_dispatcher"
)

func main() {
	q := os.Getenv("QUEUE_NAME")
	if q == "" {
		q = "genie:followfinder"
	}

	handler := func(job map[string]string) (string, error) {
		s := f.NewScraper(job["id"])
		_, err := s.IsFollowing("")
		if err != nil {
			return "", err
		}
		fmt.Println("found", s.Followees())

		return strings.Join(s.Followees(), ","), nil
	}

	redis := rd.NewRedisClient()
	w := rd.Watcher{
		QueueName:   q,
		Redis:       redis,
		HandlerFunc: handler,
	}
	w.Watch()
}
