import { Term } from 'lib/search-parser/parser';
import { And, Filter } from 'lib/filter-dsl/filter';

export function TermsToFilter(terms: Term[]): Filter {
  let filter: Filter;
  terms.forEach((term: Term) => {
    if (!filter) {
      filter = term.toFilter();
    } else {
      filter = And(filter, term.toFilter());
    }
  });
  return filter;
}
