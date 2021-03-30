package tscript

import "testing"

func TestCorrectContent(t *testing.T) {
	type args struct {
		c string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "lowercase", args: args{c: "this is lowercase."}, want: "This is lowercase."},
		{name: "already capital case", args: args{c: "This is capitalcase."}, want: "This is capitalcase."},
		{name: "first char is a number", args: args{c: "5 is a number."}, want: "5 is a number."},
		{name: "first char is a special character", args: args{c: "£ is a symbol."}, want: "£ is a symbol."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CorrectContent(tt.args.c); got != tt.want {
				t.Errorf("CorrectContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
