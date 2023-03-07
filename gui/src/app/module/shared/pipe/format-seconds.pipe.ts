import { Pipe, PipeTransform } from '@angular/core';
import { formatSecondsAsTimestamp } from 'src/app/lib/util';

@Pipe({
  name: 'formatSeconds'
})
export class FormatSecondsPipe implements PipeTransform {

  constructor() {
  }

  transform(seconds: number | string, secondsAreMilliseconds?: boolean): string {
    return formatSecondsAsTimestamp(seconds, secondsAreMilliseconds);
  }
}
