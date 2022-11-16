import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { Subject } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';

export interface TranscriptMetadata {
  summary: string;
}

@Component({
  selector: 'app-metadata-editor',
  templateUrl: './metadata-editor.component.html',
  styleUrls: ['./metadata-editor.component.scss']
})
export class MetadataEditorComponent implements OnInit, OnDestroy {

  @Input()
  set transcriptMeta(value: TranscriptMetadata) {
    if (!value) {
      return;
    }
    this.form.get('summary').setValue(value.summary, {emitEvent: false});
  }
  get transcriptMeta(): TranscriptMetadata {
    return { summary: this.form.get('summary').value };
  }

  @Input()
  allowEdit: boolean = false;

  @Output()
  metadataUpdated: EventEmitter<TranscriptMetadata> = new EventEmitter<TranscriptMetadata>();

  form: FormGroup = new FormGroup({
    summary: new FormControl()
  });

  destroy$: Subject<void> = new Subject<void>();

  constructor() {
  }

  ngOnInit(): void {
    this.form.valueChanges.pipe(takeUntil(this.destroy$)).subscribe((res) => {
      this.metadataUpdated.next(this.transcriptMeta);
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
