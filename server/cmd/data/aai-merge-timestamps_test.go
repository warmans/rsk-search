package data

import "testing"

func Test_withinNPcnt(t *testing.T) {
	type args struct {
		x    int
		y    int
		pcnt float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "within 10 percent over",
			args: args{
				x:    100,
				y:    109,
				pcnt: 0.1,
			},
			want: true,
		},
		{
			name: "exactly 10 percent over",
			args: args{
				x:    100,
				y:    110,
				pcnt: 0.1,
			},
			want: true,
		},
		{
			name: "more than 10 percent over",
			args: args{
				x:    100,
				y:    120,
				pcnt: 0.1,
			},
			want: false,
		},
		{
			name: "within 10 percent under",
			args: args{
				x:    100,
				y:    91,
				pcnt: 0.1,
			},
			want: true,
		},
		{
			name: "exactly 10 percent under",
			args: args{
				x:    100,
				y:    90,
				pcnt: 0.1,
			},
			want: true,
		},
		{
			name: "more than 10 percent under",
			args: args{
				x:    100,
				y:    80,
				pcnt: 0.1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := distanceWithinNPcnt(tt.args.x, tt.args.y, tt.args.pcnt); got != tt.want {
				t.Errorf("distanceWithinNPcnt() = %v, want %v", got, tt.want)
			}
		})
	}
}
