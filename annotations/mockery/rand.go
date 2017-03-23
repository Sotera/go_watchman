package mockery

import (
	"math/rand"
	"time"
)

// takeOne randomly takes from args or use default set.
func takeOne(vals []string) string {
	rand.Seed(time.Now().Unix())
	if len(vals) == 0 {
		vals = []string{
			"riot",
			"protest",
			"concert",
			"natural disaster",
		}
	}
	n := rand.Int() % len(vals)
	return vals[n]
}

// generateStr returns n-length list of random fixed-length lowercase, alphanum strings.
func generateStr(n, length int) []string {
	// don't include hard-to-read alphanums
	letters := []rune("abcdefghijkmnpqrstuvwxyz023456789")
	out := []string{}

	for i := 0; i < n; i++ {
		s := make([]rune, length)
		for i := range s {
			s[i] = letters[rand.Intn(len(letters))]
		}
		out = append(out, string(s))
	}

	return out
}
