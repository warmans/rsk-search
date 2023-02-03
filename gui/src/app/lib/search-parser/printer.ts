import { Term } from 'src/app/lib/search-parser/parser';
import { TermsToFilter } from 'src/app/lib/search-parser/util';
import { CompOp } from 'src/app/lib/filter-dsl/filter';

export function PrintFilterString(terms: Term[]): string {
  let filter = TermsToFilter(terms);
  return filter ? filter.print() : '';
}

export function PrintPlaintext(terms: Term[]): string {
  let out: string[] = [];
  terms.forEach((t) => {
    switch (t.field) {
      case 'content':
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
