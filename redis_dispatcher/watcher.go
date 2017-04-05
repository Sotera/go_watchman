package redis_dispatcher

import (
	"log"
	"net"
	"time"
)

type Watcher struct {
	QueueName   string
	Redis       *RedisClient
	HandlerFunc func(job map[string]string) (string, error)
}

func (w *Watcher) Watch() {
	log.Printf("watching %s...\n", w.QueueName)
	for {
		res, err := w.Redis.C.BRPop(10*time.Second, w.QueueName).Result()
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				log.Fatal(err)
			}
			// keep trying for all other errors
			continue
		}
		log.Println("job recvd", res)

		handler := &JobHandler{
			key:         res[1],
			redis:       w.Redis,
			handlerFunc: w.HandlerFunc,
		}
		go runHandler(handler)
	}
}

func runHandler(handler *JobHandler) {
	err := handler.handle()
	if err != nil {
		log.Println(err)
	}
}
