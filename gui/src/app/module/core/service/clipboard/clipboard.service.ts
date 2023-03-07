import { Injectable } from '@angular/core';
import { AlertService } from 'src/app/module/core/service/alert/alert.service';

@Injectable({
  providedIn: 'root'
})
export class ClipboardService {

  constructor(private alertService: AlertService) {

  }

  public copyTextToClipboard(text: string) {
    if (!navigator.clipboard) {
      this.alertService.danger('Clipboard is not available on your device.');
      return;
    }
    navigator.clipboard.writeText(text).then(() => {
      this.alertService.success("Copied")
    }, (err) => {
      this.alertService.danger(`Failed to copy to clipboard`, err);
    });
  }
}
