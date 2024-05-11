package searchterms

import "github.com/warmans/rsk-search/pkg/filter"

func TermsToFilter(terms []Term) filter.Filter {
	var fil filter.Filter
	for _, t := range terms {
		if fil == nil {
			fil = &filter.CompFilter{
				Field: t.Field,
				Op:    t.Op,
				Value: filter.String(t.Value),
			}
		} else {
			fil = filter.And(fil, &filter.CompFilter{
				Field: t.Field,
				Op:    t.Op,
				Value: filter.String(t.Value),
			})
		}
	}
	return fil
}
