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

  danger(content: string, ...details: string[]) {
    this.setAlert('danger', content, details);
  }

  success(content: string, ...details: string[]) {
    this.setAlert('success', content, details);
  }

  private setAlert(level: string, content: string, details: string[]) {
    this.alerts.push({ level: level, content: content, details: details });
    this.alertsUpdated.next(this.alerts);
    this.cleanup(content);
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
