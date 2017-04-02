package follow_along

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	rd "github.com/Sotera/go_watchman/redis_dispatcher"
)

func TestHTTPFetcher_Fetch(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		f       *HTTPFetcher
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &HTTPFetcher{}
			got, err := f.Fetch(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPFetcher.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTTPFetcher.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewScraper(t *testing.T) {
	type args struct {
		follower string
	}
	tests := []struct {
		name string
		args args
		want *Scraper
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewScraper(tt.args.follower); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewScraper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScraper_Followees(t *testing.T) {
	type fields struct {
		F            Fetcher
		follower     string
		url          string
		followees    Set
		currPage     int
		maxFollowees int
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scraper{
				F:            tt.fields.F,
				follower:     tt.fields.follower,
				url:          tt.fields.url,
				followees:    tt.fields.followees,
				currPage:     tt.fields.currPage,
				maxFollowees: tt.fields.maxFollowees,
			}
			if got := s.Followees(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scraper.Followees() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScraper_URL(t *testing.T) {
	type fields struct {
		F            Fetcher
		follower     string
		url          string
		followees    Set
		currPage     int
		maxFollowees int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scraper{
				F:            tt.fields.F,
				follower:     tt.fields.follower,
				url:          tt.fields.url,
				followees:    tt.fields.followees,
				currPage:     tt.fields.currPage,
				maxFollowees: tt.fields.maxFollowees,
			}
			if got := s.URL(); got != tt.want {
				t.Errorf("Scraper.URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScraper_SetURL(t *testing.T) {
	type fields struct {
		F            Fetcher
		follower     string
		url          string
		followees    Set
		currPage     int
		maxFollowees int
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scraper{
				F:            tt.fields.F,
				follower:     tt.fields.follower,
				url:          tt.fields.url,
				followees:    tt.fields.followees,
				currPage:     tt.fields.currPage,
				maxFollowees: tt.fields.maxFollowees,
			}
			s.SetURL(tt.args.path)
		})
	}
}

func TestScraper_SetMaxFollowees(t *testing.T) {
	type fields struct {
		F            Fetcher
		follower     string
		url          string
		followees    Set
		currPage     int
		maxFollowees int
	}
	type args struct {
		limit int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scraper{
				F:            tt.fields.F,
				follower:     tt.fields.follower,
				url:          tt.fields.url,
				followees:    tt.fields.followees,
				currPage:     tt.fields.currPage,
				maxFollowees: tt.fields.maxFollowees,
			}
			s.SetMaxFollowees(tt.args.limit)
		})
	}
}

func TestScraper_IsFollowing(t *testing.T) {
	type fields struct {
		F            Fetcher
		follower     string
		url          string
		followees    Set
		currPage     int
		maxFollowees int
	}
	type args struct {
		followee string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scraper{
				F:            tt.fields.F,
				follower:     tt.fields.follower,
				url:          tt.fields.url,
				followees:    tt.fields.followees,
				currPage:     tt.fields.currPage,
				maxFollowees: tt.fields.maxFollowees,
			}
			got, err := s.IsFollowing(tt.args.followee)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scraper.IsFollowing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scraper.IsFollowing() = %v, want %v", got, tt.want)
			}
		})
	}

	redis := rd.NewRedisClient()
	var err error

	redis.C.FlushAll()

	job := map[string]interface{}{"state": "new", "id": "hillaryclinton"}

	_, err = redis.C.HMSet("1", job).Result()
	if err != nil {
		panic(err)
	}

	_, err = redis.C.LPush("queue", "1").Result()
	if err != nil {
		panic(err)
	}

	res, err := redis.C.BRPop(5*time.Second, "queue").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	job1, err := redis.C.HGetAll(res[1]).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(job1)

	s := NewScraper(job1["id"])
	s.SetMaxFollowees(10)
	_, err = s.IsFollowing("")
	fmt.Println(err, s.Followees())

	j := map[string]interface{}{}
	for k, v := range job1 {
		j[k] = v
	}
	j["data"] = strings.Join(s.Followees(), ",")
	j["state"] = "processed"
	fmt.Println(j)
	fmt.Println(len(s.Followees()))
	_, err = redis.C.HMSet("1", j).Result()
	if err != nil {
		panic(err)
	}
}

func TestScraper_findFollowee(t *testing.T) {
	type fields struct {
		F            Fetcher
		follower     string
		url          string
		followees    Set
		currPage     int
		maxFollowees int
	}
	type args struct {
		markup   io.Reader
		followee string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scraper{
				F:            tt.fields.F,
				follower:     tt.fields.follower,
				url:          tt.fields.url,
				followees:    tt.fields.followees,
				currPage:     tt.fields.currPage,
				maxFollowees: tt.fields.maxFollowees,
			}
			got, got1 := s.findFollowee(tt.args.markup, tt.args.followee)
			if got != tt.want {
				t.Errorf("Scraper.findFollowee() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Scraper.findFollowee() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_strmatch(t *testing.T) {
	type args struct {
		s      string
		substr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strmatch(tt.args.s, tt.args.substr); got != tt.want {
				t.Errorf("strmatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
