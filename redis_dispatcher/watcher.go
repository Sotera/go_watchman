package redis_dispatcher

import (
	"fmt"
	"time"
)

type Watcher struct {
	QueueName   string
	Redis       *RedisClient
	HandlerFunc func(job map[string]string) (string, error)
}

func (w *Watcher) Watch() {
	for {
		res, _ := w.Redis.C.BRPop(10*time.Second, w.QueueName).Result()
		if res == nil {
			continue
		}
		fmt.Println("job recvd", res)

		handler := &JobHandler{
			key:         res[1],
			redis:       w.Redis,
			handlerFunc: w.HandlerFunc,
			finalState:  "processed",
		}
		go runHandler(handler)
	}
}

func runHandler(handler *JobHandler) {
	err := handler.handle()
	if err != nil {
		fmt.Println(err)
	}
}
