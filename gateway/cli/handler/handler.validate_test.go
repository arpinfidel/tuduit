package handler

import (
	"reflect"
	"testing"

	"github.com/arpinfidel/tuduit/app"
)

func TestHandler_validateRaw(t *testing.T) {
	type args struct {
		ctx     *app.Context
		runArgs []string
		v       any
	}
	tests := []struct {
		name    string
		h       *Handler
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "valid simple",
			h:    &Handler{},
			args: args{
				ctx:     &app.Context{},
				runArgs: []string{"test"},
				v: &struct {
					Title string `tuduit:"title,required"`
				}{},
			},
			want: &struct {
				Title string `tuduit:"title,required"`
			}{
				Title: "test",
			},
			wantErr: false,
		},
		{
			name: "invalid simple",
			h:    &Handler{},
			args: args{
				ctx:     &app.Context{},
				runArgs: []string{},
				v: &struct {
					Title string `tuduit:"title,required"`
				}{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid variadic",
			h:    &Handler{},
			args: args{
				ctx:     &app.Context{},
				runArgs: []string{"test", "test2"},
				v: &struct {
					Title []string `tuduit:"title,required"`
				}{},
			},
			want: &struct {
				Title []string `tuduit:"title,required"`
			}{
				Title: []string{"test", "test2"},
			},
			wantErr: false,
		},
		{
			name: "invalid variadic",
			h:    &Handler{},
			args: args{
				ctx:     &app.Context{},
				runArgs: []string{},
				v: &struct {
					Title []string `tuduit:"title,required"`
				}{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid complex",
			h:    &Handler{},
			args: args{
				ctx:     &app.Context{},
				runArgs: []string{"test", "1", "test2", "test3"},
				v: &struct {
					Title  string   `tuduit:"title,required"`
					Number int      `tuduit:"number,required"`
					Tags   []string `tuduit:"tags"`
				}{},
			},
			want: &struct {
				Title  string   `tuduit:"title,required"`
				Number int      `tuduit:"number,required"`
				Tags   []string `tuduit:"tags"`
			}{
				Title:  "test",
				Number: 1,
				Tags:   []string{"test2", "test3"},
			},
			wantErr: false,
		},
		{
			name: "invalid complex",
			h:    &Handler{},
			args: args{
				ctx:     &app.Context{},
				runArgs: []string{"test"},
				v: &struct {
					Title  string   `tuduit:"title,required"`
					Number int      `tuduit:"number,required"`
					Tags   []string `tuduit:"tags"`
				}{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.h.validateArgsRaw(tt.args.ctx, tt.args.runArgs, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Handler.validateRaw() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.validateRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}
