/// <reference types="w3c-web-usb" />

import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-webusb-device',
  templateUrl: './webusb-device.component.html',
  styleUrls: ['./webusb-device.component.scss']
})
export class WebusbDeviceComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
  }

  async connect() {
    let device = await navigator.usb.requestDevice({
      filters: []
    });
    console.log(device);
    await device.open();
    await device.selectConfiguration(1);
    await device.claimInterface(0);

    const listen = async () => {
      const result = await device.transferIn(3, 64);
      const decoder = new TextDecoder();
      const message = decoder.decode(result.data);
      console.log(message);

      const messageParts = message.split(' = ');
      // if (messageParts[0] === 'Count') {
      //   deviceHeartbeat.innerText = messageParts[1];
      // } else if (messageParts[0] === 'Button'
      //   && messageParts[1] === '1') {
      //
      //   deviceButtonPressed.innerText = new Date()
      //     .toLocaleString('en-ZA', {
      //       hour: 'numeric',
      //       minute: 'numeric',
      //       second: 'numeric',
      //     });
      // }
      listen();
    };


  }

}
