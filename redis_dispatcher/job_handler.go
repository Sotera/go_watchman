package redis_dispatcher

import (
	"fmt"
	"log"
)

type JobHandler struct {
	key          string
	redis        *RedisClient
	job          map[string]string
	handlerFunc  func(job map[string]string) (string, error)
	initialState string
	finalState   string
}

func (jh *JobHandler) handle() error {
	var err error
	if jh.initialState == "" {
		jh.initialState = "new"
	}
	if jh.finalState == "" {
		jh.finalState = "processed"
	}

	// TODO: check for valid initial state, like "new"?

	jh.job, err = jh.redis.C.HGetAll(jh.key).Result()
	if err != nil {
		return err
	}
	log.Println("job", jh.job)

	data, err := jh.handlerFunc(jh.job)
	if err != nil {
		_, err = jh.update("", "error", err)
		return err
	}

	_, err = jh.update(data, jh.finalState, nil)
	if err != nil {
		return err
	}

	return nil
}

func (jh *JobHandler) update(data string, state string, err error) (string, error) {
	// conversion for hmset argument
	job := map[string]interface{}{}
	for k, v := range jh.job {
		job[k] = v
	}
	job["data"] = data
	job["state"] = state
	if err != nil {
		job["error"] = fmt.Sprintf("%v", err)
	}

	log.Println("updated job", job)
	return jh.redis.C.HMSet(jh.key, job).Result()
}
