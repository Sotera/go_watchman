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
