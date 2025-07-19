import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import {FormControl, UntypedFormControl, UntypedFormGroup, Validators} from '@angular/forms';
import { Subject } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';

export interface TranscriptMetadata {
  name: string;
  summary: string;
}

@Component({
    selector: 'app-metadata-editor',
    templateUrl: './metadata-editor.component.html',
    styleUrls: ['./metadata-editor.component.scss'],
    standalone: false
})
export class MetadataEditorComponent implements OnInit, OnDestroy {

  @Input()
  set transcriptMeta(value: TranscriptMetadata) {
    if (!value) {
      return;
    }
    this.form.get('name').setValue(value.name, {emitEvent: false});
    this.form.get('summary').setValue(value.summary, {emitEvent: false});
  }
  get transcriptMeta(): TranscriptMetadata {
    return { summary: this.form.get('summary').value, name: this.form.get('name').value };
  }

  @Input()
  allowEdit: boolean = false;

  @Output()
  metadataUpdated: EventEmitter<TranscriptMetadata> = new EventEmitter<TranscriptMetadata>();

  form: UntypedFormGroup = new UntypedFormGroup({
    name: new FormControl<string>(""),
    summary: new FormControl<string>("")
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
