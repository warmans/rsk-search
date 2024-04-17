package transcript

import (
	"bufio"
	"github.com/stretchr/testify/require"
	"github.com/warmans/rsk-search/pkg/models"
	"strings"
	"testing"
	"time"
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
						Timestamp: 0,
						Type:      models.DialogTypeChat,
						Actor:     "ricky",
						Content:   "Foo",
					},
					{
						Position:  2,
						Timestamp: 0,
						Type:      models.DialogTypeChat,
						Actor:     "karl",
						Content:   "Bar",
					},
					{
						Position:  3,
						Timestamp: 0,
						Type:      models.DialogTypeChat,
						Actor:     "steve",
						Content:   "Baz",
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
			want: `ricky: Foo
karl: Bar
steve: Baz
`,
			wantErr: false,
		},
		{
			name: "single line trivia renders OK",
			args: args{
				dialog: []models.Dialog{
					{
						Position:  1,
						Timestamp: time.Second * 1,
						Type:      models.DialogTypeChat,
						Actor:     "ricky",
						Content:   "Foo",
					},
					{
						Position:  2,
						Timestamp: time.Second * 2,
						Type:      models.DialogTypeChat,
						Actor:     "karl",
						Content:   "Bar",
					},
					{
						Position:  3,
						Timestamp: time.Second * 3,
						Type:      models.DialogTypeChat,
						Actor:     "steve",
						Content:   "Baz",
					},
				},
				trivia: []models.Trivia{
					{
						Description: "Single line of trivia",
						StartPos:    1,
						EndPos:      3,
					},
				},
			},
			want: `#OFFSET: 1.00
#TRIVIA: Single line of trivia
ricky: Foo
#OFFSET: 2.00
karl: Bar
#OFFSET: 3.00
#/TRIVIA
steve: Baz
`,
			wantErr: false,
		}, {
			name: "multi-line trivia renders OK",
			args: args{
				dialog: []models.Dialog{
					{
						Position:  1,
						Timestamp: time.Second * 1,
						Type:      models.DialogTypeChat,
						Actor:     "ricky",
						Content:   "Foo",
					},
					{
						Position:  2,
						Timestamp: time.Second * 2,
						Type:      models.DialogTypeChat,
						Actor:     "karl",
						Content:   "Bar",
					},
					{
						Position:  3,
						Timestamp: time.Second * 3,
						Type:      models.DialogTypeChat,
						Actor:     "steve",
						Content:   "Baz",
					},
				},
				trivia: []models.Trivia{
					{
						Description: "Many\nlines of\ntrivia",
						StartPos:    1,
						EndPos:      3,
					},
				},
			},
			want: `#OFFSET: 1.00
#TRIVIA: Many
# lines of
# trivia
ricky: Foo
#OFFSET: 2.00
karl: Bar
#OFFSET: 3.00
#/TRIVIA
steve: Baz
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

			// also check that importing the exported content works as expected
			dialog, synopsis, trivia, err := Import(bufio.NewScanner(strings.NewReader(got)), "", 0)
			require.NoError(t, err)
			for k := range dialog {
				require.EqualValues(t, tt.args.dialog[k].Timestamp, dialog[k].Timestamp)
				require.EqualValues(t, tt.args.dialog[k].Content, dialog[k].Content)
				require.EqualValues(t, tt.args.dialog[k].Actor, dialog[k].Actor)
			}
			for k := range trivia {
				require.EqualValues(t, tt.args.trivia[k].Description, trivia[k].Description)
			}
			for k := range synopsis {
				require.EqualValues(t, tt.args.synopsis[k].Description, synopsis[k].Description)
			}
		})
	}
}
