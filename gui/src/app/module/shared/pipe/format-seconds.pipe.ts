import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'formatSeconds'
})
export class FormatSecondsPipe implements PipeTransform {
  constructor() {
  }

  transform(seconds: number | string): string {
    if (!seconds) {
      return '-';
    }

    const secondsNum: number = (typeof seconds === 'string') ? parseInt(seconds) : seconds;
    const minsNum: number = secondsNum / 60;

    const mins: string = String(minsNum.toFixed(0)).padStart(2, '0');
    const secs: string = String(((minsNum % 1) * 60).toFixed(0)).padStart(2, '0');
    return `${mins}:${secs}`;
  }
}
