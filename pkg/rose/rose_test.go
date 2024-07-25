package rose

import (
	"reflect"
	"testing"
)

func Test_parseArgs(t *testing.T) {
	type example struct {
		Name string `rose:"name"`
		Age  int    `rose:"age"`
		Data []int  `rose:"data"`
	}

	type args struct {
		args  []string
		flags map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    example
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				args: []string{"test", "10"},
				flags: map[string]string{
					"data": "[1,2,3]",
				},
			},
			want: example{
				Name: "test",
				Age:  10,
				Data: []int{1, 2, 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs[example](tt.args.args, tt.args.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseJSON(t *testing.T) {
	type example struct {
		Name string `rose:"name"`
		Age  int    `rose:"age"`
		Data []int  `rose:"data"`
	}

	type args struct {
		jsonBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    example
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				jsonBytes: []byte(`{"name":"test","age":10,"data":[1,2,3]}`),
			},
			want: example{
				Name: "test",
				Age:  10,
				Data: []int{1, 2, 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseJSON[example](tt.args.jsonBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
