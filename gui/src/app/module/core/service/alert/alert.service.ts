import { EventEmitter, Injectable } from '@angular/core';

export interface Alert {
  level: string;
  content: string;
  details?: string[];
}

const AUTO_REMOVE_AFTER_MS: number = 1000 * 10;

@Injectable({
  providedIn: 'root'
})
export class AlertService {

  alertsUpdated: EventEmitter<Alert[]> = new EventEmitter<Alert[]>();

  alerts: Alert[] = [];

  constructor() {
  }

  danger(c: string, ...d: string[]) {
    this.setAlert('danger', c, d);
  }

  success(c: string, ...d: string[]) {
    this.setAlert('success', c, d);
  }

  private setAlert(l: string, c: string, d: string[]) {
    this.alerts.push({ level: l, content: c, details: d });
    this.alertsUpdated.next(this.alerts);
    this.cleanup(c);
  }

  cleanup(c: string) {
    setTimeout(() => {
      for (const k of this.alerts.keys()) {
        if (this.alerts[k].content === c) {
          this.alerts.splice(k, 1);
        }
      }
      this.alertsUpdated.next(this.alerts);
    }, AUTO_REMOVE_AFTER_MS);
  }
}
