package redis_dispatcher

import (
	"fmt"
)

type JobHandler struct {
	key         string
	redis       *RedisClient
	job         map[string]string
	handlerFunc func(job map[string]string) (string, error)
}

func (jh *JobHandler) handle() error {
	var err error
	jh.job, err = jh.redis.C.HGetAll(jh.key).Result()
	if err != nil {
		return err
	}
	fmt.Println("job", jh.job)

	data, err := jh.handlerFunc(jh.job)
	if err != nil {
		_, err = jh.update("", "error", err)
		return err
	}

	_, err = jh.update(data, "processed", nil)
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

	fmt.Println("updated job", job)
	return jh.redis.C.HMSet(jh.key, job).Result()
}
