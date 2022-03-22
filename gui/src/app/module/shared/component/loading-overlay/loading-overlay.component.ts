import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-loading-overlay',
  templateUrl: './loading-overlay.component.html',
  styleUrls: ['./loading-overlay.component.scss']
})
export class LoadingOverlayComponent implements OnInit {

  @Input()
  loading: boolean = false;

  @Input()
  fullScreen: boolean = true;

  constructor() {

  }

  ngOnInit(): void {
  }
}
