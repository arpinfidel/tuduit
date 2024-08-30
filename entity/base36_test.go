package entity

import (
	"testing"
)

func TestBase36_UnmarshalYAML(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				s: "",
			},
			wantErr: true,
		},
		{
			name: "zero",
			args: args{
				s: "0",
			},
			want: 0,
		},
		{
			name: "one",
			args: args{
				s: "1",
			},
			want: 1,
		},
		{
			name: "361",
			args: args{
				s: "a1",
			},
			want: 361,
		},
		{
			name: "uppercase",
			args: args{
				s: "A1",
			},
			want: 361,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base36[uint64]{}
			if err := b.UnmarshalYAML(func(a any) error {
				*a.(*string) = tt.args.s
				return nil
			}); (err != nil) != tt.wantErr {
				t.Errorf("Base36.UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if b.V != tt.want {
				t.Errorf("Base36.UnmarshalYAML() = %v, want %v", b.V, tt.want)
			}
		})
	}
}

func TestBase36_String(t *testing.T) {
	type args struct {
		v uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				v: 0,
			},
			want: "0",
		},
		{
			name: "1",
			args: args{
				v: 1,
			},
			want: "1",
		},
		{
			name: "10",
			args: args{
				v: 10,
			},
			want: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Base36[uint]{V: tt.args.v}
			if got := b.String(); got != tt.want {
				t.Errorf("Base36.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
