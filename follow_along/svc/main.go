package main

import (
	"log"
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

	cacheSize := os.Getenv("CACHE_SIZE")
	if cacheSize == "" {
		cacheSize = "100"
	}

	// maintain an internal cache to help clients.
	cs, _ := strconv.Atoi(cacheSize)
	c := cache{maxSize: cs}

	handler := func(job map[string]string) (string, error) {
		if ci, hit := c.get(job["id"]); hit {
			log.Println("cache hit:", job["id"])
			return ci.value.(string), nil
		}
		scr := f.NewScraper(job["id"])
		max, err := strconv.Atoi(job["max"])
		if err == nil {
			scr.SetMaxFollowees(max)
		}
		_, err = scr.IsFollowing("")
		if err != nil {
			return "", err
		}
		log.Println("found", scr.Followees())

		followees := strings.Join(scr.Followees(), ",")
		c.add(cacheItem{key: job["id"], value: followees})
		return followees, nil
	}

	redis := rd.NewRedisClient()
	w := rd.Watcher{
		QueueName:   q,
		Redis:       redis,
		HandlerFunc: handler,
	}
	w.Watch()
}
