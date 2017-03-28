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

	found, err := isFollowing(*follower, *followee)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(found)
}

func isFollowing(follower, followee string) (bool, error) {
	// use mobile site b/c desktop site loads followers info via javascript.
	res, err := http.Get(fmt.Sprintf("https://mobile.twitter.com/%s/following", follower))
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer res.Body.Close()

	return findScreenName(res, followee), nil
}

func findScreenName(res *http.Response, screenName string) bool {
	z := html.NewTokenizer(res.Body)
	found := false
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return found
		case tt == html.StartTagToken:
			t := z.Token()

			// look for <a href="/POTUS44?p=s" data-scribe-action="profile_click">
			if t.Data == "a" {
				for _, a := range t.Attr {
					// fmt.Println(a.Key)
					if a.Key == "href" {
						// fmt.Println("Found user:", a.Val)
						//TODO: use ToLowerSpecial() ?
						if strings.Contains(strings.ToLower(a.Val), strings.ToLower(screenName)) {
							found = true
							break
						}
					}
				}
			}
		}
	}
}
