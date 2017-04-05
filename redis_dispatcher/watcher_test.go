package redis_dispatcher

import "testing"

func TestWatcher_Watch(t *testing.T) {
	type fields struct {
		QueueName   string
		Redis       *RedisClient
		HandlerFunc func(job map[string]string) (string, error)
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Watcher{
				QueueName:   tt.fields.QueueName,
				Redis:       tt.fields.Redis,
				HandlerFunc: tt.fields.HandlerFunc,
			}
			w.Watch()
		})
	}

	// c := NewRedisClient()
	// h := func(job map[string]string) (string, error) {
	// 	return "abc,def", nil
	// }

	// w := Watcher{QueueName: "test", Redis: c, HandlerFunc: h}
	// w.Watch()
}

func Test_runHandler(t *testing.T) {
	type args struct {
		handler *JobHandler
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runHandler(tt.args.handler)
		})
	}
}
