import {Term} from 'src/app/lib/search-parser/parser';
import {TermsToFilter} from 'src/app/lib/search-parser/util';
import {CompOp, Filter} from 'src/app/lib/filter-dsl/filter';
import {Tag} from "./scanner";

export function PrintFilterString(terms: Term[]): string {
  let filter: Filter = TermsToFilter(terms);
  return filter ? filter.print() : '';
}

export function PrintPlaintext(terms: Term[]): string {
  let out: string[] = [];
  terms.forEach((t: Term): void => {
    switch (t.field) {
      case 'content':
        if (t.tok.tag === Tag.Regexp) {
          out.push(t.op === CompOp.Eq ? `/${t.value}/` : t.value);
          break;
        }
        out.push(t.op === CompOp.Eq ? `"${t.value}"` : t.value);
        break;
      case 'actor':
        out.push(`@${t.value}`);
        break;
      case 'publication':
        out.push(`~${t.value}`);
        break;
    }
  });

  return out.join(' ');
}
