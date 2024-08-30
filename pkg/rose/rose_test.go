package rose

import (
	"reflect"
	"testing"

	"github.com/arpinfidel/tuduit/entity"
)

func Test_parseArgs(t *testing.T) {
	type example struct {
		Name string `rose:"name"`
		Age  int    `rose:"age"`
		Data []int  `rose:"data"`
	}

	type example2 struct {
		Help    bool    `rose:"help"`
		Example example `rose:"flatten="`
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
			got := example{}
			_, err := parseArgs(tt.args.args, tt.args.flags, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArgs() = %v, want %v", got, tt.want)
			}
		})
	}

	tests2 := []struct {
		name    string
		args    args
		want    example2
		wantErr bool
	}{
		{
			name: "success - flatten",
			args: args{
				flags: map[string]string{
					"data": "[1,2,3]",
					"help": "true",
				},
			},
			want: example2{
				Help: true,
				Example: example{
					Data: []int{1, 2, 3},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			got := example2{}
			_, err := parseArgs(tt.args.args, tt.args.flags, &got)
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
			got := example{}
			_, err := parseJSON(tt.args.jsonBytes, &got)
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

func Test_help(t *testing.T) {
	type example struct {
		Name string `rose:"name,required="`
		Age  int    `rose:"age,default=10"`
		Data []int  `rose:"data"`
	}

	type example2 struct {
		Help    bool    `rose:"help"`
		Example example `rose:"flatten="`
	}

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "success",
			want: "help: bool optional\nname: string required\nage: int optional, default=10\ndata: slice optional\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := help(example2{})
			if (err != nil) != tt.wantErr {
				t.Errorf("help() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("help() = %#v, want %v", got, tt.want)
			}
		})
	}
}

func Test_castType(t *testing.T) {
	type args struct {
		v string
		t reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		wantVal any
		wantErr bool
	}{
		{
			name: "base 36",
			args: args{
				v: "1",
				t: reflect.TypeOf(entity.Base36[uint64]{V: 0}),
			},
			wantVal: &entity.Base36[uint64]{V: 1},
			wantErr: false,
		},
		{
			name: "base 36 slice",
			args: args{
				v: "1,2,3",
				t: reflect.TypeOf([]entity.Base36[uint64]{}),
			},
			wantVal: []entity.Base36[uint64]{{V: 1}, {V: 2}, {V: 3}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := castType(tt.args.v, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("castType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("castType() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}
