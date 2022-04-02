import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'formatSeconds'
})
export class FormatSecondsPipe implements PipeTransform {
  constructor() {
  }
  transform(seconds: number): string {
    if (!seconds) {
      return "-";
    }
    return (new Date(seconds * 1000)).toISOString().substr(14, 5);
  }
}
