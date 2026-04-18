import { Component, Input, OnInit } from '@angular/core';
import { NgClass } from '@angular/common';

@Component({
  selector: 'app-loading-overlay',
  templateUrl: './loading-overlay.component.html',
  styleUrls: ['./loading-overlay.component.scss'],
  imports: [NgClass],
})
export class LoadingOverlayComponent implements OnInit {
  @Input()
  loading: boolean = false;

  @Input()
  fullScreen: boolean = true;

  constructor() {}

  ngOnInit(): void {}
}
