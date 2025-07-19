import { Component, Input, OnInit } from '@angular/core';

@Component({
    selector: 'app-loading-spinner',
    templateUrl: './loading-spinner.component.html',
    styleUrls: ['./loading-spinner.component.scss'],
    standalone: false
})
export class LoadingSpinnerComponent implements OnInit {

  @Input()
  loading: boolean = false;

  @Input()
  fullScreen: boolean = true;

  constructor() {

  }

  ngOnInit(): void {
  }
}
