import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { SearchAPIClient } from 'src/app/lib/api-client/services/search';
import { AlertService } from '../../../core/service/alert/alert.service';
import { RskTscriptImport, RskTscriptImportList, RskChunkedTranscriptList, RskChunkedTranscriptStats } from 'src/app/lib/api-client/models';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-import',
  templateUrl: './import.component.html',
  styleUrls: ['./import.component.scss']
})
export class ImportComponent implements OnInit, OnDestroy {

  importForm: FormGroup = new FormGroup({
    'epid': new FormControl('preview-S1E06', [Validators.required]),
    'epname': new FormControl('', []),
    'mp3_uri': new FormControl('https://scrimpton.com/dl/media/episode/preview-S1E06.mp3', [Validators.required]),
  });

  chunkedTranscripts: RskChunkedTranscriptStats[] = [];
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

    this.importForm.get('mp3_uri').valueChanges.pipe(takeUntil(this.unsubscribe$)).subscribe((value: string) => {

      const filePattern = new RegExp(/.+\/(.*)\.mp3$/gm);
      let matches = filePattern.exec(value);
      if ((matches || []).length > 0) {
        this.importForm.get('epid').setValue(matches[1]);
      }
    });
  }

  updateTscriptList() {
    this.apiClient.listChunkedTranscripts().pipe(takeUntil(this.unsubscribe$)).subscribe((v: RskChunkedTranscriptList) => {
      this.chunkedTranscripts = v.chunked;
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
