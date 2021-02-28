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
}
