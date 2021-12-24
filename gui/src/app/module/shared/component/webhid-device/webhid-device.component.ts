/// <reference types="w3c-web-hid" />

import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-webhid-device',
  templateUrl: './webhid-device.component.html',
  styleUrls: ['./webhid-device.component.scss']
})
export class WebhidDeviceComponent implements OnInit {

  hidAvailable: boolean = (navigator.hid !== undefined);

  constructor() {
  }

  ngOnInit(): void {

  }

  async listDevices() {

    if (!this.hidAvailable) {
      console.error("hid not possible");
      return;
    }

    let device: any;
    try {
      const devices = await navigator.hid.requestDevice({ filters: [] });
      device = devices[1];
    } catch (error) {
      console.warn('No device access granted', error);
      return;
    }

    if (device) {
      await device.open();
    }
  }

}
