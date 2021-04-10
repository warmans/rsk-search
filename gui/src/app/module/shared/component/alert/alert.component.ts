import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { Alert, AlertService } from '../../../core/service/alert/alert.service';

@Component({
  selector: 'app-alert',
  templateUrl: './alert.component.html',
  styleUrls: ['./alert.component.scss']
})
export class AlertComponent implements OnInit, OnDestroy {

  alerts: Alert[] = [];

  destroy: EventEmitter<void> = new EventEmitter<void>();

  constructor(public alertService: AlertService) {
    this.alertService.alertsUpdated.subscribe(
      (alerts: Alert[]) => {
        this.alerts = alerts;
      }
    );
  }

  ngOnDestroy(): void {
    this.destroy.complete();
  }

  ngOnInit(): void {
  }

  remove(c: string): void {
    this.alertService.remove(c);
  }

  createTestData() {
    this.alerts.push({level:"success", content: "Update success", details: ["detail 1", "detail 2"]});
    this.alerts.push({level:"danger", content: "Update danger", details: ["detail 1", "detail 2"]});
    this.alerts.push({level:"warning", content: "Update warning", details: ["detail 1", "detail 2"]});
  }
}
