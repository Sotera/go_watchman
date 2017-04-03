package follow_along

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"strings"

	"github.com/Sotera/go_watchman/util/collections"
	"golang.org/x/net/html"
)

// system-defined max to prevent runaway paging
const MAX_PAGES = 150

type Fetcher interface {
	Fetch(url string) (io.ReadCloser, error)
}

type Scraper struct {
	F            Fetcher
	follower     string
	url          string
	followees    collections.Set
	currPage     int
	maxFollowees int // 0 indicates no limit
}

type HTTPFetcher struct {
}

func (f *HTTPFetcher) Fetch(url string) (io.ReadCloser, error) {
	// Default client has no timeout
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func NewScraper(follower string) *Scraper {
	return &Scraper{
		F:        &HTTPFetcher{},
		follower: follower,
	}
}

func (s *Scraper) Followees() []string {
	f := make([]string, len(s.followees.Items()))
	i := 0
	for _, v := range s.followees.Items() {
		f[i] = v.(string)
		i++
	}
	return f
}

func (s *Scraper) URL() string {
	// use mobile site b/c desktop site loads followings info via javascript.
	if s.url == "" {
		// set initial following URL.
		s.url = fmt.Sprintf("https://mobile.twitter.com/%s/following", s.follower)
	}
	return s.url
}

func (s *Scraper) SetURL(path string) {
	// update url with cursor/paging path
	s.url = fmt.Sprintf("https://mobile.twitter.com%s", path)
}

func (s *Scraper) SetMaxFollowees(limit int) {
	if limit < 0 {
		limit = 0
	}
	s.maxFollowees = limit
}

func (s *Scraper) IsFollowing(followee string) (bool, error) {
	rc, err := s.F.Fetch(s.URL())
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// must close after tokenizing
	defer rc.Close()

	found, nextPagePath := s.findFollowee(rc, followee)
	s.currPage++
	if found {
		return found, nil
	} else if nextPagePath != "" {
		// not too many pages, like u might find with bots.
		if s.currPage >= MAX_PAGES {
			return found, nil
		}
		s.SetURL(nextPagePath)
		return s.IsFollowing(followee)
	} else {
		return false, nil
	}
}

func (s *Scraper) findFollowee(markup io.Reader, followee string) (bool, string) {
	z := html.NewTokenizer(markup)
	var found bool
	var nextPagePath string
	//screen name regex
	snRegex := regexp.MustCompile(`(?i)/(.+)\?p=s`)

tokens:
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return found, nextPagePath
		case tt == html.StartTagToken:
			t := z.Token()

			// look for <a href="/POTUS44?p=s" data-scribe-action="profile_click">
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						// is screen name link?
						m := snRegex.FindStringSubmatch(a.Val)
						if len(m) > 1 {
							// HACK: empty followee == collect all
							if followee == "" {
								s.followees.Add(m[1])
								numFollowees := len(s.followees.Items())
								if s.maxFollowees > 0 && numFollowees >= s.maxFollowees {
									break tokens
								}
							} else { // try to match followee
								if strmatch(followee, m[1]) {
									found = true
									break tokens
								}
							}

						}
						// look for next page anchor
						if strmatch(a.Val, "?cursor") {
							nextPagePath = a.Val
						}
					}
				}
			}
		}
	}

	return found, nextPagePath
}

func strmatch(s, substr string) bool {
	//TODO: use ToLowerSpecial() ?
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
