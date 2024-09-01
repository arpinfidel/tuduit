package errs

import (
	"errors"
	"testing"
)

// func TestError_WithTrace(t *testing.T) {
// 	type fields struct {
// 		Base       error
// 		Attributes []error
// 		Trace      []string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   *Error
// 	}{
// 		{
// 			name: "with trace",
// 			fields: fields{
// 				Base:       nil,
// 				Attributes: nil,
// 				Trace:      nil,
// 			},
// 			want: &Error{
// 				Base:       nil,
// 				Attributes: nil,
// 				Trace:      nil,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e := &Error{
// 				Base:       tt.fields.Base,
// 				Attributes: tt.fields.Attributes,
// 				Trace:      tt.fields.Trace,
// 			}
// 			if got := e.WithTrace(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Error.WithTrace() =\n%#v, want\n%#v\n", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestDeferTrace(t *testing.T) {
	someErr := errors.New("some error")
	type args struct {
		err *error
		f   func() error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "success - nil",
			args: args{
				err: nil,
				f: func() (err error) {
					defer DeferTrace(&err)()
					return nil
				},
			},
		},
		{
			name: "success - is self",
			args: args{
				err: nil,
				f: func() (err error) {
					defer DeferTrace(&err)()
					return someErr
				},
			},
			want: someErr,
		},
		{
			name: "success - is base",
			args: args{
				err: nil,
				f: func() (err error) {
					defer DeferTrace(&err)()
					return someErr
				},
			},
			want: ErrTypeBase,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.f()
			if !errors.Is(got, tt.want) {
				t.Errorf("DeferTrace() = %#v, want %#v", got, tt.want)
			}

			// fmt.Printf(" >> debug >> GetTrace(got): %#v\n", GetTrace(got))
		})
	}
}
