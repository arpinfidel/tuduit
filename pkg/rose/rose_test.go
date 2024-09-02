package rose

import (
	"reflect"
	"testing"
	"time"

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
			p := &Parser{}
			_, err := p.parseArgs(tt.args.args, tt.args.flags, &got)
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
			p := &Parser{}
			_, err := p.parseArgs(tt.args.args, tt.args.flags, &got)
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
			p := &Parser{}
			_, err := p.parseJSON(tt.args.jsonBytes, &got)
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
				t: reflect.TypeOf(&entity.Base36[uint64]{V: 0}),
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
			p := &Parser{}
			gotVal, err := p.castType(tt.args.v, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("castType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("castType() = %#v, want %#v", gotVal, tt.wantVal)
			}
		})
	}
}

func Test_ChangeTimezone(t *testing.T) {
	loc1, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	loc2, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	v1 := struct {
		T time.Time
	}{
		T: time.Date(2020, 1, 1, 0, 0, 0, 0, loc1),
	}

	tm := v1.T
	ChangeTimezone(&v1, loc2)
	if v1.T.Location() != loc2 {
		t.Errorf("ChangeTimezone() = %v, want %v", v1.T.Location(), loc2)
	}
	if !tm.Equal(v1.T) {
		t.Errorf("ChangeTimezone() = %v, want %v", tm, v1.T)
	}

	tm = v1.T
	ChangeTimezone(&v1, time.UTC)
	if v1.T.Location() != time.UTC {
		t.Errorf("ChangeTimezone() = %v, want %v", v1.T.Location(), time.UTC)
	}
	if !tm.Equal(v1.T) {
		t.Errorf("ChangeTimezone() = %v, want %v", tm, v1.T)
	}

	tm = time.Date(2020, 1, 1, 0, 0, 0, 0, loc1)
	v2 := struct {
		T *time.Time
	}{
		T: &tm,
	}

	tm = *v2.T
	ChangeTimezone(&v2, loc2)
	if v2.T.Location() != loc2 {
		t.Errorf("ChangeTimezone() = %v, want %v", v2.T.Location(), loc2)
	}
	if !tm.Equal(*v2.T) {
		t.Errorf("ChangeTimezone() = %v, want %v", tm, *v2.T)
	}

	v3 := []entity.Task{
		{
			StartDate: &tm,
		},
		{
			StartDate: &tm,
		},
	}

	tm = *v3[0].StartDate
	ChangeTimezone(&v3, loc2)
	if v3[0].StartDate.Location() != loc2 {
		t.Errorf("ChangeTimezone() = %v, want %v", v3[0].StartDate.Location(), loc2)
	}
	if !tm.Equal(*v3[0].StartDate) {
		t.Errorf("ChangeTimezone() = %v, want %v", tm, *v3[0].StartDate)
	}
}
