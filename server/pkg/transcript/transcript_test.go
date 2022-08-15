package transcript

import (
	"github.com/stretchr/testify/require"
	"github.com/warmans/rsk-search/pkg/models"
	"testing"
)

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

func TestExport(t *testing.T) {
	type args struct {
		dialog   []models.Dialog
		synopsis []models.Synopsis
		trivia   []models.Trivia
		opts     []ExportOption
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "strip metadata",
			args: args{
				dialog: []models.Dialog{
					{
						Position:  1,
						OffsetSec: 1,
						Type:      models.DialogTypeChat,
						Actor:     "ricky",
						Content:   "foo",
					},
					{
						Position:  2,
						OffsetSec: 2,
						Type:      models.DialogTypeChat,
						Actor:     "karl",
						Content:   "bar",
					},
					{
						Position:  3,
						OffsetSec: 3,
						Type:      models.DialogTypeChat,
						Actor:     "steve",
						Content:   "baz",
					},
				},
				synopsis: []models.Synopsis{
					{
						Description: "whatever",
						StartPos:    2,
						EndPos:      3,
					},
				},
				trivia: []models.Trivia{
					{
						Description: "manysome trivia",
						StartPos:    1,
						EndPos:      3,
					},
				},
				opts: []ExportOption{WithStripMetadata()},
			},
			want: `ricky: foo
karl: bar
steve: baz
`,
			wantErr: false,
		},
		{
			name: "single line trivia renders OK",
			args: args{
				dialog: []models.Dialog{
					{
						Position:  1,
						OffsetSec: 1,
						Type:      models.DialogTypeChat,
						Actor:     "ricky",
						Content:   "foo",
					},
					{
						Position:  2,
						OffsetSec: 2,
						Type:      models.DialogTypeChat,
						Actor:     "karl",
						Content:   "bar",
					},
					{
						Position:  3,
						OffsetSec: 3,
						Type:      models.DialogTypeChat,
						Actor:     "steve",
						Content:   "baz",
					},
				},
				trivia: []models.Trivia{
					{
						Description: "single line of trivia",
						StartPos:    1,
						EndPos:      3,
					},
				},
			},
			want: `#OFFSET: 1
#TRIVIA: single line of trivia
ricky: foo
#OFFSET: 2
karl: bar
#OFFSET: 3
#/TRIVIA
steve: baz
`,
			wantErr: false,
		}, {
			name: "multi-line trivia renders OK",
			args: args{
				dialog: []models.Dialog{
					{
						Position:  1,
						OffsetSec: 1,
						Type:      models.DialogTypeChat,
						Actor:     "ricky",
						Content:   "foo",
					},
					{
						Position:  2,
						OffsetSec: 2,
						Type:      models.DialogTypeChat,
						Actor:     "karl",
						Content:   "bar",
					},
					{
						Position:  3,
						OffsetSec: 3,
						Type:      models.DialogTypeChat,
						Actor:     "steve",
						Content:   "baz",
					},
				},
				trivia: []models.Trivia{
					{
						Description: "many\nlines of\n trivia",
						StartPos:    1,
						EndPos:      3,
					},
				},
			},
			want: `#OFFSET: 1
#TRIVIA: many
# lines of
# trivia
ricky: foo
#OFFSET: 2
karl: bar
#OFFSET: 3
#/TRIVIA
steve: baz
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Export(tt.args.dialog, tt.args.synopsis, tt.args.trivia, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Export() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.EqualValues(t, tt.want, got)
		})
	}
}
