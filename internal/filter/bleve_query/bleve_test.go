package bleve_query

import (
	"fmt"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/stretchr/testify/require"
	"github.com/warmans/rsk-search/internal/filter"
	"testing"
)

func TestFilterToQuery(t *testing.T) {
	f, err := filter.Parse(`foo = "bar" and bar ~= "baz" and baz != 1 and cat > 10`)
	require.NoError(t, err)
	q, err := FilterToQuery(f)
	require.NoError(t, err)

	str, err := query.DumpQuery(mapping.NewIndexMapping(), q)
	require.NoError(t, err)

	fmt.Println(str)
}
