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

func TestExportImport(t *testing.T) {
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
		}, {
			name: "multi-line synopsis renders OK",
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
				synopsis: []models.Synopsis{
					{
						Description: "Many\nlines of\nsynopsis",
						StartPos:    1,
						EndPos:      3,
					},
				},
			},
			want: `#OFFSET: 1.00
#SYN: Many
# lines of
# synopsis
ricky: Foo
#OFFSET: 2.00
karl: Bar
#OFFSET: 3.00
#/SYN
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
			ts, err := Import(bufio.NewScanner(strings.NewReader(got)), "", 0)
			require.NoError(t, err)
			for k := range ts.Transcript {
				require.EqualValues(t, tt.args.dialog[k].Timestamp, ts.Transcript[k].Timestamp)
				require.EqualValues(t, tt.args.dialog[k].Content, ts.Transcript[k].Content)
				require.EqualValues(t, tt.args.dialog[k].Actor, ts.Transcript[k].Actor)
			}
			for k := range ts.Trivia {
				require.EqualValues(t, tt.args.trivia[k].Description, ts.Trivia[k].Description)
			}
			for k := range ts.Synopsis {
				require.EqualValues(t, tt.args.synopsis[k].Description, ts.Synopsis[k].Description)
			}
		})
	}
}

func TestReadTagNearOffset(t *testing.T) {
	text := `
#TRIVIA: Many
#OFFSET: 123
`

	ts, err := Import(bufio.NewScanner(strings.NewReader(text)), "", 0)
	require.NoError(t, err)
	require.EqualValues(t, ts.Trivia[0].Description, "Many")
}

func TestReadMultilineTagNearOffset(t *testing.T) {
	text := `
#TRIVIA: Many
# lines of
# text
#OFFSET: 123
`

	ts, err := Import(bufio.NewScanner(strings.NewReader(text)), "", 0)
	require.NoError(t, err)
	require.EqualValues(t, "Many\nlines of\ntext", ts.Trivia[0].Description)
}

func TestReadGap(t *testing.T) {
	text := `
#OFFSET: 1
Steve: 1
#OFFSET: 12
Ricky: 2
#GAP: 1m5s
Karl: 3
`

	ts, err := Import(bufio.NewScanner(strings.NewReader(text)), "", 0)
	require.NoError(t, err)
	require.NotNil(t, ts.Transcript)

	require.EqualValues(t, models.DialogTypeChat, ts.Transcript[0].Type)
	require.EqualValues(t, "1", ts.Transcript[0].Content)
	require.EqualValues(t, time.Second, ts.Transcript[0].Timestamp)

	require.EqualValues(t, models.DialogTypeChat, ts.Transcript[1].Type)
	require.EqualValues(t, "2", ts.Transcript[1].Content)
	require.EqualValues(t, time.Second*12, ts.Transcript[1].Timestamp)

	require.EqualValues(t, models.DialogTypeGap, ts.Transcript[2].Type)
	require.EqualValues(t, time.Minute+(5*time.Second), ts.Transcript[2].Duration)

	require.EqualValues(t, models.DialogTypeChat, ts.Transcript[3].Type)
	require.EqualValues(t, "3", ts.Transcript[3].Content)
}
