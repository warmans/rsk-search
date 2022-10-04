import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { SearchAPIClient } from '../../../../lib/api-client/services/search';
import { AlertService } from '../../../core/service/alert/alert.service';
import { RskTscriptImport, RskTscriptImportList, RskTscriptList, RskTscriptStats } from '../../../../lib/api-client/models';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-import',
  templateUrl: './import.component.html',
  styleUrls: ['./import.component.scss']
})
export class ImportComponent implements OnInit, OnDestroy {

  importForm: FormGroup = new FormGroup({
    'epid': new FormControl('sample-S0E0', [Validators.required]),
    'epname': new FormControl('', []),
    'mp3_uri': new FormControl('https://storage.googleapis.com/scrimpton-raw-audio/sample-S0E0.mp3', [Validators.required]),
  });

  tscripts: RskTscriptStats[] = [];
  imports: RskTscriptImport[] = [];

  private unsubscribe$: EventEmitter<void> = new EventEmitter<void>();

  constructor(private apiClient: SearchAPIClient, private alerts: AlertService) {
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  ngOnInit(): void {
    this.updateTscriptList();
    this.updateTscriptImportList();
  }

  updateTscriptList() {
    this.apiClient.listTscripts().pipe(takeUntil(this.unsubscribe$)).subscribe((v: RskTscriptList) => {
      this.tscripts = v.tscripts;
    });
  }

  updateTscriptImportList() {
    this.apiClient.listTscriptImports({ sortField: 'created_at', sortDirection: 'desc' }).pipe(takeUntil(this.unsubscribe$)).subscribe((v: RskTscriptImportList) => {
      this.imports = v.imports;
    });
  }

  startImport() {
    if (!this.importForm.valid) {
      return;
    }
    this.apiClient.createTscriptImport({
      body: {
        epid: this.importForm.get('epid').value,
        epname: this.importForm.get('epname').value,
        mp3Uri: this.importForm.get('mp3_uri').value,
      }
    }).pipe(takeUntil(this.unsubscribe$)).subscribe((done) => {
      this.alerts.success('Importing has started...');
    });
  }

  deleteTscript(id: string) {
    this.apiClient.deleteTscript({ id: id }).pipe(takeUntil(this.unsubscribe$)).subscribe(() => {
      this.alerts.success('Deleted');
      this.updateTscriptList();
    });
  }
}
