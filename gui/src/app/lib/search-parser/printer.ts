import { Term } from 'src/app/lib/search-parser/parser';
import { TermsToFilter } from 'src/app/lib/search-parser/util';

export function PrintFilterString(terms: Term[]): string {
  let filter = TermsToFilter(terms);
  return filter ? filter.print() : '';
}

export function PrintPlaintext(terms: Term[]): string {
  let out: string[] = [];
  terms.forEach((t) => {
    switch (t.field) {
      case 'content':
        out.push(t.value);
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
