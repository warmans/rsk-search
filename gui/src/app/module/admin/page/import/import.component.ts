import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { AlertService } from '../../../core/service/alert/alert.service';

@Component({
  selector: 'app-import',
  templateUrl: './import.component.html',
  styleUrls: ['./import.component.scss']
})
export class ImportComponent implements OnInit {

  importForm: FormGroup = new FormGroup({
    'epid': new FormControl('xfm-S1E01', [Validators.required]),
    'mp3_uri': new FormControl('https://storage.googleapis.com/scrimpton-raw-audio/xfm-S1E01.mp3', [Validators.required]),
  });

  constructor(private apiClient: SearchAPIClient, private alerts: AlertService) {
  }

  ngOnInit(): void {
  }

  startImport() {
    if (!this.importForm.valid) {
      return;
    }
    this.apiClient.createTscriptImport({
      body: {
        epid: this.importForm.get('epid').value,
        mp3Uri: this.importForm.get('mp3_uri').value,
      }
    }).subscribe((done) => {
      this.alerts.success("Importing has started...")
    });
  }
}
