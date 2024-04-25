package data

import "testing"

func Test_withinNPcnt(t *testing.T) {
	type args struct {
		x     int
		y     int
		total int
		pcnt  float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "within 10 percent over",
			args: args{
				x:     90,
				y:     99,
				total: 100,
				pcnt:  0.1,
			},
			want: true,
		},
		{
			name: "exactly 10 percent over",
			args: args{
				x:     90,
				y:     100,
				total: 100,
				pcnt:  0.1,
			},
			want: true,
		},
		{
			name: "more than 10 percent over",
			args: args{
				x:     80,
				y:     100,
				total: 100,
				pcnt:  0.1,
			},
			want: false,
		},
		{
			name: "within 10 percent under",
			args: args{
				x:     100,
				y:     91,
				total: 100,
				pcnt:  0.1,
			},
			want: true,
		},
		{
			name: "exactly 10 percent under",
			args: args{
				x:     100,
				y:     90,
				total: 100,
				pcnt:  0.1,
			},
			want: true,
		},
		{
			name: "more than 10 percent under",
			args: args{
				x:     100,
				y:     80,
				total: 100,
				pcnt:  0.1,
			},
			want: false,
		},
		{
			name: "percentages are accurate with low numbers",
			args: args{
				x:     1,
				y:     2,
				total: 100,
				pcnt:  0.01,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := distanceWithinNPcnt(tt.args.x, tt.args.y, tt.args.total, tt.args.pcnt); got != tt.want {
				t.Errorf("distanceWithinNPcnt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeComparisonSameLength(t *testing.T) {
	type args struct {
		targetString  string
		compareString string
		padding       []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "correct length already",
			args: args{
				targetString:  "foo",
				compareString: "foo",
				padding:       []string{"bar", "baz"},
			},
			want: "foo",
		},
		{
			name: "target shorter than compare shortens compare",
			args: args{
				targetString:  "foo",
				compareString: "foo bar",
				padding:       []string{"bar", "baz"},
			},
			want: "foo",
		},
		{
			name: "needs padding",
			args: args{
				targetString:  "foo bar",
				compareString: "foo",
				padding:       []string{"bar", "baz", "quix"},
			},
			want: "foo bar",
		},
		{
			name: "needs more padding",
			args: args{
				targetString:  "foo bar bar bar",
				compareString: "foo",
				padding:       []string{"bar", "baz", "quix"},
			},
			want: "foo bar baz quix",
		},
		{
			name: "too many spaces",
			args: args{
				targetString:  "foo   bar bar    bar",
				compareString: "foo",
				padding:       []string{"bar", "baz", "quix"},
			},
			want: "foo bar baz quix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeComparisonSameLengthAsTarget(tt.args.targetString, tt.args.compareString, tt.args.padding...); got != tt.want {
				t.Errorf("makeComparisonSameLengthAsTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withoutWords(t *testing.T) {
	type args struct {
		str   string
		words []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "remove words",
			args: args{
				str:   "yeah sure okay",
				words: []string{"okay", "sure"},
			},
			want: "yeah",
		},
		{
			name: "remove no words",
			args: args{
				str:   "yeah sure okay",
				words: []string{"foo", "bar"},
			},
			want: "yeah sure okay",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := withoutWords(tt.args.str, tt.args.words...); got != tt.want {
				t.Errorf("withoutWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
