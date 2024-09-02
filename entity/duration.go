package entity

import (
	"fmt"
	"strings"
	"time"
)

var _ fmt.Stringer = Duration(0)

type Duration time.Duration

func (d Duration) Format(s fmt.State, v rune) {
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

	switch v {
	default:
		type hideMethods Duration
		type Duration hideMethods
		fmtDirective := fmt.FormatString(s, v)
		fmt.Fprint(s, fmtDirective, Duration(d))

	case 's':
		var parts []string
		i := 0
		for _, u := range units {
			dd := time.Duration(d)
			if dd < u.val {
				continue
			}

			parts = append(parts, fmt.Sprintf("%d%s", dd/u.val, u.name))
			dd = dd % u.val
			d = Duration(dd)

			i++
			w, ok := s.Width()
			if ok && i+1 > w {
				break
			}
		}

		fmt.Fprint(s, strings.Join(parts, " "))
	}
}

func (d Duration) String() string {
	return fmt.Sprintf("%s", d)
}

func (d Duration) Elapsed() string {
	s := ""
	if d < 0 {
		s += fmt.Sprintf("in %2s", -d)
	} else if d > 0 {
		s += fmt.Sprintf("%2s ago", d)
	} else {
		s += "now"
	}

	return s
}
