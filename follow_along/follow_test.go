package main

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func Test_findScreenName(t *testing.T) {
	type args struct {
		res        *http.Response
		screenName string
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
			if got := findScreenName(tt.args.res, tt.args.screenName); got != tt.want {
				t.Errorf("findScreenName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isFollowing(t *testing.T) {
	// type args struct {
	// 	follower string
	// 	followee string
	// }
	// tests := []struct {
	// 	name string
	// 	args args
	// 	want bool
	// }{
	// // TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		if got := isFollowing(tt.args.follower, tt.args.followee); got != tt.want {
	// 			t.Errorf("isFollowing() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }

	for i := 0; i < 1000; i++ {
		strconv.Itoa(i)
		found, err := isFollowing("POTUS44", strconv.Itoa(i+1*1000))
		fmt.Println(found, err, i)
		found, err = isFollowing("lukewendling", "POTUS44")
		fmt.Println(found, err, i)
	}
}
