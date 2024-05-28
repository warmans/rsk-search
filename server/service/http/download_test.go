package http

import "testing"

func TestFormatGifText(t *testing.T) {
	type args struct {
		maxLineLength int
		lines         []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no words given",
			args: args{
				maxLineLength: 5,
				lines:         []string{},
			},
			want: "",
		},
		{
			name: "single line fits within max length",
			args: args{
				maxLineLength: 10,
				lines:         []string{"foo"},
			},
			want: "foo",
		},
		{
			name: "many lines where each fits within max length",
			args: args{
				maxLineLength: 10,
				lines:         []string{"foo", "bar", "baz"},
			},
			want: "foo\nbar\nbaz",
		},
		{
			name: "line does not fit within max length",
			args: args{
				maxLineLength: 7,
				lines:         []string{"foo", "bar bar bar", "baz baz"},
			},
			want: "foo\nbar bar\nbar\nbaz baz",
		},
		{
			name: "single word exceeds line length",
			args: args{
				maxLineLength: 10,
				lines:         []string{"foo baaaaaaaaaar", "baz baz"},
			},
			want: "foo\nbaaaaaaaaaar\nbaz baz",
		},
		{
			name: "single word exceeds line length",
			args: args{
				maxLineLength: 10,
				lines:         []string{"foo baaaaaaaaaar", "baz baz"},
			},
			want: "foo\nbaaaaaaaaaar\nbaz baz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatGifText(tt.args.maxLineLength, tt.args.lines); got != tt.want {
				t.Errorf("FormatGifText() = %v, want %v", got, tt.want)
			}
		})
	}
}
