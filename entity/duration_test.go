package entity

import (
	"fmt"
	"testing"
	"time"
)

func TestDuration_Format(t *testing.T) {
	dd := 365*24*time.Hour + 7*24*time.Hour + 24*time.Hour + 1*time.Hour + 1*time.Minute + 1*time.Second
	type args struct {
		format string
		d      Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "%s",
			args: args{
				format: "%s",
				d:      Duration(dd),
			},
			want: "1y 1w 1d 1h 1m 1s",
		},
		{
			name: "%2s",
			args: args{
				format: "%2s",
				d:      Duration(dd),
			},
			want: "1y 1w",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf(tt.args.format, tt.args.d); got != tt.want {
				t.Errorf("Duration.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
