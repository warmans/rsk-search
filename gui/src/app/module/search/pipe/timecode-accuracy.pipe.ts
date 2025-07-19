import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'timecodeAccuracy',
    standalone: false
})
export class TimecodeAccuracyPipe implements PipeTransform {
  constructor() {
  }

  transform(accuracyPcnt: number): string {
    switch (true) {
      case accuracyPcnt == 0:
        return 'Very Poor';
      case accuracyPcnt < 5 && accuracyPcnt > 0:
        return 'Poor';
      case accuracyPcnt >= 5 && accuracyPcnt < 10:
        return 'Average';
      case accuracyPcnt >= 10 && accuracyPcnt < 20:
        return 'Good';
      case accuracyPcnt >= 20 && accuracyPcnt < 40:
        return 'Very Good';
      case accuracyPcnt >= 40 :
        return 'Great';
    }
  }
}
