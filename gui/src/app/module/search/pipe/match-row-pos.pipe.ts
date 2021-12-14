import { Pipe, PipeTransform } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { RskDialog } from '../../../lib/api-client/models';

@Pipe({
  name: 'matchedRowPos'
})
export class MatchedRowPosPipe implements PipeTransform {
  constructor() {
  }
  transform(lines: RskDialog[]): string {
    const pos = lines.find((l: RskDialog) => l.isMatchedRow)?.pos || "0";
    return `pos-${pos}`
  }
}
