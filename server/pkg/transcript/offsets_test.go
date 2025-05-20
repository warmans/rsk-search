package transcript

import (
	"bufio"
	"github.com/stretchr/testify/require"
	"github.com/warmans/rsk-search/pkg/models"
	"strings"
	"testing"
	"time"
)

func TestInferOffsetsSimple(t *testing.T) {
	ts, err := Import(bufio.NewScanner(strings.NewReader("#OFFSET: 1.00\nricky: Hello\nsteve: Hi")), "", 0)
	require.NoError(t, err)
	require.NotNil(t, ts)

	ts.Media.AudioDurationMs = (time.Second * 2).Milliseconds()

	ts = InferOffsets(ts)
	require.NotNil(t, ts)

	require.EqualValues(t, models.DialogTypeChat, ts.Transcript[0].Type)
	require.EqualValues(t, "1s", ts.Transcript[0].Timestamp.String())
	require.EqualValues(t, "ricky", ts.Transcript[0].Actor)

	require.EqualValues(t, models.DialogTypeChat, ts.Transcript[1].Type)
	require.EqualValues(t, "1.714s", ts.Transcript[1].Timestamp.String())
	require.EqualValues(t, "steve", ts.Transcript[1].Actor)
}

func TestInferOffsetsWithGap(t *testing.T) {
	ts, err := Import(bufio.NewScanner(strings.NewReader("#OFFSET: 1.00\nricky: Hello\n#GAP: 1m\nsteve: Hi")), "", 0)
	require.NoError(t, err)
	require.NotNil(t, ts)

	// audio is 1m5s
	ts.Media.AudioDurationMs = (time.Minute + (time.Second * 5)).Milliseconds()

	ts = InferOffsets(ts)
	require.NotNil(t, ts)

	require.EqualValues(t, models.DialogTypeChat, ts.Transcript[0].Type)
	require.EqualValues(t, "1s", ts.Transcript[0].Timestamp.String())
	require.EqualValues(t, "ricky", ts.Transcript[0].Actor)
	require.EqualValues(t, "2.857s", ts.Transcript[0].Duration.String())

	require.EqualValues(t, models.DialogTypeGap, ts.Transcript[1].Type)
	require.EqualValues(t, "1m0s", ts.Transcript[1].Duration.String())

	require.EqualValues(t, models.DialogTypeChat, ts.Transcript[2].Type)
	require.EqualValues(t, "1m3.857s", ts.Transcript[2].Timestamp.String()) // this just needs to be plausible.
	require.EqualValues(t, "steve", ts.Transcript[2].Actor)
	require.EqualValues(t, "1.143s", ts.Transcript[2].Duration.String())

	// the total duration less the initial 1s offset
	require.EqualValues(t, "1m4s", (ts.Transcript[0].Duration + ts.Transcript[1].Duration + ts.Transcript[2].Duration).String())
}
