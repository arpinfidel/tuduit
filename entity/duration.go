package entity

import (
	"fmt"
	"strings"
	"time"
)

var _ fmt.Stringer = Duration(0)

type Duration time.Duration

func (d Duration) String() string {
	type unit struct {
		name string
		val  time.Duration
	}
	units := []unit{
		{"y", time.Hour * 24 * 365},
		// {"mo", time.Hour * 24 * 30},
		{"w", time.Hour * 24 * 7},
		{"d", time.Hour * 24},
		{"h", time.Hour},
		{"m", time.Minute},
		{"s", time.Second},
	}

	var parts []string
	for _, u := range units {
		dd := time.Duration(d)
		if dd > u.val {
			parts = append(parts, fmt.Sprintf("%d%s", dd/u.val, u.name))
			dd = dd % u.val
			d = Duration(dd)
		}
	}

	return strings.Join(parts, "")
}
