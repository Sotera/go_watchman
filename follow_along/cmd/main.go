package main

import (
	"errors"
	"fmt"

	"flag"

	f "github.com/Sotera/go_watchman/follow_along"
)

func main() {
	follower := flag.String("follower", "", "screen_name is following")
	followee := flag.String("followee", "", "screen_name being followed")

	flag.Parse()

	if *followee == "" || *follower == "" {
		panic(errors.New("missing required args"))
	}

	scraper := f.NewScraper(*follower)

	found, err := scraper.IsFollowing(*followee)

	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(found)
}
