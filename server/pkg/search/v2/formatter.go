package v2

import "github.com/blugelabs/bluge/search/highlight"

const StartIdentifier = "{{"
const EndIdentifier = "}}"

func NewBracketFragmentFormatter() *BracketFragmentFormatter {
	return &BracketFragmentFormatter{}
}

type BracketFragmentFormatter struct {
}

func (a *BracketFragmentFormatter) Format(f *highlight.Fragment, orderedTermLocations highlight.TermLocations) string {
	rv := ""
	curr := f.Start
	for _, termLocation := range orderedTermLocations {
		if termLocation == nil {
			continue
		}
		if termLocation.Start < curr {
			continue
		}
		if termLocation.End > f.End {
			break
		}
		// add the stuff before this location
		rv += string(f.Orig[curr:termLocation.Start])
		// add the color
		rv += StartIdentifier
		// add the term itself
		rv += string(f.Orig[termLocation.Start:termLocation.End])
		// reset the color
		rv += EndIdentifier
		// update current
		curr = termLocation.End
	}
	// add any remaining text after the last token
	rv += string(f.Orig[curr:f.End])

	return rv
}
