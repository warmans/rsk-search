import { Pipe, PipeTransform } from '@angular/core';
import { RskDialog } from 'src/app/lib/api-client/models';

@Pipe({
  name: 'matchedRowPos'
})
export class MatchedRowPosPipe implements PipeTransform {
  constructor() {
  }

  transform(lines: RskDialog[]): string {
    if ((lines || []).length === 0) {
      return '';
    }
    const posStart: string = `${lines[0].pos}` || '0';
    let endPos: string;
    if (lines.length > 1) {
      endPos = `${lines[lines.length - 1].pos}`;
    }
    return `pos-${posStart}${endPos ? '-' + endPos : ''}`;
  }
}
