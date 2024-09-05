package rose

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/go-chi/chi/v5"
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

func TestParser_ParseHTTP(t *testing.T) {
	type fields struct {
		ctx           *ctxx.Context
		textMsgPrefix string
	}
	type args struct {
		r      *http.Request
		target any
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       Rose
		wantTarget any
		wantErr    bool
	}{
		{
			name: "valid - request body",
			fields: fields{
				ctx:           ctxx.Background(),
				textMsgPrefix: "",
			},
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name": "test"}`))
					return r
				}(),
				target: &struct {
					Name string `rose:"name"`
				}{},
			},
			want: Rose{
				Valid: true,
			},
			wantTarget: &struct {
				Name string `rose:"name"`
			}{
				Name: "test",
			},
			wantErr: false,
		},
		{
			name: "valid - request query",
			fields: fields{
				ctx:           ctxx.Background(),
				textMsgPrefix: "",
			},
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("GET", "/?name=test", nil)
					return r
				}(),
				target: &struct {
					Name string `rose:"name"`
				}{},
			},
			want: Rose{
				Valid: true,
			},
			wantTarget: &struct {
				Name string `rose:"name"`
			}{
				Name: "test",
			},
			wantErr: false,
		},
		{
			name: "valid - url params",
			fields: fields{
				ctx:           ctxx.Background(),
				textMsgPrefix: "",
			},
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("GET", "/test?name=test", nil)
					rctx := chi.NewRouteContext()
					rctx.URLParams.Keys = []string{"name"}
					rctx.URLParams.Values = []string{"test"}
					ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
					r = r.WithContext(ctx)
					return r
				}(),
				target: &struct {
					Name string `rose:"name"`
				}{},
			},
			want: Rose{
				Valid: true,
			},
			wantTarget: &struct {
				Name string `rose:"name"`
			}{
				Name: "test",
			},
			wantErr: false,
		},
		{
			name: "invalid - request body",
			fields: fields{
				ctx:           ctxx.Background(),
				textMsgPrefix: "",
			},
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"name": "abc"}`))
					return r
				}(),
				target: &struct {
					Name int `rose:"name"`
				}{},
			},
			want: Rose{
				Valid: false,
			},
			wantTarget: &struct {
				Name int `rose:"name"`
			}{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				ctx:           tt.fields.ctx,
				textMsgPrefix: tt.fields.textMsgPrefix,
			}
			got, err := p.ParseHTTP(tt.args.r, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.ParseHTTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.ParseHTTP() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(tt.wantTarget, tt.args.target) {
				t.Errorf("Parser.ParseHTTP().target = %#v, want %#v", tt.args.target, tt.wantTarget)
			}
		})
	}
}

func TestJSONMarshal(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				v: struct {
					Name string `rose:"name"`
				}{
					Name: "test",
				},
			},
			want:    []byte(`{"name":"test"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONMarshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONMarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONMarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
