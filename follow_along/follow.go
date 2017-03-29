package main

import (
	"errors"
	"fmt"
	"net/http"

	"flag"

	"strings"

	"golang.org/x/net/html"
)

func main() {
	follower := flag.String("follower", "", "screen_name is following")
	followee := flag.String("followee", "", "screen_name being followed")

	flag.Parse()

	if *followee == "" || *follower == "" {
		panic(errors.New("missing required args"))
	}

	found, err := isFollowing(*follower, *followee, fmt.Sprintf("https://mobile.twitter.com/%s/following", *follower))

	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(found)
}

func isFollowing(follower, followee, url string) (bool, error) {
	// use mobile site b/c desktop site loads followers info via javascript.
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer res.Body.Close()

	found, nextPageURL := findScreenName(res, followee)
	if found {
		return found, nil
	} else if nextPageURL != "" {
		url := fmt.Sprintf("https://mobile.twitter.com%s", nextPageURL)
		return isFollowing(follower, followee, url)
	} else {
		return false, nil
	}
}

func findScreenName(res *http.Response, screenName string) (bool, string) {
	z := html.NewTokenizer(res.Body)
	var found, morePages bool
	var nextPageURL string

	_ = morePages

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return found, nextPageURL
		case tt == html.TextToken:
			t := z.Token()
			// fmt.Println(t)

			if strings.ToLower(t.Data) == "show more people" {
				fmt.Println("more pages")
				morePages = true
			}
		case tt == html.StartTagToken:
			t := z.Token()

			// look for <a href="/POTUS44?p=s" data-scribe-action="profile_click">
			if t.Data == "a" {
				for _, a := range t.Attr {
					// fmt.Println(a.Key)
					if a.Key == "href" {
						//TODO: use ToLowerSpecial() ?
						if strings.Contains(strings.ToLower(a.Val), strings.ToLower(screenName)) {
							found = true
							break
						}
						if strings.Contains(strings.ToLower(a.Val), "/following?cursor") {
							nextPageURL = a.Val
						}
					}
				}
			}
		}
	}
}
