import { Pipe, PipeTransform } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';

@Pipe({
  name: 'formatSeconds'
})
export class FormatSecondsPipe implements PipeTransform {
  constructor() {
  }
  transform(seconds: number): string {
    if (seconds == 0) {
      return "";
    }
    return (new Date(seconds * 1000)).toISOString().substr(14, 5);
  }
}
